package iothub

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (i *iotHubManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, nil
}

func (i *iotHubManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*instanceDetails)

	hostName := fmt.Sprintf("%s.azure-devices.net", dt.IoTHubName)
	c := credentials{
		IoTHubName: dt.IoTHubName,
		HostName:   hostName,
		KeyName:    dt.Keys[0].KeyName,
		Key:        string(dt.Keys[0].PrimaryKey),
	}
	c.ConnectionString = fmt.Sprintf(
		"HostName=%s;SharedAccessKeyName=%s;SharedAccessKey=%s",
		hostName,
		c.KeyName,
		c.Key,
	)

	return c, nil
}
