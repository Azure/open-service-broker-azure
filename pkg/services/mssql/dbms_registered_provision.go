package mssql

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsRegisteredManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("getServer", d.getServer),
		service.NewProvisioningStep("testConnection", d.testConnection),
	)
}

func (d *dbmsRegisteredManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	return &dbmsInstanceDetails{
		ARMDeploymentName:          "",
		ServerName:                 pp.GetString("server"),
		AdministratorLogin:         pp.GetString("administratorLogin"),
		AdministratorLoginPassword: service.SecureString(pp.GetString("administratorLoginPassword")), // nolint: lll
	}, nil
}

func (d *dbmsRegisteredManager) getServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	pp := instance.ProvisioningParameters
	dt := instance.Details.(*dbmsInstanceDetails)
	resourceGroup := pp.GetString("resourceGroup")
	result, err := d.serversClient.Get(
		ctx,
		resourceGroup,
		dt.ServerName,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting sql server: %s", err)
	}
	expectedVersion :=
		instance.Service.GetProperties().Extended["version"].(string)
	if *result.Version != expectedVersion {
		return nil, fmt.Errorf(
			"sql server version validation failed, "+
				"expected version: %s, actual version: %s",
			expectedVersion,
			*result.Version,
		)
	}
	expectedLocation := strings.Replace(
		strings.ToLower(pp.GetString("location")),
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
		return nil, fmt.Errorf(
			"sql server location validation failed, "+
				"expected location: %s, actual location: %s",
			expectedLocation,
			actualLocation,
		)
	}
	dt.FullyQualifiedDomainName = *result.FullyQualifiedDomainName
	return instance.Details, nil
}

func (d *dbmsRegisteredManager) testConnection(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsInstanceDetails)
	// connect to master database
	masterDb, err := getDBConnection(
		dt.AdministratorLogin,
		string(dt.AdministratorLoginPassword),
		fmt.Sprintf("%s.%s", dt.ServerName, d.sqlDatabaseDNSSuffix),
		"master",
	)
	if err != nil {
		return nil, err
	}
	defer masterDb.Close() // nolint: errcheck

	// Is there a better approach to verify if it is a sys admin?
	rows, err := masterDb.Query("SELECT 1 FROM fn_my_permissions (NULL, 'DATABASE') WHERE permission_name='ALTER ANY USER'") // nolint: lll
	if err != nil {
		return nil, fmt.Errorf(
			`error querying SELECT from table fn_my_permissions: %s`,
			err,
		)
	}
	defer rows.Close() // nolint: errcheck
	if !rows.Next() {
		return nil, fmt.Errorf(
			`error user doesn't have permission 'ALTER ANY USER'`,
		)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			`error iterating rows`,
		)
	}

	return instance.Details, nil
}
