// +build experimental

package rediscache

import (
	"fmt"
	"net/url"

	"open-service-broker-azure/pkg/service"
)

func (s *serviceManager) Bind(
	service.Instance,
	service.BindingParameters,
	service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return nil, nil, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}
	sdt := secureInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, err
	}

	redisPort := 6379
	return credentials{
		Host:     dt.FullyQualifiedDomainName,
		Password: sdt.PrimaryKey,
		Port:     redisPort,
		URI: fmt.Sprintf(
			"redis://:%s@%s:%d",
			url.QueryEscape(sdt.PrimaryKey),
			dt.FullyQualifiedDomainName,
			redisPort,
		),
	}, nil
}
