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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	stash "stash.appscode.dev/apimachinery/client/clientset/versioned"
	"stash.appscode.dev/apimachinery/pkg/restic"

	core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	meta_util "kmodules.xyz/client-go/meta"
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

	namespace         string
	backupSessionName string
	appBindingName    string
	esArgs            string
	interimDataDir    string
	outputDir         string
	waitTimeout       int32

	setupOptions   restic.SetupOptions
	backupOptions  restic.BackupOptions
	restoreOptions restic.RestoreOptions
}

func waitForDBReady(host string, port, waitTimeout int32) {
	klog.Infoln("Checking database connection")
	cmd := fmt.Sprintf(`nc "%s" "%d" -w %d`, host, port, waitTimeout)
	for {
		if err := exec.Command(cmd).Run(); err != nil {
			break
		}
		klog.Infoln("Waiting... database is not ready yet")
		time.Sleep(5 * time.Second)
	}
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
	return ioutil.WriteFile(filename, []byte(authKeys), 0400) // only redable to owner
}
