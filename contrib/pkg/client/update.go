package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Update initiates updating of an existing service instance
func Update(
	host string,
	port int,
	username string,
	password string,
	instanceID string,
	serviceID string,
	planID string,
	params map[string]interface{},
) error {
	url := fmt.Sprintf(
		"%s/v2/service_instances/%s",
		getBaseURL(host, port),
		instanceID,
	)
	updatingRequest := UpdatingRequest{
		ServiceID:  serviceID,
		PlanID:     planID,
		Parameters: params,
	}
	jsonBytes, err := json.Marshal(updatingRequest)
	if err != nil {
		return fmt.Errorf("error encoding request body: %s", err)
	}
	req, err := newRequest(http.MethodPatch, url, username, password, jsonBytes)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("accepts_incomplete", "true")
	req.URL.RawQuery = q.Encode()
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing update call: %s", err)
	}
	defer resp.Body.Close() // nolint: errcheck
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf(
			"unanticipated http response code %d",
			resp.StatusCode,
		)
	}
	return nil
}
