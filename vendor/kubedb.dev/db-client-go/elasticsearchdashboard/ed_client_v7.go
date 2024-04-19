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
	"fmt"
	"io"
	"net/http"
	"strings"

	esapi "kubedb.dev/apimachinery/apis/elasticsearch/v1alpha1"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

type EDClientV7 struct {
	Client *resty.Client
	Config *Config
}

func (h *EDClientV7) GetHealthStatus() (*Health, error) {
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
func (h *EDClientV7) GetStateFromHealthResponse(health *Health) (esapi.DashboardServerState, error) {
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
		return "", errors.Wrap(err, "Failed to parse response body")
	}

	if overallStatus, ok := responseBody.Status["overall"].(map[string]interface{}); ok {
		if overallState, ok := overallStatus["state"].(string); ok {
			health.OverallState = overallState
		} else {
			return "", errors.New("Failed to parse overallState")
		}
	} else {
		return "", errors.New("Failed to parse overallStatus")
	}

	// get the statuses for plugins stored,
	// so that the plugins which are not available or ready can be shown from condition message
	if statuses, ok := responseBody.Status["statuses"].([]interface{}); ok {
		for _, sts := range statuses {
			if curr, ok := sts.(map[string]interface{}); ok {
				if curr["state"].(string) != string(esapi.StateGreen) {
					health.StateFailedReason[curr["id"].(string)] = strings.Join([]string{curr["state"].(string), curr["message"].(string)}, ",")
				}
			} else {
				return "", errors.New("Failed to convert statuses to map[string]interface{}")
			}
		}
	} else {
		return "", errors.New("Failed to convert statuses to []interface{}")
	}

	return esapi.DashboardServerState(health.OverallState), nil
}

func (h *EDClientV7) ExportSavedObjects(spaceName string) (*Response, error) {
	req := h.Client.R().
		SetDoNotParseResponse(true).
		SetHeaders(jsonHeaderForKibanaAPI).
		SetBody([]byte(SavedObjectsReqBodyES))
	res, err := req.Post("/s/" + spaceName + SavedObjectsExportURL)
	if err != nil {
		klog.Error(err, "Failed to send http request")
		return nil, err
	}

	return &Response{
		Code: res.StatusCode(),
		Body: res.RawBody(),
	}, nil
}

func (h *EDClientV7) ImportSavedObjects(spaceName, filepath string) (*Response, error) {
	req := h.Client.R().
		SetDoNotParseResponse(true).
		SetHeader("kbn-xsrf", "true").
		SetFile("file", filepath).
		SetQueryParam("overwrite", "true")
	res, err := req.Post("/s/" + spaceName + SavedObjectsImportURL)
	if err != nil {
		klog.Error(err, "Failed to send http request")
		return nil, err
	}

	return &Response{
		Code: res.StatusCode(),
		Body: res.RawBody(),
	}, nil
}

func (h *EDClientV7) ListSpaces() ([]Space, error) {
	req := h.Client.R().
		SetDoNotParseResponse(true).
		SetHeaders(jsonHeaderForKibanaAPI)
	res, err := req.Get(SpacesURL)
	if err != nil {
		klog.Error("Failed to send http request")
		return nil, err
	}

	body, err := io.ReadAll(res.RawBody())
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to list dashboard spaces %s", string(body))
	}

	var spaces []Space
	if err = json.Unmarshal(body, &spaces); err != nil {
		return nil, err
	}

	return spaces, nil
}

func (h *EDClientV7) CreateSpace(space Space) error {
	req := h.Client.R().
		SetDoNotParseResponse(true).
		SetHeaders(jsonHeaderForKibanaAPI).
		SetBody(space)
	res, err := req.Post(SpacesURL)
	if err != nil {
		klog.Error(err, "Failed to send http request")
		return err
	}

	body, err := io.ReadAll(res.RawBody())
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to create dashboard space %s: %s", space.Name, string(body))
	}

	return nil
}
