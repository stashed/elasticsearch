package pkg

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/appscode/go/log"
	"k8s.io/client-go/kubernetes"
	appcatalog_cs "kmodules.xyz/custom-resources/client/clientset/versioned"
	"stash.appscode.dev/stash/pkg/restic"
)

const (
	JobESBackup  = "stash-es-backup"
	ESUser       = "ADMIN_USERNAME"
	ESPassword   = "ADMIN_PASSWORD"
	ESDumpCMD    = "multielasticdump"
	ESCACertFile = "root.pem"
	ESDataDir    = "/var/pv/data"
)

type esOptions struct {
	kubeClient    kubernetes.Interface
	catalogClient appcatalog_cs.Interface

	namespace      string
	appBindingName string
	esArgs         string
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
