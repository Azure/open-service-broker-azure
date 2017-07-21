package echo

type echoProvisioningParameters struct {
	Message string `json:"message"`
}

type echoProvisioningContext struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

type echoBindingParameters struct {
	Message string `json:"message"`
}

type echoBindingContext struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

func (m *module) GetEmptyProvisioningParameters() interface{} {
	return &echoProvisioningParameters{}
}

func (m *module) GetEmptyProvisioningContext() interface{} {
	return &echoProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() interface{} {
	return &echoBindingParameters{}
}

func (m *module) GetEmptyBindingContext() interface{} {
	return &echoBindingContext{}
}
