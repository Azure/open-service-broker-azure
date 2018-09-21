package rediscache

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

const enabled = "enabled"

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
) (service.InstanceDetails, error) {
	return &instanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
		ServerName:        uuid.NewV4().String(),
	}, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)

	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}

	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		instance.ProvisioningParameters.GetString("location"),
		armTemplateBytes,
		buildGoTemplate(instance, provision),
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	var ok bool
	dt.FullyQualifiedDomainName, ok = outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving primary key from deployment: %s",
			err,
		)
	}
	dt.PrimaryKey = service.SecureString(primaryKey)

	dt.NonSSLEnabled = (instance.ProvisioningParameters.GetString("enableNonSslPort") == enabled) // nolint: lll

	return dt, err
}

const provision = "provision"
const update = "update"

func buildGoTemplate(
	instance service.Instance,
	mode string,
) map[string]interface{} {
	var pp *service.ProvisioningParameters
	if mode == provision {
		pp = instance.ProvisioningParameters
	} else if mode == update {
		pp = instance.UpdatingParameters
	}

	dt := instance.Details.(*instanceDetails)
	plan := instance.Plan

	var enableNonSslPort string
	if pp.GetString("enableNonSslPort") == enabled {
		enableNonSslPort = "true"
	} else {
		enableNonSslPort = "false"
	}

	subnetSettings := pp.GetObject("subnetSettings")

	return map[string]interface{}{ // ARM template params
		"location":           pp.GetString("location"),
		"serverName":         dt.ServerName,
		"enableNonSslPort":   enableNonSslPort,
		"redisCacheSKU":      plan.GetProperties().Extended["redisCacheSKU"],
		"redisCacheFamily":   plan.GetProperties().Extended["redisCacheFamily"],
		"redisCacheCapacity": pp.GetInt64("skuCapacity"),
		"subnetId":           subnetSettings.GetString("subnetId"),
		"staticIP":           subnetSettings.GetString("staticIP"),
	}
}
