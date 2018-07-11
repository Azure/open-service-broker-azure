package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// Provision initiates provisioning of a new service instance
func Provision(
	useSSL bool,
	skipCertValidation bool,
	host string,
	port int,
	username string,
	password string,
	serviceID string,
	planID string,
	params map[string]interface{},
) (string, error) {
	instanceID := uuid.NewV4().String()
	url := fmt.Sprintf(
		"%s/v2/service_instances/%s",
		getBaseURL(useSSL, host, port),
		instanceID,
	)
	provisioningRequest := ProvisioningRequest{
		ServiceID:  serviceID,
		PlanID:     planID,
		Parameters: params,
	}
	jsonBytes, err := json.Marshal(provisioningRequest)
	if err != nil {
		return "", fmt.Errorf("error encoding request body: %s", err)
	}
	req, err := newRequest(http.MethodPut, url, username, password, jsonBytes)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add("accepts_incomplete", "true")
	req.URL.RawQuery = q.Encode()
	httpClient := getHTTPClient(skipCertValidation)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing provision call: %s", err)
	}
	defer resp.Body.Close() // nolint: errcheck
	if resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf(
			"unanticipated http response code %d",
			resp.StatusCode,
		)
	}
	return instanceID, nil
}
