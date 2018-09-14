package rediscache

import "github.com/Azure/open-service-broker-azure/pkg/service"

type instanceDetails struct {
	ARMDeploymentName        string               `json:"armDeployment"`
	ServerName               string               `json:"server"`
	FullyQualifiedDomainName string               `json:"fullyQualifiedDomainName"`
	NonSSLEnabled            bool                 `json:"nonSSLEnabled"`
	PrimaryKey               service.SecureString `json:"primaryKey"`
}

type credentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	URI      string `json:"uri"`
}

type bindingDetails struct {
	LoginName string               `json:"loginName"`
	Password  service.SecureString `json:"password"`
}

func (s *serviceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &instanceDetails{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &bindingDetails{}
}
