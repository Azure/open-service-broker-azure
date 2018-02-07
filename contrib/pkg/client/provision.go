package client

import (
	"fmt"
	"net/http"

	"github.com/Azure/open-service-broker-azure/pkg/api"
	uuid "github.com/satori/go.uuid"
)

// Provision initiates provisioning of a new service instance
func Provision(
	host string,
	port int,
	username string,
	password string,
	serviceID string,
	planID string,
	params map[string]interface{},
	tags map[string]string,
) (string, error) {
	instanceID := uuid.NewV4().String()
	url := fmt.Sprintf(
		"%s/v2/service_instances/%s",
		getBaseURL(host, port),
		instanceID,
	)
	params["tags"] = tags
	provisioningRequest := &api.ProvisioningRequest{
		ServiceID:  serviceID,
		PlanID:     planID,
		Parameters: params,
	}
	json, err := provisioningRequest.ToJSON()
	if err != nil {
		return "", fmt.Errorf("error encoding request body: %s", err)
	}
	req, err := newRequest(http.MethodPut, url, username, password, json)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add("accepts_incomplete", "true")
	req.URL.RawQuery = q.Encode()
	httpClient := &http.Client{}
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
