// +build experimental

package servicebus

import (
	"context"
	"fmt"

	"open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *serviceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
	)
}

func (s *serviceManager) preProvision(
	context.Context,
	service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := instanceDetails{
		ARMDeploymentName:       uuid.NewV4().String(),
		ServiceBusNamespaceName: "sb-" + uuid.NewV4().String(),
	}
	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, nil, err
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := instanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"serviceBusNamespaceName": dt.ServiceBusNamespaceName,
			"serviceBusSku": instance.Plan.
				GetProperties().Extended["serviceBusSku"],
		},
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	sdt := secureInstanceDetails{}

	var ok bool
	sdt.ConnectionString, ok = outputs["connectionString"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving connection string from deployment: %s",
			err,
		)
	}

	sdt.PrimaryKey, ok = outputs["primaryKey"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving primary key from deployment: %s",
			err,
		)
	}

	sdtMap, err := service.GetMapFromStruct(sdt)
	return instance.Details, sdtMap, err
}
