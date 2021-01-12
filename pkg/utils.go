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
	"os"
	"os/exec"
	"time"

	stash "stash.appscode.dev/apimachinery/client/clientset/versioned"
	"stash.appscode.dev/apimachinery/pkg/restic"

	"gomodules.xyz/x/log"
	"k8s.io/client-go/kubernetes"
	appcatalog_cs "kmodules.xyz/custom-resources/client/clientset/versioned"
)

const (
	ESUser              = "ADMIN_USERNAME"
	ESPassword          = "ADMIN_PASSWORD"
	MultiElasticDumpCMD = "multielasticdump"
	ESCACertFile        = "root.pem"
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
	log.Infoln("Checking database connection")
	cmd := fmt.Sprintf(`nc "%s" "%d" -w %d`, host, port, waitTimeout)
	for {
		if err := exec.Command(cmd).Run(); err != nil {
			break
		}
		log.Infoln("Waiting... database is not ready yet")
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
