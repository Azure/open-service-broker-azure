package iothub

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (i *iotHubManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", i.preProvision),
		service.NewProvisioningStep("deployARMTemplate", i.deployARMTemplate),
	)
}

func (i *iotHubManager) preProvision(
	context.Context,
	service.Instance,
) (service.InstanceDetails, error) {
	return &instanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
		IoTHubName:        uuid.NewV4().String(),
	}, nil
}

func (i *iotHubManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}

	outputs, err := i.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		instance.ProvisioningParameters.GetString("location"),
		armTemplateBytes,
		buildGoTemplateParams(instance),
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	keyInfo := outputs["keyInfo"].(map[string]interface{})
	var ok bool

	dt.KeyName, ok = keyInfo["keyName"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving keyName from deployment: %s",
			err,
		)
	}

	key, ok := keyInfo["primaryKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving key from deployment: %s",
			err,
		)
	}
	dt.Key = service.SecureString(key)
	return dt, err
}

func buildGoTemplateParams(instance service.Instance) map[string]interface{} {
	dt := instance.Details.(*instanceDetails)
	params := map[string]interface{}{
		"iotHubName": dt.IoTHubName,
		"location":   instance.ProvisioningParameters.GetString("location"),
	}

	var skuName string
	switch instance.Plan.GetName() {
	case planF1:
		skuName = "F1"
	case planS1:
		skuName = "S1"
	case planS2:
		skuName = "S2"
	case planS3:
		skuName = "S3"
	case planB1:
		skuName = "B1"
	case planB2:
		skuName = "B2"
	case planB3:
		skuName = "B3"
	}
	params["skuName"] = skuName

	if instance.Plan.GetName() == planF1 {
		params["skuUnits"] = 1
		params["partitionCount"] = 2
	} else {
		params["skuUnits"] = instance.ProvisioningParameters.GetInt64("units")
		params["partitionCount"] = instance.ProvisioningParameters.GetInt64("partitionCount")
	}
	return params
}