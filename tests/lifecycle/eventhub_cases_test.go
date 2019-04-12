// +build !unit

package lifecycle

import (
	"context"
	"fmt"
	"github.com/Azure/azure-event-hubs-go"
	"time"
)

var eventhubsTestCases = []serviceLifecycleTestCase{
	{
		group:     "eventhubs",
		name:      "eventhubs",
		serviceID: "7bade660-32f1-4fd7-b9e6-d416d975170b",
		planID:    "80756db5-a20c-495d-ae70-62cf7d196a3c",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"tags": map[string]interface{}{
				"latest-operation": "provision",
			},
		},
		testCredentials: testEventhubCreds,
	},
}

func testEventhubCreds(credentials map[string]interface{}) error {

	hubName := credentials["eventHubName"].(string)
	connStrP1 := credentials["connectionString"].(string)
	connStr := connStrP1 + ";EntityPath=" + hubName

	hub, err := eventhub.NewHubFromConnectionString(connStr)

	if err != nil {
		return fmt.Errorf("error creating eventhub client object: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Second)
	defer cancel()

	err = hub.Send(ctx, eventhub.NewEventFromString("hello, world!"))

	if err != nil {
		return fmt.Errorf("error sending message to eventhub: %s", err)

	}

	handler := func(c context.Context, event *eventhub.Event) error {
		fmt.Println(string(event.Data))
		return nil
	}

	runtimeInfo, err := hub.GetRuntimeInformation(ctx)

	if err != nil {
		return fmt.Errorf("error getting RunTimeInformation from eventhub: %s", err)
	}

	for _, partitionID := range runtimeInfo.PartitionIDs {
		_, err := hub.Receive(ctx, partitionID, handler, eventhub.ReceiveWithLatestOffset())
		if err != nil {
			return fmt.Errorf("error receiving messages from eventhub: %s", err)
		}

	}

	hub.Close(context.Background())

	return nil
}
