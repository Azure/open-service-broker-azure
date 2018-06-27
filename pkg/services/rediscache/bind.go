package rediscache

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"fmt"
	"net/url"
)

func (s *serviceManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*instanceDetails)

	redisPort := 6379
	return credentials{
		Host:     dt.FullyQualifiedDomainName,
		Password: string(dt.PrimaryKey),
		Port:     redisPort,
		URI: fmt.Sprintf(
			"redis://:%s@%s:%d",
			url.QueryEscape(string(dt.PrimaryKey)),
			dt.FullyQualifiedDomainName,
			redisPort,
		),
	}, nil
}
