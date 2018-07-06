package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetCatalog retrieves the catalog from the broker specified by host name
// and port number
func GetCatalog(
	useSSL bool,
	skipCertValidation bool,
	host string,
	port int,
	username string,
	password string,
) (Catalog, error) {
	catalog := Catalog{}
	url := fmt.Sprintf("%s/v2/catalog", getBaseURL(useSSL, host, port))
	req, err := newRequest(http.MethodGet, url, username, password, nil)
	if err != nil {
		return catalog, err
	}
	httpClient := getHTTPClient(skipCertValidation)
	resp, err := httpClient.Do(req)
	if err != nil {
		return catalog, fmt.Errorf("error requesting catalog: %s", err)
	}
	defer resp.Body.Close() // nolint: errcheck
	if resp.StatusCode != http.StatusOK {
		return catalog, fmt.Errorf(
			"unanticipated http response code %d",
			resp.StatusCode,
		)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return catalog, fmt.Errorf("error reading response body: %s", err)
	}
	if err := json.Unmarshal(bodyBytes, &catalog); err != nil {
		return catalog, fmt.Errorf("error decoding catalog: %s", err)
	}
	return catalog, nil
}
