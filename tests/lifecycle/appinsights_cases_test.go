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
		name:            "appinsights",
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
