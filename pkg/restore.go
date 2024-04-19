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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"stash.appscode.dev/apimachinery/apis"
	api_v1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	"stash.appscode.dev/apimachinery/pkg/restic"

	"github.com/spf13/cobra"
	license "go.bytebuilders.dev/license-verifier/kubernetes"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gomodules.xyz/flags"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	appcatalog "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	appcatalog_cs "kmodules.xyz/custom-resources/client/clientset/versioned"
	v1 "kmodules.xyz/offshoot-api/api/v1"
	"kubedb.dev/db-client-go/elasticsearchdashboard"
)

func NewCmdRestore() *cobra.Command {
	var (
		masterURL      string
		kubeconfigPath string
		opt            = esOptions{
			waitTimeout: 300,
			setupOptions: restic.SetupOptions{
				ScratchDir:  restic.DefaultScratchDir,
				EnableCache: false,
			},
			restoreOptions: restic.RestoreOptions{
				Host: restic.DefaultHost,
			},
		}
	)

	cmd := &cobra.Command{
		Use:               "restore-es",
		Short:             "Restores Elasticsearch DB Backup",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.EnsureRequiredFlags(cmd, "appbinding", "provider", "storage-secret-name", "storage-secret-namespace")

			// prepare client
			config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
			if err != nil {
				return err
			}
			opt.config = config

			opt.kubeClient, err = kubernetes.NewForConfig(config)
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
				Namespace:  opt.appBindingNamespace,
			}
			var restoreOutput *restic.RestoreOutput
			restoreOutput, err = opt.restoreElasticsearch(targetRef)
			if err != nil {
				restoreOutput = &restic.RestoreOutput{
					RestoreTargetStatus: api_v1beta1.RestoreMemberStatus{
						Ref: targetRef,
						Stats: []api_v1beta1.HostRestoreStats{
							{
								Hostname: opt.restoreOptions.Host,
								Phase:    api_v1beta1.HostRestoreFailed,
								Error:    err.Error(),
							},
						},
					},
				}
			}
			// If output directory specified, then write the output in "output.json" file in the specified directory
			if opt.outputDir != "" {
				return restoreOutput.WriteOutput(filepath.Join(opt.outputDir, restic.DefaultOutputFileName))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&opt.esArgs, "es-args", opt.esArgs, "Additional arguments")
	cmd.Flags().Int32Var(&opt.waitTimeout, "wait-timeout", opt.waitTimeout, "Number of seconds to wait for the database to be ready")

	cmd.Flags().StringVar(&masterURL, "master", masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&opt.namespace, "namespace", "default", "Namespace of Backup/Restore Session")
	cmd.Flags().StringVar(&opt.appBindingName, "appbinding", opt.appBindingName, "Name of the app binding")
	cmd.Flags().StringVar(&opt.appBindingNamespace, "appbinding-namespace", opt.appBindingNamespace, "Namespace of the app binding")
	cmd.Flags().StringVar(&opt.storageSecret.Name, "storage-secret-name", opt.storageSecret.Name, "Name of the storage secret")
	cmd.Flags().StringVar(&opt.storageSecret.Namespace, "storage-secret-namespace", opt.storageSecret.Namespace, "Namespace of the storage secret")

	cmd.Flags().StringVar(&opt.setupOptions.Provider, "provider", opt.setupOptions.Provider, "Backend provider (i.e. gcs, s3, azure etc)")
	cmd.Flags().StringVar(&opt.setupOptions.Bucket, "bucket", opt.setupOptions.Bucket, "Name of the cloud bucket/container (keep empty for local backend)")
	cmd.Flags().StringVar(&opt.setupOptions.Endpoint, "endpoint", opt.setupOptions.Endpoint, "Endpoint for s3/s3 compatible backend or REST server URL")
	cmd.Flags().BoolVar(&opt.setupOptions.InsecureTLS, "insecure-tls", opt.setupOptions.InsecureTLS, "InsecureTLS for TLS secure s3/s3 compatible backend")
	cmd.Flags().StringVar(&opt.setupOptions.Region, "region", opt.setupOptions.Region, "Region for s3/s3 compatible backend")
	cmd.Flags().StringVar(&opt.setupOptions.Path, "path", opt.setupOptions.Path, "Directory inside the bucket where backup will be stored")
	cmd.Flags().StringVar(&opt.setupOptions.ScratchDir, "scratch-dir", opt.setupOptions.ScratchDir, "Temporary directory")
	cmd.Flags().BoolVar(&opt.setupOptions.EnableCache, "enable-cache", opt.setupOptions.EnableCache, "Specify whether to enable caching for restic")
	cmd.Flags().Int64Var(&opt.setupOptions.MaxConnections, "max-connections", opt.setupOptions.MaxConnections, "Specify maximum concurrent connections for GCS, Azure and B2 backend")

	cmd.Flags().StringVar(&opt.restoreOptions.Host, "hostname", opt.restoreOptions.Host, "Name of the host machine")
	cmd.Flags().StringVar(&opt.restoreOptions.SourceHost, "source-hostname", opt.restoreOptions.SourceHost, "Name of the host whose data will be restored")
	cmd.Flags().StringSliceVar(&opt.restoreOptions.Snapshots, "snapshot", opt.restoreOptions.Snapshots, "Snapshots to restore")
	cmd.Flags().StringVar(&opt.interimDataDir, "interim-data-dir", opt.interimDataDir, "Directory where the restored data will be stored temporarily before injecting into the desired database.")
	cmd.Flags().BoolVar(&opt.enableDashboard, "enable-dashboard-restore", opt.enableDashboard, "Specify whether to enable kibana dashboard restore")

	cmd.Flags().StringVar(&opt.outputDir, "output-dir", opt.outputDir, "Directory where output.json file will be written (keep empty if you don't need to write output in file)")

	return cmd
}

