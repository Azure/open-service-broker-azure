package servicebus

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (nm *namespaceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", nm.preProvision),
		service.NewProvisioningStep("deployARMTemplate", nm.deployARMTemplate),
	)
}

func (nm *namespaceManager) preProvision(
	context.Context,
	service.Instance,
) (service.InstanceDetails, error) {
	return &namespaceInstanceDetails{
		ARMDeploymentName:       uuid.NewV4().String(),
		ServiceBusNamespaceName: "sb-" + uuid.NewV4().String(),
	}, nil
}

func (nm *namespaceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*namespaceInstanceDetails)
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	outputs, err := nm.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		instance.ProvisioningParameters.GetString("location"),
		armTemplateBytes,
		map[string]interface{}{
			"location":                instance.ProvisioningParameters.GetString("location"), // nolint: lll
			"serviceBusNamespaceName": dt.ServiceBusNamespaceName,
			"serviceBusSku": instance.Plan.
				GetProperties().Extended["serviceBusSku"],
		},
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	var ok bool
	connectionString, ok := outputs["connectionString"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving connection string from deployment: %s",
			err,
		)
	}
	dt.ConnectionString = service.SecureString(connectionString)

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving primary key from deployment: %s",
			err,
		)
	}
	dt.PrimaryKey = service.SecureString(primaryKey)

	return dt, err
}
