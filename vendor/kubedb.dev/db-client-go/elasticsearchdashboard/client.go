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
	"io"
	"net/http"

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	esapi "kubedb.dev/apimachinery/apis/elasticsearch/v1alpha1"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	core "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	SavedObjectsReqBodyES = `{"type": ["dashboard", "config", "index-pattern", "url", "query", "tag", "canvas-element", "canvas-workpad", "action", "alert", "visualization",
"graph-workspace", "map", "lens", "cases", "search", "osquery-saved-query", "osquery-pack", "uptime-dynamic-settings", "infrastructure-ui-source", "metrics-explorer-view",
"inventory-view", "apm-indices"]}`
	SavedObjectsReqBodyOS = `{"type": ["config", "url", "index-pattern", "query", "dashboard", "visualization", "visualization-visbuilder", "map", "observability-panel",
"observability-visualization", "search"]}`
	SavedObjectsExportURL = "/api/saved_objects/_export"
	SavedObjectsImportURL = "/api/saved_objects/_import"
	SpacesURL             = "/api/spaces/space"
)

var jsonHeaderForKibanaAPI = map[string]string{
	"Content-Type": "application/json",
	"kbn-xsrf":     "true",
}

type Client struct {
	EDClient
}

type ClientOptions struct {
	KClient   client.Client
	Dashboard *esapi.ElasticsearchDashboard
	ESVersion *catalog.ElasticsearchVersion
	DB        *dbapi.Elasticsearch
	Ctx       context.Context
	Secret    *core.Secret
}

type DbVersionInfo struct {
	Name       string
	Version    string
	AuthPlugin catalog.ElasticsearchAuthPlugin
}

type Config struct {
	host             string
	api              string
	username         string
	password         string
	connectionScheme string
	transport        *http.Transport
	dbVersionInfo    *DbVersionInfo
}

type Health struct {
	ConnectionResponse Response
	OverallState       string
	StateFailedReason  map[string]string
}

type Response struct {
	Code   int
	header http.Header
	Body   io.ReadCloser
}

type ResponseBody struct {
	Name    string                 `json:"name"`
	UUID    string                 `json:"uuid"`
	Version map[string]interface{} `json:"version"`
	Status  map[string]interface{} `json:"status"`
	Metrics map[string]interface{} `json:"metrics"`
}

type Space struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	Color            string   `json:"color,omitempty"`
	Initials         string   `json:"initials,omitempty"`
	DisabledFeatures []string `json:"disabledFeatures,omitempty"`
	ImageUrl         string   `json:"imageUrl,omitempty"`
}
