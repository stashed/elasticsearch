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
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	meta_util "kmodules.xyz/client-go/meta"
	appcatalog "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	appcatalog_cs "kmodules.xyz/custom-resources/client/clientset/versioned"
)

const (
	ESUser              = "ADMIN_USERNAME"
	ESPassword          = "ADMIN_PASSWORD"
	MultiElasticDumpCMD = "multielasticdump"
	ESCACertFile        = "root.pem"
	ESAuthFile          = "auth.txt"
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
