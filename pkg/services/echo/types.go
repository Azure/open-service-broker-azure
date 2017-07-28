package echo

type ProvisioningParameters struct {
	Message string `json:"message"`
}

type ProvisioningContext struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

type BindingParameters struct {
	Message string `json:"message"`
}

type BindingContext struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

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
