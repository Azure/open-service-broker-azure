package rediscache

import (
	"fmt"
	"net/url"

	"github.com/Azure/open-service-broker-azure/pkg/service"
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

	var redisPort int
	var scheme string
	if dt.NonSSLEnabled {
		redisPort = 6379
		scheme = "redis"
	} else {
		redisPort = 6380
		scheme = "rediss"
	}

	return credentials{
		Host:     dt.FullyQualifiedDomainName,
		Password: string(dt.PrimaryKey),
		Port:     redisPort,
		URI: fmt.Sprintf(
			"%s://:%s@%s:%d",
			scheme,
			url.QueryEscape(string(dt.PrimaryKey)),
			dt.FullyQualifiedDomainName,
			redisPort,
		),
	}, nil
}
