package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/open-service-broker-azure/pkg/api"
)

// Poll polls the status of an instance
func Poll(
	host string,
	port int,
	username string,
	password string,
	instanceID string,
	operation string,
) (string, error) {
	url := fmt.Sprintf(
		"%s/v2/service_instances/%s/last_operation",
		getBaseURL(host, port),
		instanceID,
	)
	req, err := newRequest(http.MethodGet, url, username, password, nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add("operation", operation)
	req.URL.RawQuery = q.Encode()
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing polling call: %s", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %s", err)
	}
	defer resp.Body.Close() // nolint: errcheck
	if operation == api.OperationDeprovisioning &&
		resp.StatusCode == http.StatusGone {
		return "gone", nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"unanticipated http response code %d",
			resp.StatusCode,
		)
	}
	responseMap := make(map[string]string)
	if err := json.Unmarshal(bodyBytes, &responseMap); err != nil {
		return "", fmt.Errorf("error unmarshaling response body: %s", err)
	}
	state, ok := responseMap["state"]
	if !ok {
		return "", errors.New("polling response did not include operation state")
	}
	return state, nil
}
