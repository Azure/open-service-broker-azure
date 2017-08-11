package postgresql

type postgresqlProvisioningParameters struct {
	Location string            `json:"location"`
	Tags     map[string]string `json:"tags"`
}

type postgresqlProvisioningContext struct {
	ResourceGroupName          string `json:"resourceGroup"`
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	DatabaseName               string `json:"database"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
}

type postgresqlBindingParameters struct {
}

type postgresqlBindingContext struct {
	LoginName string `json:"loginName"`
}

type postgresqlCredentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (m *module) GetEmptyProvisioningParameters() interface{} {
	return &postgresqlProvisioningParameters{}
}

func (m *module) GetEmptyProvisioningContext() interface{} {
	return &postgresqlProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() interface{} {
	return &postgresqlBindingParameters{}
}

func (m *module) GetEmptyBindingContext() interface{} {
	return &postgresqlBindingContext{}
}

func (m *module) GetEmptyCredentials() interface{} {
	return &postgresqlCredentials{}
}
