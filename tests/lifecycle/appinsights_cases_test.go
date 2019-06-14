// +build !unit

package lifecycle

import (
	"fmt"
	"net/http"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
)

var appinsightsTestCases = []serviceLifecycleTestCase{
	{
		group:           "appinsights",
		name:            "asp-dot-net-web",
		serviceID:       "66130ee7-451b-4c61-8b78-d5c426a06f3e",
		planID:          "c14826d6-87c4-45de-94a0-52fad0893799",
		testCredentials: testAppinsightsCreds,
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"tags": map[string]interface{}{
				"test": "test",
			},
		},
	},
	{
		group:           "appinsights",
		name:            "java-web",
		serviceID:       "66130ee7-451b-4c61-8b78-d5c426a06f3e",
		planID:          "be62781f-5750-49c1-be8f-56e9c804c5fa",
		testCredentials: testAppinsightsCreds,
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"tags": map[string]interface{}{
				"test": "test",
			},
		},
	},
	{
		group:           "appinsights",
		name:            "node-dot-js",
		serviceID:       "66130ee7-451b-4c61-8b78-d5c426a06f3e",
		planID:          "7dc0b4ee-2322-4fec-88d1-1cce63e47fd8",
		testCredentials: testAppinsightsCreds,
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"tags": map[string]interface{}{
				"test": "test",
			},
		},
	},
	{
		group:           "appinsights",
		name:            "general",
		serviceID:       "66130ee7-451b-4c61-8b78-d5c426a06f3e",
		planID:          "a75e8854-591a-4ef2-b3f1-b311d2a02902",
		testCredentials: testAppinsightsCreds,
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"tags": map[string]interface{}{
				"test": "test",
			},
		},
	},
	{
		group:           "appinsights",
		name:            "app-center",
		serviceID:       "66130ee7-451b-4c61-8b78-d5c426a06f3e",
		planID:          "b730dc3e-6928-4d05-9193-edff71790095",
		testCredentials: testAppinsightsCreds,
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"tags": map[string]interface{}{
				"test": "test",
			},
		},
	},
}

func testAppinsightsCreds(credentials map[string]interface{}) error {
	ik, ok := credentials["instrumentationKey"].(string)
	if !ok {
		return fmt.Errorf(
			"can't find instrumentation key in the credentials",
		)
	}
	client := appinsights.NewTelemetryClient(ik)
	client.TrackEvent("Client connected")

	// API Key test
	appInsightsName, ok := credentials["appInsightsName"].(string)
	if !ok {
		return fmt.Errorf(
			"can't find app insights name in the credentials",
		)
	}
	APIKey, ok := credentials["APIKey"].(string)
	if !ok {
		return fmt.Errorf(
			"can't find API Key in the credentials",
		)
	}

	requestsCountMetricsAPIUrl := fmt.Sprintf(
		"https://api.applicationinsights.io/v1/apps/%s"+
			"/metrics/requests/count",
		appInsightsName,
	)
	req, err := http.NewRequest(
		"GET",
		requestsCountMetricsAPIUrl,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error validating the API Key usage: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(
		"X-Api-Key",
		APIKey,
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error validating the API Key usage: %s", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf(
			"error validating the API Key usage: response code not = 200",
		)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("error validating the API Key usage: %s", err)
	}
	return nil
}
