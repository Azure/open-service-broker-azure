package mssqldr

import (
	"context"
	"fmt"
	"strings"

	sqlSDK "github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
)

func validateServer(
	ctx context.Context,
	serversClient *sqlSDK.ServersClient,
	resourceGroup string,
	serverName string,
	expectedVersion string,
	expectedLocation string,
) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	result, err := serversClient.Get(
		ctx,
		resourceGroup,
		serverName,
	)
	if err != nil {
		return "", fmt.Errorf("error getting the sql server: %s", err)
	}
	if result.Name == nil {
		return "", fmt.Errorf(
			"can't find sql server %s in the resource group %s",
			serverName,
			resourceGroup,
		)
	}
	if *result.Version != expectedVersion {
		return "", fmt.Errorf(
			"sql server version validation failed, "+
				"expected version: %s, actual version: %s",
			expectedVersion,
			*result.Version,
		)
	}
	expectedLocation = strings.Replace(
		strings.ToLower(expectedLocation),
		" ",
		"",
		-1,
	)
	actualLocation := strings.Replace(
		strings.ToLower(*result.Location),
		" ",
		"",
		-1,
	)
	if expectedLocation != actualLocation {
		return "", fmt.Errorf(
			"sql server location validation failed, "+
				"expected location: %s, actual location: %s",
			expectedLocation,
			actualLocation,
		)
	}
	if *result.FullyQualifiedDomainName == "" {
		return "", fmt.Errorf(
			"sql server details doesn't contain FQDN",
		)
	}
	return *result.FullyQualifiedDomainName, nil
}
