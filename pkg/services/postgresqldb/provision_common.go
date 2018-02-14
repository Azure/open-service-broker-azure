package postgresqldb

import (
	"context"
	"fmt"

	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-04-30-preview/postgresql" // nolint: lll
	uuid "github.com/satori/go.uuid"
)

func getAvailableServerName(
	ctx context.Context,
	checkNameAvailabilityClient postgresSDK.CheckNameAvailabilityClient,
) (string, error) {
	for {
		serverName := uuid.NewV4().String()
		nameAvailability, err := checkNameAvailabilityClient.Execute(
			ctx,
			postgresSDK.NameAvailabilityRequest{
				Name: &serverName,
			},
		)
		if err != nil {
			return "", fmt.Errorf(
				"error determining server name availability: %s",
				err,
			)
		}
		if *nameAvailability.NameAvailable {
			return serverName, nil
		}
	}
}
