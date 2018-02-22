package sqldb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

//TODO: implement db only validation
func (d *dbOnlyManager) ValidateProvisioningParameters(
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
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
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlDBOnlyInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *mssqlDBOnlyInstanceDetails",
		)
	}

	//Parent should be set by the framework, but return an error if it is not set.
	if instance.Parent == nil {
		return nil, nil, errors.New("parent instance not set")
	}
	//Assume refererence instance is a vm only instance. Fail if not
	pdt, ok := instance.Parent.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.Details as " +
				"*mssqlVMOnlyInstanceDetails",
		)
	}

	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.DatabaseName = generate.NewIdentifier()
	dt.FullyQualifiedDomainName = fmt.Sprintf(
		"%s.%s",
		pdt.ServerName,
		d.sqlDatabaseDNSSuffix,
	)

	return dt, instance.SecureDetails, nil
}

func (d *dbOnlyManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlDBOnlyInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *mssqlDBOnlyInstanceDetails",
		)
	}

	//Parent should be set by the framework, but return an error if it is not set.
	if instance.Parent == nil {
		return nil, nil, errors.New("parent instance not set")
	}
	pdt, ok := instance.Parent.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Parent.Details as " +
				"*mssqlVMOnlyInstanceDetails",
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
	//No output, so ignore the output
	_, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.Parent.ResourceGroup,
		instance.Parent.Location,
		armTemplateDBOnlyBytes,
		nil, // Go template params
		p,
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return dt, instance.SecureDetails, nil
}
