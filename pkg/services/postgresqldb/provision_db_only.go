package postgresqldb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *dbOnlyManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
	_ service.SecureProvisioningParameters,
) error {
	return nil
}

func (d *dbOnlyManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
		service.NewProvisioningStep("setupDatabase", d.setupDatabase),
		service.NewProvisioningStep("createExtensions", d.createExtensions),
	)
}

func (d *dbOnlyManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*dbOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details " +
				"as *dbOnlyPostgresqlInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.DatabaseName = generate.NewIdentifier()
	return dt, instance.SecureDetails, nil
}

func (d *dbOnlyManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.Details " +
				"as *dbmsOnlyPostgresqlInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details " +
				"as *dbOnlyPostgresqlInstanceDetails",
		)
	}
	armTemplateParameters := map[string]interface{}{
		"serverName":   pdt.ServerName,
		"databaseName": dt.DatabaseName,
	}
	_, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.Parent.ResourceGroup,
		instance.Parent.Location,
		armTemplateDBOnlyBytes,
		nil, // Go template params
		armTemplateParameters,
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

func (d *dbOnlyManager) setupDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.Details " +
				"as *dbmsOnlyPostgresqlInstanceDetails",
		)
	}
	spdt, ok :=
		instance.Parent.SecureDetails.(*dbmsOnlyPostgresqlSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.SecureDetails " +
				"as *dbmsOnlyPostgresqlSecureInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details " +
				"as *dbOnlyPostgresqlInstanceDetails",
		)
	}
	err := setupDatabase(
		pdt.EnforceSSL,
		pdt.ServerName,
		spdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	if err != nil {
		return nil, nil, err
	}
	return dt, instance.SecureDetails, nil
}

func (d *dbOnlyManager) createExtensions(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.Details " +
				"as *dbmsOnlyPostgresqlInstanceDetails",
		)
	}
	spdt, ok :=
		instance.Parent.SecureDetails.(*dbmsOnlyPostgresqlSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.SecureDetails " +
				"as *dbmsOnlyPostgresqlSecureInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details " +
				"as *dbOnlyPostgresqlInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*DatabaseProvisioningParameters)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*postgresql.DatabaseProvisioningParameters",
		)
	}

	if len(pp.Extensions) > 0 {
		err := createExtensions(
			pdt.EnforceSSL,
			pdt.ServerName,
			spdt.AdministratorLoginPassword,
			pdt.FullyQualifiedDomainName,
			dt.DatabaseName,
			pp.Extensions,
		)
		if err != nil {
			return nil, nil, err
		}
	}
	return dt, instance.SecureDetails, nil
}
