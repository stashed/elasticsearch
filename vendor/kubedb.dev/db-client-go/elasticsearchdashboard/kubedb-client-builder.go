/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package elasticsearchdashboard

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"time"

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	esapi "kubedb.dev/apimachinery/apis/elasticsearch/v1alpha1"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	"github.com/Masterminds/semver/v3"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KubeDBClientBuilder struct {
	kc            client.Client
	dashboard     *esapi.ElasticsearchDashboard
	db            *dbapi.Elasticsearch
	dbVersion     *catalog.ElasticsearchVersion
	dbVersionInfo *DbVersionInfo
	authSecret    *core.Secret
	url           string
	podName       string
	ctx           context.Context
}

func NewKubeDBClientBuilder(kc client.Client, db *esapi.ElasticsearchDashboard) *KubeDBClientBuilder {
	return &KubeDBClientBuilder{
		kc:        kc,
		dashboard: db,
	}
}

func (o *KubeDBClientBuilder) WithPod(podName string) *KubeDBClientBuilder {
	o.podName = podName
	return o
}

func (o *KubeDBClientBuilder) WithURL(url string) *KubeDBClientBuilder {
	o.url = url
	return o
}

func (o *KubeDBClientBuilder) WithAuthSecret(secret *core.Secret) *KubeDBClientBuilder {
	o.authSecret = secret
	return o
}

func (o *KubeDBClientBuilder) WithDatabaseRef(db *dbapi.Elasticsearch) *KubeDBClientBuilder {
	o.db = db
	return o
}

func (o *KubeDBClientBuilder) WithDbVersion(version *catalog.ElasticsearchVersion) *KubeDBClientBuilder {
	o.dbVersion = version
	return o
}

func (o *KubeDBClientBuilder) WithDbVersionInfo(versionInfo *DbVersionInfo) *KubeDBClientBuilder {
	o.dbVersionInfo = versionInfo
	return o
}

func (o *KubeDBClientBuilder) WithContext(ctx context.Context) *KubeDBClientBuilder {
	o.ctx = ctx
	return o
}

func (o *KubeDBClientBuilder) GetElasticsearchDashboardClient() (*Client, error) {
	config := Config{
		host: getHostPath(o.dashboard),
		api:  esapi.KibanaStatusEndpoint,
		transport: &http.Transport{
			IdleConnTimeout: time.Second * 3,
			DialContext: (&net.Dialer{
				Timeout: time.Second * 30,
			}).DialContext,
		},
		connectionScheme: o.dashboard.GetConnectionScheme(),
		dbVersionInfo:    o.getDbVersionInfo(),
	}
	// If EnableSSL is true set tls config,
	// provide client certs and root CA
	if o.dashboard.Spec.EnableSSL {
		var certSecret core.Secret
		err := o.kc.Get(o.ctx, types.NamespacedName{
			Namespace: o.dashboard.Namespace,
			Name:      o.dashboard.CertificateSecretName(esapi.ElasticsearchDashboardServerCert),
		}, &certSecret)
		if err != nil {
			klog.Error(err, "failed to get serverCert secret")
			return nil, err
		}

		// get tls cert, clientCA and rootCA for tls config
		// use server cert ca for rootca as issuer ref is not taken into account
		clientCA := x509.NewCertPool()
		rootCA := x509.NewCertPool()

		crt, err := tls.X509KeyPair(certSecret.Data[core.TLSCertKey], certSecret.Data[core.TLSPrivateKeyKey])
		if err != nil {
			klog.Error(err, "failed to create certificate for TLS config")
			return nil, err
		}
		clientCA.AppendCertsFromPEM(certSecret.Data[esapi.CaCertKey])
		rootCA.AppendCertsFromPEM(certSecret.Data[esapi.CaCertKey])

		config.transport.TLSClientConfig = &tls.Config{
			Certificates: []tls.Certificate{crt},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    clientCA,
			RootCAs:      rootCA,
			MaxVersion:   tls.VersionTLS13,
		}
	}

	var username, password string

	// if security is enabled set database credentials in clientConfig
	if !o.db.Spec.DisableSecurity {

		if value, ok := o.authSecret.Data[core.BasicAuthUsernameKey]; ok {
			username = string(value)
		} else {
			klog.Info(fmt.Sprintf("Failed for secret: %s/%s, username is missing", o.authSecret.Namespace, o.authSecret.Name))
			return nil, errors.New("username is missing")
		}

		if value, ok := o.authSecret.Data[core.BasicAuthPasswordKey]; ok {
			password = string(value)
		} else {
			klog.Info(fmt.Sprintf("Failed for secret: %s/%s, password is missing", o.authSecret.Namespace, o.authSecret.Name))
			return nil, errors.New("password is missing")
		}

		config.username = username
		config.password = password
	}

	// parse version
	version, err := semver.NewVersion(config.dbVersionInfo.Version)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse version")
	}

	switch {
	// for Elasticsearch 7.x.x
	case config.dbVersionInfo.AuthPlugin == catalog.ElasticsearchAuthPluginXpack && version.Major() == 7:
		newClient := resty.New()
		newClient.SetTransport(config.transport).SetScheme(config.connectionScheme).SetBaseURL(config.host)
		newClient.SetHeader("Accept", "application/json")
		newClient.SetBasicAuth(config.username, config.password)
		newClient.SetTimeout(time.Second * 30)

		return &Client{
			&EDClientV7{
				Client: newClient,
				Config: &config,
			},
		}, nil

	case config.dbVersionInfo.AuthPlugin == catalog.ElasticsearchAuthPluginXpack && version.Major() == 8:
		newClient := resty.New()
		newClient.SetTransport(config.transport).SetScheme(config.connectionScheme).SetBaseURL(config.host)
		newClient.SetHeader("Accept", "application/json")
		newClient.SetBasicAuth(config.username, config.password)
		newClient.SetTimeout(time.Second * 30)

		return &Client{
			&EDClientV8{
				Client: newClient,
				Config: &config,
			},
		}, nil

	case config.dbVersionInfo.AuthPlugin == catalog.ElasticsearchAuthPluginOpenSearch:
		newClient := resty.New()
		newClient.SetTransport(config.transport).SetScheme(config.connectionScheme).SetBaseURL(config.host)
		newClient.SetHeader("Accept", "application/json")
		newClient.SetBasicAuth(config.username, config.password)
		newClient.SetTimeout(time.Second * 30)

		return &Client{
			&OSClient{
				Client: newClient,
				Config: &config,
			},
		}, nil
	}

	return nil, fmt.Errorf("unknown version: %s", config.dbVersionInfo.Name)
}

func (o *KubeDBClientBuilder) getDbVersionInfo() *DbVersionInfo {
	if o.dbVersionInfo == nil {
		return &DbVersionInfo{
			Name:       o.dbVersion.Name,
			Version:    o.dbVersion.Spec.Version,
			AuthPlugin: o.dbVersion.Spec.AuthPlugin,
		}
	}

	return o.dbVersionInfo
}

// return host path in
// format https://svc_name.namespace.svc:5601/api/status
func getHostPath(dashboard *esapi.ElasticsearchDashboard) string {
	return fmt.Sprintf("%v://%s.%s.svc:%d", dashboard.GetConnectionScheme(), dashboard.ServiceName(), dashboard.GetNamespace(), esapi.ElasticsearchDashboardRESTPort)
}
