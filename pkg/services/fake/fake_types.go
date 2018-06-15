package fake

import "github.com/Azure/open-service-broker-azure/pkg/service"

// GetEmptyInstanceDetails returns an "empty" service-specific object that
// can be populated with data during unmarshaling of JSON to an Instance
func (s *ServiceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return nil
}

// GetEmptyBindingDetails returns an "empty" service-specific object that
// can be populated with data during unmarshaling of JSON to a Binding
func (s *ServiceManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}
