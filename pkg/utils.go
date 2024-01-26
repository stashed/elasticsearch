/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Free Trial License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Free-Trial-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pkg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	stash "stash.appscode.dev/apimachinery/client/clientset/versioned"
	"stash.appscode.dev/apimachinery/pkg/restic"

	shell "gomodules.xyz/go-sh"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	meta_util "kmodules.xyz/client-go/meta"
	appcatalog "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	appcatalog_cs "kmodules.xyz/custom-resources/client/clientset/versioned"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	esapi "kubedb.dev/apimachinery/apis/elasticsearch/v1alpha1"
	kubedbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	es_dashboard "kubedb.dev/db-client-go/elasticsearchdashboard"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

const (
	ESUser               = "ADMIN_USERNAME"
	ESPassword           = "ADMIN_PASSWORD"
	MultiElasticDumpCMD  = "multielasticdump"
	ESCACertFile         = "root.pem"
	ESAuthFile           = "auth.txt"
	DashboardObjectsFile = "dashboard.ndjson"
)

type esOptions struct {
	kubeClient    kubernetes.Interface
	stashClient   stash.Interface
	catalogClient appcatalog_cs.Interface

	namespace           string
	backupSessionName   string
	appBindingName      string
	appBindingNamespace string
	esArgs              string
	interimDataDir      string
	enableDashboard     bool
	outputDir           string
	storageSecret       kmapi.ObjectReference
	waitTimeout         int32

	setupOptions   restic.SetupOptions
	backupOptions  restic.BackupOptions
	restoreOptions restic.RestoreOptions
	config         *restclient.Config
}
type sessionWrapper struct {
	sh  *shell.Session
	cmd *restic.Command
}

func (opt *esOptions) newSessionWrapper(cmd string) *sessionWrapper {
	return &sessionWrapper{
		sh: shell.NewSession(),
		cmd: &restic.Command{
			Name: cmd,
		},
	}
}

func (opt *esOptions) setDatabaseCredentials(appBinding *appcatalog.AppBinding, cmd *restic.Command) error {
	appBindingSecret, err := opt.kubeClient.CoreV1().Secrets(appBinding.Namespace).Get(context.TODO(), appBinding.Spec.Secret.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	err = appBinding.TransformSecret(opt.kubeClient, appBindingSecret.Data)
	if err != nil {
		return err
	}
	// write the credential ifo into a file
	// TODO: support backup without authentication
	httpAuthFilePath := filepath.Join(opt.setupOptions.ScratchDir, ESAuthFile)
	err = writeAuthFile(httpAuthFilePath, appBindingSecret)
	if err != nil {
		return err
	}
	cmd.Args = append(cmd.Args, fmt.Sprintf("--httpAuthFile=%s", httpAuthFilePath))
	return nil
}

func (session *sessionWrapper) setUserArgs(args string) {
	for _, arg := range strings.Fields(args) {
		session.cmd.Args = append(session.cmd.Args, arg)
	}
}

func (session *sessionWrapper) setTLSParameters(appBinding *appcatalog.AppBinding, scratchDir string) error {
	session.sh.SetEnv("NODE_TLS_REJECT_UNAUTHORIZED", "0") // xref: https://github.com/taskrabbit/elasticsearch-dump#bypassing-self-sign-certificate-errors
	if appBinding.Spec.ClientConfig.CABundle != nil {
		if err := os.WriteFile(filepath.Join(scratchDir, ESCACertFile), appBinding.Spec.ClientConfig.CABundle, os.ModePerm); err != nil {
			return err
		}
		session.cmd.Args = append(session.cmd.Args, fmt.Sprintf("--ca-input=%v", filepath.Join(scratchDir, ESCACertFile)))
	}
	return nil
}

func (opt esOptions) waitForDBReady(appBinding *appcatalog.AppBinding) error {
	hostname, err := appBinding.Hostname()
	if err != nil {
		return err
	}

	port, err := appBinding.Port()
	if err != nil {
		return err
	}

	klog.Infoln("Checking database connection")
	cmd := fmt.Sprintf(`nc "%s" "%d" -w %d`, hostname, port, opt.waitTimeout)
	for {
		if err := exec.Command(cmd).Run(); err != nil {
			break
		}
		klog.Infoln("Waiting... database is not ready yet")
		time.Sleep(5 * time.Second)
	}
	klog.Infoln("Performing multielasticdump on", hostname)
	return nil
}

func clearDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("unable to clean datadir: %v. Reason: %v", dir, err)
	}

	return os.MkdirAll(dir, os.ModePerm)
}

