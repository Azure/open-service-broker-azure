package echo

// ProvisioningParameters represents parameters specific to provisioning a
// service using the echo service module. Note that, ordinarily, service module-
// specific types such as this do not need to be exported. An exception is made
// here because the echo service module is used to facilitate testing of the
// broker framework itself.
type ProvisioningParameters struct {
	Message string `json:"message"`
}

// ProvisioningContext represents context collected and modified over the course
// of the echo service module's provisioning and deprovisioning processes. Note
// that, ordinarily, service module-specific types such as this do not need to
// be exported. An exception is made here because the echo service module is
// used to facilitate testing of the broker framework itself.
type ProvisioningContext struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

// BindingParameters represents parameters specific to binding to a service
// instance using the echo service module. Note that, ordinarily, service
// module-specific types such as this do not need to be exported. An exception
// is made here because the echo service module is used to facilitate testing of
// the broker framework itself.
type BindingParameters struct {
	Message string `json:"message"`
}

// BindingContext represents context collected and modified over the course
// of the echo service module's binding and unbinding processes. Note that,
// ordinarily, service module-specific types such as this do not need to be
// exported. An exception is made here because the echo service module is used
// to facilitate testing of the broker framework itself.
type BindingContext struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

// Credentials generally represent credentials AND/OR ANY OTHER DETAILS (e.g.
// URLs, port numbers, etc.) that will be conveyed back to the client upon
// successful completion of a bind. In the specific case of the echo service
// module, which doesn't do much of anything (other than generate a messageID
// and write some logs), there are no important details to convey back to the
// client. The messageID is included just to provide an example of HOW details
// such as these can be conveyed to the client. Note that, ordinarily, service
// module-specific types such as this do not need to be exported. An exception
// is made here because the echo service module is used to facilitate testing of
// the broker framework itself.
type Credentials struct {
	MessageID string `json:"messageId"`
}

func (m *module) GetEmptyProvisioningParameters() interface{} {
	return &ProvisioningParameters{}
}

func (m *module) GetEmptyProvisioningContext() interface{} {
	return &ProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() interface{} {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() interface{} {
	return &BindingContext{}
}

func (m *module) GetEmptyCredentials() interface{} {
	return &Credentials{}
}