func (opt *esOptions) restoreElasticsearch(targetRef api_v1beta1.TargetRef) (*restic.RestoreOutput, error) {
	var err error
	err = license.CheckLicenseEndpoint(opt.config, licenseApiService, SupportedProducts)
	if err != nil {
		return nil, err
	}

	opt.setupOptions.StorageSecret, err = opt.kubeClient.CoreV1().Secrets(opt.storageSecret.Namespace).Get(context.TODO(), opt.storageSecret.Name, metav1.GetOptions{})
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

	appBinding, err := opt.catalogClient.AppcatalogV1alpha1().AppBindings(opt.appBindingNamespace).Get(context.TODO(), opt.appBindingName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// clear directory before running multielasticdump
	session := opt.newSessionWrapper(MultiElasticDumpCMD)

	err = opt.setDatabaseCredentials(appBinding, session.cmd)
	if err != nil {
		return nil, err
	}

	// clear directory before running multielasticdump
	klog.Infoln("Cleaning up directory: ", opt.interimDataDir)
	if err := clearDir(opt.interimDataDir); err != nil {
		return nil, err
	}

	url, err := appBinding.URL()
	if err != nil {
		return nil, err
	}

	session.cmd.Args = append(session.cmd.Args, []interface{}{
		"--direction=load",
		fmt.Sprintf(`--input=%v`, opt.interimDataDir),
		fmt.Sprintf(`--output=%v`, url),
	}...)

	err = session.setTLSParameters(appBinding, opt.setupOptions.ScratchDir)
	if err != nil {
		return nil, err
	}

	err = opt.waitForDBReady(appBinding)
	if err != nil {
		return nil, err
	}

	// we will restore the desired data into interim data dir before injecting into the desired database
	opt.restoreOptions.RestorePaths = []string{opt.interimDataDir}

	// init restic wrapper
	resticWrapper, err := restic.NewResticWrapper(opt.setupOptions)
	if err != nil {
		return nil, err
	}

	restoreOutput, err := resticWrapper.RunRestore(opt.restoreOptions, targetRef)
	if err != nil {
		return nil, err
	}

	// delete the metadata file as it is not required for restoring the dumps
	if err := clearFile(filepath.Join(opt.interimDataDir, apis.ESMetaFile)); err != nil {
		return nil, err
	}

	if err = opt.restoreDashboardObjects(appBinding); err != nil {
		return nil, fmt.Errorf("failed to restore dashboard objects %w", err)
	}

	// run separate shell to restore indices
	// klog.Infoln("Performing multielasticdump on", hostname)
	session.sh.ShowCMD = false
	session.setUserArgs(opt.esArgs)
	session.sh.Command(session.cmd.Name, session.cmd.Args...) // xref: multielasticdump: https://github.com/taskrabbit/elasticsearch-dump#multielasticdump

	if err := session.sh.Run(); err != nil {
		return nil, err
	}
	return restoreOutput, nil
}

func clearFile(filepath string) error {
	if _, err := os.Stat(filepath); err == nil {
		if err := os.Remove(filepath); err != nil {
			return fmt.Errorf("unable to clean file: %v. Reason: %v", filepath, err)
		}
	}
	return nil
}

func (opt *esOptions) restoreDashboardObjects(appBinding *appcatalog.AppBinding) error {
	if !opt.enableDashboard {
		return nil
	}

	dashboardClient, err := opt.getDashboardClient(appBinding)
	if err != nil {
		return err
	}

	spaces, err := opt.getSpaces()
	if err != nil {
		return err
	}

	existingSpaces, err := dashboardClient.ListSpaces()
	if err != nil {
		return err
	}

	for _, space := range spaces {
		if !isExist(existingSpaces, space.Id) {
			if err = dashboardClient.CreateSpace(space); err != nil {
				return fmt.Errorf("failed to create space %s: %w", space.Id, err)
			}
		}

		response, err := dashboardClient.ImportSavedObjects(space.Id, opt.getDashboardFilePath(space.Id))
		if err != nil {
			return err
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if response.Code != http.StatusOK {
			return fmt.Errorf("failed to import dashboard saved objects %s", string(body))
		}

		// delete the dashboard file(s) as it is not required for restoring the dumps
		if err = clearFile(opt.getDashboardFilePath(space.Id)); err != nil {
			return err
		}
	}

	return nil
}

func (opt *esOptions) getSpaces() ([]elasticsearchdashboard.Space, error) {
	if _, err := os.Stat(filepath.Join(opt.interimDataDir, SpacesInfoFile)); os.IsNotExist(err) {
		spaces := make([]elasticsearchdashboard.Space, 0)
		if err = filepath.Walk(opt.interimDataDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(info.Name()) == DashboardSavedObjectsExt {
				id, _ := strings.CutSuffix(info.Name(), DashboardSavedObjectsExt)
				spaces = append(spaces, elasticsearchdashboard.Space{
					Id:   id,
					Name: cases.Title(language.English).String(strings.Replace(id, "-", " ", -1)),
				})
			}

			return nil
		}); err != nil {
			return nil, err
		}

		if len(spaces) == 0 {
			return nil, fmt.Errorf("no spaces found in interim data directory")
		}

		return spaces, nil
	} else {
		data, err := os.ReadFile(filepath.Join(opt.interimDataDir, SpacesInfoFile))
		if err != nil {
			return nil, fmt.Errorf("failed to read spaces info %w", err)
		}

		var spaces []elasticsearchdashboard.Space
		if err = json.Unmarshal(data, &spaces); err != nil {
			return nil, fmt.Errorf("failed to unmarshal spaces info %w", err)
		}

		// delete the spaces info file as it is not required for restoring the dumps
		if err = clearFile(filepath.Join(opt.interimDataDir, SpacesInfoFile)); err != nil {
			return nil, err
		}

		return spaces, nil
	}
}

func isExist(existingSpaces []elasticsearchdashboard.Space, spaceId string) bool {
	for _, space := range existingSpaces {
		if space.Id == spaceId {
			return true
		}
	}
	return false
}
