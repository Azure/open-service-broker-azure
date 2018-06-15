package mysql

import "github.com/Azure/open-service-broker-azure/pkg/service"

type allInOneInstanceDetails struct {
	dbmsInstanceDetails `json:",squash"`
	DatabaseName        string `json:"database"`
}

func (a *allInOneManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &allInOneInstanceDetails{}
}

func (a *allInOneManager) GetEmptyBindingDetails() service.BindingDetails {
	return &bindingDetails{}
}
