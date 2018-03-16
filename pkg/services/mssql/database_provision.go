package mssql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

// TODO: implement db only validation
func (d *databaseManager) ValidateProvisioningParameters(
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
) error {
	return nil
}

func (d *databaseManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
	)
}

func (d *databaseManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *mssql.databaseInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*DatabaseProvisioningParams)
	if !ok {
		return nil, nil, errors.New(
			"error casting provisioningParameters as " +
				"*mssql.DatabaseProvisioningParams",
		)
	}

	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.DatabaseName = generate.NewIdentifier()

	if !pp.DisableTDE {
		dt.TransparentDataEncryption = true
	}

	return dt, instance.SecureDetails, nil
}

func (d *databaseManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *mssql.databaseInstanceDetails",
		)
	}
	pdt, ok := instance.Parent.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.Details as *mssql.dbmsInstanceDetails",
		)
	}
	p := map[string]interface{}{ // ARM template params
		"serverName":   pdt.ServerName,
		"databaseName": dt.DatabaseName,
		"edition":      instance.Plan.GetProperties().Extended["edition"],
		"requestedServiceObjectiveName": instance.Plan.GetProperties().
			Extended["requestedServiceObjectiveName"],
		"maxSizeBytes": instance.Plan.GetProperties().Extended["maxSizeBytes"],
	}
	goTemplateParams := map[string]interface{}{}
	goTemplateParams["transparentDataEncryption"] = dt.TransparentDataEncryption

	// No output, so ignore the output
	_, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.Parent.ResourceGroup,
		instance.Parent.Location,
		databaseARMTemplateBytes,
		goTemplateParams,
		p,
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return dt, instance.SecureDetails, nil
}
