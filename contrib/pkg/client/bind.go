package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-service-broker/pkg/service"
	uuid "github.com/satori/go.uuid"
)

// Bind carries out binding to an existing service
func Bind(
	host string,
	port int,
	username string,
	password string,
	instanceID string,
	params map[string]string,
) (string, map[string]interface{}, error) {
	bindingID := uuid.NewV4().String()
	url := fmt.Sprintf(
		"%s/v2/service_instances/%s/service_bindings/%s",
		getBaseURL(host, port),
		instanceID,
		bindingID,
	)
	bindingRequest := &service.BindingRequest{
		Parameters: params,
	}
	json, err := bindingRequest.ToJSON()
	if err != nil {
		return "", nil, fmt.Errorf("error encoding request body: %s", err)
	}
	req, err := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer(json),
	)
	if err != nil {
		return "", nil, fmt.Errorf("error building request: %s", err)
	}
	if username != "" || password != "" {
		addAuthHeader(req, username, password)
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("error executing bind call: %s", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("error reading response body: %s", err)
	}
	defer resp.Body.Close() // nolint: errcheck
	if resp.StatusCode != http.StatusCreated {
		return "", nil, fmt.Errorf(
			"unanticipated http response code %d",
			resp.StatusCode,
		)
	}
	bindingResponse := &service.BindingResponse{}
	err = service.GetBindingResponseFromJSON(bodyBytes, bindingResponse)
	if err != nil {
		return "", nil, fmt.Errorf("error decoding response body: %s", err)
	}
	credsMap, ok := bindingResponse.Credentials.(map[string]interface{})
	if !ok {
		return "", nil, fmt.Errorf("error decoding response body: %s", err)
	}
	return bindingID, credsMap, nil
}
