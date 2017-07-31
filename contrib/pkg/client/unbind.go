package client

import (
	"fmt"
	"net/http"
)

// Unbind carries out unbinding
func Unbind(
	host string,
	port int,
	username string,
	password string,
	instanceID string,
	bindingID string,
) error {
	url := fmt.Sprintf(
		"%s/v2/service_instances/%s/service_bindings/%s",
		getBaseURL(host, port),
		instanceID,
		bindingID,
	)
	req, err := http.NewRequest(
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error building request: %s", err)
	}
	if username != "" || password != "" {
		addAuthHeader(req, username, password)
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing unbind call: %s", err)
	}
	defer resp.Body.Close() // nolint: errcheck
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"unanticipated http response code %d",
			resp.StatusCode,
		)
	}
	return nil
}
