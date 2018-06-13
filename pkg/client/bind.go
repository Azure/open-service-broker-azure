package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// Bind carries out binding to an existing service
func Bind(
	host string,
	port int,
	username string,
	password string,
	instanceID string,
	params map[string]interface{},
) (string, map[string]interface{}, error) {
	bindingID := uuid.NewV4().String()
	url := fmt.Sprintf(
		"%s/v2/service_instances/%s/service_bindings/%s",
		getBaseURL(host, port),
		instanceID,
		bindingID,
	)
	bindingRequest := BindingRequest{
		Parameters: params,
	}
	jsonBytes, err := json.Marshal(bindingRequest)
	if err != nil {
		return "", nil, fmt.Errorf("error encoding request body: %s", err)
	}
	req, err := newRequest(http.MethodPut, url, username, password, jsonBytes)
	if err != nil {
		return "", nil, err
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
	bindingResponse := BindingResponse{}
	if err := json.Unmarshal(bodyBytes, &bindingResponse); err != nil {
		return "", nil, fmt.Errorf("error unmarshaling response body: %s", err)
	}
	return bindingID, bindingResponse.Credentials, nil
}
