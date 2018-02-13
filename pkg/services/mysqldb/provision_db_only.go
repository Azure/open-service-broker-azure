package mysqldb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	uuid "github.com/satori/go.uuid"
)

func (d *dbOnlyManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	return nil
}

func (d *dbOnlyManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
	)
}
func (d *dbOnlyManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}
	//We aren't using any of these, but validate it can be type cast
	_, ok = instance.ProvisioningParameters.(*DatabaseProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*mysqldb.DatabaseProvisioningParameters",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.DatabaseName = generate.NewIdentifier()

	return dt, nil
}

func (d *dbOnlyManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Parent.Details " +
				"as *dbmsOnlyMysqlInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}
	_, ok = instance.ProvisioningParameters.(*DatabaseProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters " +
				"as *mysql.DatabaseProvisioningParameters",
		)
	}

	armTemplateParameters := map[string]interface{}{ // ARM template params
		"serverName":   pdt.ServerName,
		"databaseName": dt.DatabaseName,
	}
	_, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.Parent.ResourceGroup,
		instance.Parent.Location,
		dbOnlyArmTemplateBytes,
		nil, // Go template params
		armTemplateParameters,
		instance.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	return dt, nil
}
