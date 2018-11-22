// +build !unit

package lifecycle

import (
	"fmt"
	"time"

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
	ik, ok := credentials["instrumentationKey"]
	if !ok {
		return fmt.Errorf(
			"can't find instrumentation key in the credentials",
		)
	}
	client := appinsights.NewTelemetryClient(ik)
	client.TrackEvent("Client connected")
	trace := appinsights.NewTraceTelemetry("message", appinsights.Warning)
	trace.Properties["test"] = "osba"
	trace.Timestamp = time.Now().Sub(time.Minute)
	client.Track(trace)
}
