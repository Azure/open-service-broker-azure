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

func testConnection(
	fqdn string,
	administratorLogin string,
	administratorLoginPassword string,
) error {
	masterDb, err := getDBConnection(
		administratorLogin,
		administratorLoginPassword,
		fqdn,
		"master",
	)
	if err != nil {
		return err
	}
	defer masterDb.Close() // nolint: errcheck
	// TODO: Is there a better approach to verify if it is a sys admin?
	rows, err := masterDb.Query("SELECT 1 FROM fn_my_permissions (NULL, 'DATABASE') WHERE permission_name='ALTER ANY USER'") // nolint: lll
	if err != nil {
		return fmt.Errorf(
			`error querying SELECT from table fn_my_permissions: %s`,
			err,
		)
	}
	defer rows.Close() // nolint: errcheck
	if !rows.Next() {
		return fmt.Errorf(
			`error user doesn't have permission 'ALTER ANY USER'`,
		)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf(`server error iterating rows`)
	}
	return nil
}
