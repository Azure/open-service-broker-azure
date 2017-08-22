package rediscache

type redisProvisioningParameters struct {
	Location string            `json:"location"`
	Tags     map[string]string `json:"tags"`
}

type redisProvisioningContext struct {
	ResourceGroupName        string `json:"resourceGroup"`
	ARMDeploymentName        string `json:"armDeployment"`
	ServerName               string `json:"server"`
	PrimaryKey               string `json:"primaryKey"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
}

type redisBindingParameters struct {
}

type redisBindingContext struct {
}

type redisCredentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

func (m *module) GetEmptyProvisioningParameters() interface{} {
	return &redisProvisioningParameters{}
}

func (m *module) GetEmptyProvisioningContext() interface{} {
	return &redisProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() interface{} {
	return &redisBindingParameters{}
}

func (m *module) GetEmptyBindingContext() interface{} {
	return &redisBindingContext{}
}

func (m *module) GetEmptyCredentials() interface{} {
	return &redisCredentials{}
}
