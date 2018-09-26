// +build !unit

package lifecycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Sentiment struct for credentials test
type Sentiment struct {
	Documents []Documents `json:"documents"`
}

// Documents struct for credentials test
type Documents struct {
	Language string `json:"language"`
	ID       string `json:"id"`
	Text     string `json:"text"`
}

var textanalyticsTestCases = []serviceLifecycleTestCase{
	{ // Text analytics free tier
		group:     "textanalytics",
		name:      "text-analytics-free",
		serviceID: "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
		planID:    "d5a0f91f-10da-42fc-b792-656a616d9ec2",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
		testCredentials: testTextAnalyticsCreds,
	},
	{ // Text analytics standard-s0 tier
		group:     "textanalytics",
		name:      "text-analytics-s0",
		serviceID: "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
		planID:    "7f49713b-2689-4c66-bac9-85a024c0fb9e",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
		testCredentials: testTextAnalyticsCreds,
	},
	{ // Text analytics standard-s1 tier
		group:     "textanalytics",
		name:      "text-analytics-s1",
		serviceID: "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
		planID:    "55575612-482b-4260-b67e-69be36d83a54",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
		testCredentials: testTextAnalyticsCreds,
	},
	{ // Text analytics standard-s2 tier
		group:     "textanalytics",
		name:      "text-analytics-s2",
		serviceID: "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
		planID:    "76bc6a2f-1364-4ef2-8037-d7cfff48f3b6",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
		testCredentials: testTextAnalyticsCreds,
	},
	{ // Text analytics standard-s3 tier
		group:     "textanalytics",
		name:      "text-analytics-s3",
		serviceID: "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
		planID:    "1d9a9e7c-80ac-4f23-aabe-876125541f59",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
		testCredentials: testTextAnalyticsCreds,
	},
	{ // Text analytics standard-s4 tier
		group:     "textanalytics",
		name:      "text-analytics-s4",
		serviceID: "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
		planID:    "b9db834d-1350-4c50-adaf-f1e59efa2381",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
		testCredentials: testTextAnalyticsCreds,
	},
}

func testTextAnalyticsCreds(credentials map[string]interface{}) error {

	// Test data
	testText := Documents{
		Language: "en",
		ID:       "1",
		Text:     "Super positive happy text.",
	}

	documents := []Documents{}
	sentiment := Sentiment{documents}
	sentiment.Documents = append(sentiment.Documents, testText)

	j, _ := json.Marshal(sentiment)
	b := bytes.NewBuffer(j)

	req, err := http.NewRequest(
		"POST",
		credentials["textAnalyticsEndpoint"].(string)+"/sentiment", b,
	)
	if err != nil {
		return fmt.Errorf("error validating the text analytics arguments: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(
		"Ocp-Apim-Subscription-Key",
		credentials["textAnalyticsKey"].(string),
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error validating the text analytics arguments: %s", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf(
			"error validating the text analytics arguments: response code not = 200",
		)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("error validating the text analytics arguments: %s", err)
	}

	return nil
}
