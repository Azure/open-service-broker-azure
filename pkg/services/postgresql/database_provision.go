package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *databaseManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
	_ service.SecureProvisioningParameters,
) error {
	return nil
}

func (d *databaseManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
		service.NewProvisioningStep("setupDatabase", d.setupDatabase),
		service.NewProvisioningStep("createExtensions", d.createExtensions),
	)
}

func (d *databaseManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *postgresql.databaseInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.DatabaseName = generate.NewIdentifier()
	return dt, instance.SecureDetails, nil
}

func (d *databaseManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	pdt, ok := instance.Parent.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.Details as " +
				"*postgresql.dbmsInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *postgresql.databaseInstanceDetails",
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
		databaseARMTemplateBytes,
		nil, // Go template params
		armTemplateParameters,
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

func (d *databaseManager) setupDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	pdt, ok := instance.Parent.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.Details as " +
				"*postgresql.dbmsInstanceDetails",
		)
	}
	spdt, ok :=
		instance.Parent.SecureDetails.(*secureDBMSInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.SecureDetails as " +
				"*postgresql.secureDBMSInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *postgresql.databaseInstanceDetails",
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

func (d *databaseManager) createExtensions(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	pdt, ok := instance.Parent.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.Details as " +
				"*postgresql.dbmsInstanceDetails",
		)
	}
	spdt, ok :=
		instance.Parent.SecureDetails.(*secureDBMSInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.SecureDetails as " +
				"*postgresql.secureDBMSInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *postgresql.databaseInstanceDetails",
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
