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
		service.NewProvisioningStep("validateServer", d.validateServer),
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

func (d *dbmsRegisteredManager) validateServer(
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

	if err = validateServerAdmin(
		dt.AdministratorLogin,
		string(dt.AdministratorLoginPassword),
		dt.FullyQualifiedDomainName,
	); err != nil {
		return nil, err
	}

	return instance.Details, nil
}
