package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"stash.appscode.dev/stash/pkg/restic"

	"github.com/appscode/go/log"
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
	catalogClient appcatalog_cs.Interface

	namespace      string
	appBindingName string
	esArgs         string
	interimDataDir string
	outputDir      string

	setupOptions   restic.SetupOptions
	backupOptions  restic.BackupOptions
	restoreOptions restic.RestoreOptions
}

func waitForDBReady(host string, port int32) {
	log.Infoln("Checking database connection")
	cmd := fmt.Sprintf(`nc "%s" "%d" -w 30`, host, port)
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
