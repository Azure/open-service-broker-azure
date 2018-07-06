package client

import (
	"fmt"
	"net/http"
)

// Deprovision initiates deprovisioning of a service instance
func Deprovision(
	useSSL bool,
	skipCertValidation bool,
	host string,
	port int,
	username string,
	password string,
	instanceID string,
) error {
	url := fmt.Sprintf(
		"%s/v2/service_instances/%s",
		getBaseURL(useSSL, host, port),
		instanceID,
	)
	req, err := newRequest(http.MethodDelete, url, username, password, nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("accepts_incomplete", "true")
	req.URL.RawQuery = q.Encode()
	httpClient := getHTTPClient(skipCertValidation)
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing deprovision call: %s", err)
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
