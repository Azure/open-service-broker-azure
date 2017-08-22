package mysql

type mysqlProvisioningParameters struct {
	Location string            `json:"location"`
	Tags     map[string]string `json:"tags"`
}

type mysqlProvisioningContext struct {
	ResourceGroupName          string `json:"resourceGroup"`
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	DatabaseName               string `json:"database"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
}

type mysqlBindingParameters struct {
}

type mysqlBindingContext struct {
	LoginName string `json:"loginName"`
}

type mysqlCredentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (m *module) GetEmptyProvisioningParameters() interface{} {
	return &mysqlProvisioningParameters{}
}

func (m *module) GetEmptyProvisioningContext() interface{} {
	return &mysqlProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() interface{} {
	return &mysqlBindingParameters{}
}

func (m *module) GetEmptyBindingContext() interface{} {
	return &mysqlBindingContext{}
}

func (m *module) GetEmptyCredentials() interface{} {
	return &mysqlCredentials{}
}
