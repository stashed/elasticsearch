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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	api_v1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	stash "stash.appscode.dev/apimachinery/client/clientset/versioned"
	"stash.appscode.dev/apimachinery/pkg/restic"
	api_util "stash.appscode.dev/apimachinery/pkg/util"

	"github.com/spf13/cobra"
	license "go.bytebuilders.dev/license-verifier/kubernetes"
	"gomodules.xyz/flags"
	"gomodules.xyz/go-sh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	appcatalog "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	appcatalog_cs "kmodules.xyz/custom-resources/client/clientset/versioned"
	v1 "kmodules.xyz/offshoot-api/api/v1"
)

func NewCmdBackup() *cobra.Command {
	var (
		masterURL      string
		kubeconfigPath string
		opt            = esOptions{
			waitTimeout: 300,
			setupOptions: restic.SetupOptions{
				ScratchDir:  restic.DefaultScratchDir,
				EnableCache: false,
			},
			backupOptions: restic.BackupOptions{
				Host: restic.DefaultHost,
			},
		}
	)

	cmd := &cobra.Command{
		Use:               "backup-es",
		Short:             "Takes a backup of Elasticsearch DB",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.EnsureRequiredFlags(cmd, "appbinding", "provider", "secret-dir")
			time.Sleep(time.Second * 5)

			// prepare client
			config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
			if err != nil {
				return err
			}
			err = license.CheckLicenseEndpoint(config, licenseApiService, SupportedProducts)
			if err != nil {
				return err
			}
			opt.kubeClient, err = kubernetes.NewForConfig(config)
			if err != nil {
				return err
			}
			opt.stashClient, err = stash.NewForConfig(config)
			if err != nil {
				return err
			}
			opt.catalogClient, err = appcatalog_cs.NewForConfig(config)
			if err != nil {
				return err
			}

			targetRef := api_v1beta1.TargetRef{
				APIVersion: appcatalog.SchemeGroupVersion.String(),
				Kind:       appcatalog.ResourceKindApp,
				Name:       opt.appBindingName,
			}
			var backupOutput *restic.BackupOutput
			backupOutput, err = opt.backupElasticsearch(targetRef)
			if err != nil {
				backupOutput = &restic.BackupOutput{
					BackupTargetStatus: api_v1beta1.BackupTargetStatus{
						Ref: targetRef,
						Stats: []api_v1beta1.HostBackupStats{
							{
								Hostname: opt.backupOptions.Host,
								Phase:    api_v1beta1.HostBackupFailed,
								Error:    err.Error(),
							},
						},
					},
				}
			}
			// If output directory specified, then write the output in "output.json" file in the specified directory
			if opt.outputDir != "" {
				return backupOutput.WriteOutput(filepath.Join(opt.outputDir, restic.DefaultOutputFileName))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&opt.esArgs, "es-args", opt.esArgs, "Additional arguments")
	cmd.Flags().Int32Var(&opt.waitTimeout, "wait-timeout", opt.waitTimeout, "Number of seconds to wait for the database to be ready")

	cmd.Flags().StringVar(&masterURL, "master", masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&opt.namespace, "namespace", "default", "Namespace of Backup/Restore Session")
	cmd.Flags().StringVar(&opt.backupSessionName, "backupsession", opt.backupSessionName, "Name of the Backup Session")
	cmd.Flags().StringVar(&opt.appBindingName, "appbinding", opt.appBindingName, "Name of the app binding")

	cmd.Flags().StringVar(&opt.setupOptions.Provider, "provider", opt.setupOptions.Provider, "Backend provider (i.e. gcs, s3, azure etc)")
	cmd.Flags().StringVar(&opt.setupOptions.Bucket, "bucket", opt.setupOptions.Bucket, "Name of the cloud bucket/container (keep empty for local backend)")
	cmd.Flags().StringVar(&opt.setupOptions.Endpoint, "endpoint", opt.setupOptions.Endpoint, "Endpoint for s3/s3 compatible backend or REST server URL")
	cmd.Flags().StringVar(&opt.setupOptions.Region, "region", opt.setupOptions.Region, "Region for s3/s3 compatible backend")
	cmd.Flags().StringVar(&opt.setupOptions.Path, "path", opt.setupOptions.Path, "Directory inside the bucket where backup will be stored")
	cmd.Flags().StringVar(&opt.setupOptions.SecretDir, "secret-dir", opt.setupOptions.SecretDir, "Directory where storage secret has been mounted")
	cmd.Flags().StringVar(&opt.setupOptions.ScratchDir, "scratch-dir", opt.setupOptions.ScratchDir, "Temporary directory")
	cmd.Flags().BoolVar(&opt.setupOptions.EnableCache, "enable-cache", opt.setupOptions.EnableCache, "Specify whether to enable caching for restic")
	cmd.Flags().Int64Var(&opt.setupOptions.MaxConnections, "max-connections", opt.setupOptions.MaxConnections, "Specify maximum concurrent connections for GCS, Azure and B2 backend")

	cmd.Flags().StringVar(&opt.backupOptions.Host, "hostname", opt.backupOptions.Host, "Name of the host machine")
	cmd.Flags().StringVar(&opt.interimDataDir, "interim-data-dir", opt.interimDataDir, "Directory where the targeted data will be stored temporarily before uploading to the backend.")

	cmd.Flags().Int64Var(&opt.backupOptions.RetentionPolicy.KeepLast, "retention-keep-last", opt.backupOptions.RetentionPolicy.KeepLast, "Specify value for retention strategy")
	cmd.Flags().Int64Var(&opt.backupOptions.RetentionPolicy.KeepHourly, "retention-keep-hourly", opt.backupOptions.RetentionPolicy.KeepHourly, "Specify value for retention strategy")
	cmd.Flags().Int64Var(&opt.backupOptions.RetentionPolicy.KeepDaily, "retention-keep-daily", opt.backupOptions.RetentionPolicy.KeepDaily, "Specify value for retention strategy")
	cmd.Flags().Int64Var(&opt.backupOptions.RetentionPolicy.KeepWeekly, "retention-keep-weekly", opt.backupOptions.RetentionPolicy.KeepWeekly, "Specify value for retention strategy")
	cmd.Flags().Int64Var(&opt.backupOptions.RetentionPolicy.KeepMonthly, "retention-keep-monthly", opt.backupOptions.RetentionPolicy.KeepMonthly, "Specify value for retention strategy")
	cmd.Flags().Int64Var(&opt.backupOptions.RetentionPolicy.KeepYearly, "retention-keep-yearly", opt.backupOptions.RetentionPolicy.KeepYearly, "Specify value for retention strategy")
	cmd.Flags().StringSliceVar(&opt.backupOptions.RetentionPolicy.KeepTags, "retention-keep-tags", opt.backupOptions.RetentionPolicy.KeepTags, "Specify value for retention strategy")
	cmd.Flags().BoolVar(&opt.backupOptions.RetentionPolicy.Prune, "retention-prune", opt.backupOptions.RetentionPolicy.Prune, "Specify whether to prune old snapshot data")
	cmd.Flags().BoolVar(&opt.backupOptions.RetentionPolicy.DryRun, "retention-dry-run", opt.backupOptions.RetentionPolicy.DryRun, "Specify whether to test retention policy without deleting actual data")

	cmd.Flags().StringVar(&opt.outputDir, "output-dir", opt.outputDir, "Directory where output.json file will be written (keep empty if you don't need to write output in file)")

	return cmd
}

func (opt *esOptions) backupElasticsearch(targetRef api_v1beta1.TargetRef) (*restic.BackupOutput, error) {
	// if any pre-backup actions has been assigned to it, execute them
	actionOptions := api_util.ActionOptions{
		StashClient:       opt.stashClient,
		TargetRef:         targetRef,
		SetupOptions:      opt.setupOptions,
		BackupSessionName: opt.backupSessionName,
		Namespace:         opt.namespace,
	}
	err := api_util.ExecutePreBackupActions(actionOptions)
	if err != nil {
		return nil, err
	}
	// wait until the backend repository has been initialized.
	err = api_util.WaitForBackendRepository(actionOptions)
	if err != nil {
		return nil, err
	}
	// apply nice, ionice settings from env
	opt.setupOptions.Nice, err = v1.NiceSettingsFromEnv()
	if err != nil {
		return nil, err
	}
	opt.setupOptions.IONice, err = v1.IONiceSettingsFromEnv()
	if err != nil {
		return nil, err
	}

	// get app binding
	appBinding, err := opt.catalogClient.AppcatalogV1alpha1().AppBindings(opt.namespace).Get(context.TODO(), opt.appBindingName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	// get secret
	appBindingSecret, err := opt.kubeClient.CoreV1().Secrets(opt.namespace).Get(context.TODO(), appBinding.Spec.Secret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// transform secret
	err = appBinding.TransformSecret(opt.kubeClient, appBindingSecret.Data)
	if err != nil {
		return nil, err
	}

	// clear directory before running multielasticdump
	klog.Infoln("Cleaning up directory: ", opt.interimDataDir)
	if err := clearDir(opt.interimDataDir); err != nil {
		return nil, err
	}

	// write the credential ifo into a file
	// TODO: support backup without authentication
	httpAuthFilePath := filepath.Join(opt.setupOptions.ScratchDir, ESAuthFile)
	err = writeAuthFile(httpAuthFilePath, appBindingSecret)
	if err != nil {
		return nil, err
	}

	var tlsArgs string
	if appBinding.Spec.ClientConfig.CABundle != nil {
		if err := ioutil.WriteFile(filepath.Join(opt.setupOptions.ScratchDir, ESCACertFile), appBinding.Spec.ClientConfig.CABundle, os.ModePerm); err != nil {
			return nil, err
		}
		tlsArgs = fmt.Sprintf("--ca-input=%v", filepath.Join(opt.setupOptions.ScratchDir, ESCACertFile))
	}

	appSVC := appBinding.Spec.ClientConfig.Service
	esURL := fmt.Sprintf("%v://%s:%d", appSVC.Scheme, appSVC.Name, appSVC.Port)

	// wait for DB ready
	waitForDBReady(appBinding.Spec.ClientConfig.Service.Name, appBinding.Spec.ClientConfig.Service.Port, opt.waitTimeout)

	// run separate shell to dump indices
	klog.Infoln("Performing multielasticdump on ", appSVC.Name)
	esShell := sh.NewSession()
	esShell.ShowCMD = false
	esShell.SetEnv("NODE_TLS_REJECT_UNAUTHORIZED", "0") //xref: https://github.com/taskrabbit/elasticsearch-dump#bypassing-self-sign-certificate-errors

	args := []interface{}{
		fmt.Sprintf(`--input=%v`, esURL),
		fmt.Sprintf(`--output=%v`, opt.interimDataDir),
		fmt.Sprintf("--httpAuthFile=%s", httpAuthFilePath),
		tlsArgs,
	}
	for _, arg := range strings.Fields(opt.esArgs) {
		args = append(args, arg)
	}

	esShell.Command(MultiElasticDumpCMD, args...) // xref: multielasticdump: https://github.com/taskrabbit/elasticsearch-dump#multielasticdump

	if err := esShell.Run(); err != nil {
		return nil, err
	}

	// dumped data has been stored in the interim data dir. Now, we will backup this directory using Stash.
	opt.backupOptions.BackupPaths = []string{opt.interimDataDir}

	// init restic wrapper
	resticWrapper, err := restic.NewResticWrapper(opt.setupOptions)
	if err != nil {
		return nil, err
	}

	// Run backup
	return resticWrapper.RunBackup(opt.backupOptions, targetRef)
}
