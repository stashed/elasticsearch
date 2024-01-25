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

package elasticsearchdashboard

import (
	"encoding/json"
	"io"
	"strings"

	esapi "kubedb.dev/apimachinery/apis/elasticsearch/v1alpha1"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

type EDClientV8 struct {
	Client *resty.Client
	Config *Config
}

func (h *EDClientV8) GetHealthStatus() (*Health, error) {
	req := h.Client.R().SetDoNotParseResponse(true)
	res, err := req.Get(h.Config.api)
	if err != nil {
		klog.Error(err, "Failed to send http request")
		return nil, err
	}

	statesList := make(map[string]string)

	healthStatus := &Health{
		ConnectionResponse: Response{
			Code:   res.StatusCode(),
			header: res.Header(),
			Body:   res.RawBody(),
		},
		StateFailedReason: statesList,
	}

	return healthStatus, nil
}

// GetStateFromHealthResponse parse health response in json from server and
// return overall status of the server
func (h *EDClientV8) GetStateFromHealthResponse(health *Health) (esapi.DashboardServerState, error) {
	resStatus := health.ConnectionResponse

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			err1 := errors.Wrap(err, "failed to parse response body")
			if err1 != nil {
				return
			}
			return
		}
	}(resStatus.Body)

	var responseBody ResponseBody
	body, _ := io.ReadAll(resStatus.Body)
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse response body")
	}

	if overallStatus, ok := responseBody.Status["overall"].(map[string]interface{}); ok {
		if overallState, ok := overallStatus["level"].(string); ok {
			health.OverallState = overallState
		} else {
			return "", errors.New("Failed to parse overallState")
		}
	} else {
		return "", errors.New("Failed to parse overallStatus")
	}

	if pluginsStatus, ok := responseBody.Status["plugins"].(map[string]interface{}); ok {
		for plugin, pluginStatus := range pluginsStatus {
			if pstatus, ok := pluginStatus.(map[string]interface{}); ok {
				if pstatus["level"].(string) != string(esapi.StateAvailable) {
					health.StateFailedReason[plugin] = strings.Join([]string{pstatus["level"].(string), pstatus["summary"].(string)}, ",")
				}
			} else {
				return "", errors.New("Failed to parse plugin status")
			}
		}
	} else {
		return "", errors.New("Failed to parse overallStatus")
	}

	return esapi.DashboardServerState(health.OverallState), nil
}

func (h *EDClientV8) ExportSavedObjects() (*Response, error) {
	req := h.Client.R().
		SetDoNotParseResponse(true).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"kbn-xsrf":     "true",
		}).
		SetBody([]byte(SavedObjectsReqBody))
	res, err := req.Post(SavedObjectsExportURL)
	if err != nil {
		klog.Error(err, "Failed to send http request")
		return nil, err
	}

	return &Response{
		Code: res.StatusCode(),
		Body: res.RawBody(),
	}, nil
}

func (h *EDClientV8) ImportSavedObjects(filepath string) (*Response, error) {
	req := h.Client.R().
		SetDoNotParseResponse(true).
		SetHeader("kbn-xsrf", "true").
		SetFile("file", filepath).
		SetQueryParam("overwrite", "true")
	res, err := req.Post(SavedObjectsImportURL)
	if err != nil {
		klog.Error(err, "Failed to send http request")
		return nil, err
	}

	return &Response{
		Code: res.StatusCode(),
		Body: res.RawBody(),
	}, nil
}
