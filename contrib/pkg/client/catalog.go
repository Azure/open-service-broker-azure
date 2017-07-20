package client

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-service-broker/pkg/service"
)

// GetCatalog retrieves the catalog from the broker specoified by host name
// and port number
func GetCatalog(host string, port int) (service.Catalog, error) {
	url := fmt.Sprintf("%s/v2/catalog", getBaseURL(host, port))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error building request: %s", err)
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting catalog: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unanticipated http response code %d",
			resp.StatusCode,
		)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}
	catalog, err := service.NewCatalogFromJSONString(string(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error decoding catalog: %s", err)
	}
	return catalog, nil
}