func must(v []byte, err error) string {
	if err != nil {
		panic(err)
	}
	return string(v)
}

func writeAuthFile(filename string, cred *core.Secret) error {
	authKeys := fmt.Sprintf("user=%s\npassword=%q",
		must(meta_util.GetBytesForKeys(cred.Data, core.BasicAuthUsernameKey, ESUser)),
		must(meta_util.GetBytesForKeys(cred.Data, core.BasicAuthPasswordKey, ESPassword)),
	)
	return os.WriteFile(filename, []byte(authKeys), 0o400) // only readable to owner
}

func newRuntimeClient(cfg *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	utilruntime.Must(core.AddToScheme(scheme))
	utilruntime.Must(esapi.AddToScheme(scheme))
	utilruntime.Must(kubedbapi.AddToScheme(scheme))

	hc, err := rest.HTTPClientFor(cfg)
	if err != nil {
		return nil, err
	}
	mapper, err := apiutil.NewDynamicRESTMapper(cfg, hc)
	if err != nil {
		return nil, err
	}

	return client.New(cfg, client.Options{
		Scheme: scheme,
		Mapper: mapper,
	})
}

func getElasticSearchDashboard(klient client.Client, appBinding *appcatalog.AppBinding) (*esapi.ElasticsearchDashboard, error) {
	dashboards := &esapi.ElasticsearchDashboardList{}
	opts := []client.ListOption{client.InNamespace(appBinding.Namespace)}
	if err := klient.List(context.TODO(), dashboards, opts...); err != nil {
		return nil, err
	}

	for _, dashboard := range dashboards.Items {
		if dashboard.Spec.DatabaseRef != nil {
			if dashboard.Spec.DatabaseRef.Name == appBinding.Name {
				return &dashboard, nil
			}
		}
	}

	return nil, fmt.Errorf("no elasticsearch dashboard found")
}

func getElasticSearch(klient client.Client, appBinding *appcatalog.AppBinding) (*kubedbapi.Elasticsearch, error) {
	es := &kubedbapi.Elasticsearch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appBinding.Name,
			Namespace: appBinding.Namespace,
		},
	}

	if err := klient.Get(context.TODO(), client.ObjectKeyFromObject(es), es); err != nil {
		return nil, err
	}

	return es, nil
}

func getSecret(klient client.Client, obj kmapi.ObjectReference) (*core.Secret, error) {
	sec := &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      obj.Name,
			Namespace: obj.Namespace,
		},
	}

	if err := klient.Get(context.TODO(), client.ObjectKeyFromObject(sec), sec); err != nil {
		return nil, err
	}

	return sec, nil
}

func getVersionInfo(es *kubedbapi.Elasticsearch, appBinding *appcatalog.AppBinding) *es_dashboard.DbVersionInfo {
	authPlugin := catalog.ElasticsearchAuthPluginOpenSearch
	segments := strings.Split(es.Spec.Version, "-")
	if segments[0] == "xpack" {
		authPlugin = catalog.ElasticsearchAuthPluginXpack
	}
	return &es_dashboard.DbVersionInfo{
		Name:       es.Spec.Version,
		Version:    appBinding.Spec.Version,
		AuthPlugin: authPlugin,
	}
}

func (opt esOptions) getDashboardClient(appBinding *appcatalog.AppBinding) (*es_dashboard.Client, error) {
	klient, err := newRuntimeClient(opt.config)
	if err != nil {
		return nil, err
	}

	esDashboard, err := getElasticSearchDashboard(klient, appBinding)
	if err != nil {
		return nil, err
	}

	es, err := getElasticSearch(klient, appBinding)
	if err != nil {
		return nil, err
	}

	sec, err := getSecret(klient, kmapi.ObjectReference{
		Name:      es.Spec.AuthSecret.Name,
		Namespace: es.Namespace,
	})
	if err != nil {
		return nil, err
	}

	versionInfo := getVersionInfo(es, appBinding)

	return es_dashboard.NewKubeDBClientBuilder(klient, esDashboard).
		WithContext(context.TODO()).
		WithDatabaseRef(es).
		WithAuthSecret(sec).
		WithDbVersionInfo(versionInfo).
		GetElasticsearchDashboardClient()
}
