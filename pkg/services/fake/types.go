package fake

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters represents parameters specific to provisioning a
// service using the fake service module. Note that, ordinarily, service module-
// specific types such as this do not need to be exported. An exception is made
// here because the fake service module is used to facilitate testing of the
// broker framework itself.
type ProvisioningParameters struct {
	SomeParameter string `json:"someParameter"`
}

// InstanceDetails represents details collected and modified over the course
// of a fake service instance's provisioning and deprovisioning processes. Note
// that, ordinarily, service-specific types such as this do not need to be
// exported. An exception is made here because the fake service module is used
// to facilitate testing of the broker framework itself.
type InstanceDetails struct {
	ResourceGroupName string `json:"resourceGroup"`
}

// UpdatingParameters represents parameters specific to binding to a service
// instance using the fake service module. Note that, ordinarily, service
// module-specific types such as this do not need to be exported. An exception
// is made here because the fake service module is used to facilitate testing of
// the broker framework itself.
type UpdatingParameters struct {
	SomeParameter string `json:"someParameter"`
}

// BindingParameters represents parameters specific to binding to a service
// instance using the fake service module. Note that, ordinarily,
// service-specific types such as this do not need to be exported. An exception
// is made here because the fake service module is used to facilitate testing of
// the broker framework itself.
type BindingParameters struct {
	SomeParameter string `json:"someParameter"`
}

// BindingDetails represents details collected and modified over the course
// of a fake service instance's binding and unbinding processes. Note that,
// ordinarily, service-specific types such as this do not need to be exported.
// An exception is made here because the fake service module is used to
// facilitate testing of the broker framework itself.
type BindingDetails struct {
}

// Credentials generally represent credentials AND/OR ANY OTHER DETAILS (e.g.
// URLs, port numbers, etc.) that will be conveyed back to the client upon
// successful completion of a bind. In the specific case of the fake service
// module, which doesn't do much of anything (other than generate a messageID
// and write some logs), there are no important details to convey back to the
// client. The messageID is included just to provide an example of HOW details
// such as these can be conveyed to the client. Note that, ordinarily, service
// module-specific types such as this do not need to be exported. An exception
// is made here because the fake service module is used to facilitate testing of
// the broker framework itself.
type Credentials struct {
}

// GetEmptyProvisioningParameters returns an empty instance of module-specific
// provisioningParameters
func (
	s *ServiceManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

// GetEmptyInstanceDetails returns an empty instance of a service-specific
// instance details
func (
	s *ServiceManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &InstanceDetails{}
}

// GetEmptyUpdatingParameters returns an empty instance of module-specific
// updatingParameters
func (
	s *ServiceManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

// GetEmptyBindingParameters returns an empty instance of module-specific
// bindingParameters
func (s *ServiceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

// GetEmptyBindingDetails returns an empty instance of a service-specific
// bindingDetails
func (s *ServiceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &BindingDetails{}
}

// GetEmptyCredentials returns an empty instance of module-specific
// credentials
func (s *ServiceManager) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}
