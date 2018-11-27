package iothub

import "github.com/Azure/open-service-broker-azure/pkg/service"

type instanceDetails struct {
	ARMDeploymentName string               `json:"armDeployment"`
	IoTHubName        string               `json:"iotHubName"`
	KeyName           string               `json:"keyName"`
	Key               service.SecureString `json:"key"`
}

type credentials struct {
	IoTHubName       string `json:"iotHubName"`
	HostName         string `json:"hostName"`
	KeyName          string `json:"keyName"`
	Key              string `json:"key"`
	ConnectionString string `json:"connectionString"`
}

func (i *iotHubManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &instanceDetails{}
}

func (i *iotHubManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}