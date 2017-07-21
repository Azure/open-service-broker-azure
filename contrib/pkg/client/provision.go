package client

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Azure/azure-service-broker/pkg/service"
	uuid "github.com/satori/go.uuid"
)

// Provision initiates provisioning of a new service instance
func Provision(
	host string,
	port int,
	serviceID string,
	planID string,
	params map[string]string,
) (string, error) {
	instanceID := uuid.NewV4().String()
	url := fmt.Sprintf(
		"%s/v2/service_instances/%s",
		getBaseURL(host, port),
		instanceID,
	)
	provisioningRequest := &service.ProvisioningRequest{
		ServiceID:  serviceID,
		PlanID:     planID,
		Parameters: params,
	}
	jsonStr, err := provisioningRequest.ToJSONString()
	if err != nil {
		return "", fmt.Errorf("error encoding request body: %s", err)
	}
	req, err := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		return "", fmt.Errorf("error building request: %s", err)
	}
	q := req.URL.Query()
	q.Add("accepts_incomplete", "true")
	req.URL.RawQuery = q.Encode()
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error requesting catalog: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf(
			"unanticipated http response code %d",
			resp.StatusCode,
		)
	}
	return instanceID, nil
}
