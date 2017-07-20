package echo

type echoProvisioningParameters struct {
	Message string `json:"message"`
}

type echoProvisioningResult struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

type echoBindingParameters struct {
	Message string `json:"message"`
}

type echoBindingResult struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

func (m *module) GetEmptyProvisioningParameters() interface{} {
	return &echoProvisioningParameters{}
}

func (m *module) GetEmptyProvisioningResult() interface{} {
	return &echoProvisioningResult{}
}

func (m *module) GetEmptyBindingParameters() interface{} {
	return &echoBindingParameters{}
}

func (m *module) GetEmptyBindingResult() interface{} {
	return &echoBindingResult{}
}
