package postgresql

type allInOneProvisioningParameters struct {
	dbmsProvisioningParameters     `json:",squash"`
	databaseProvisioningParameters `json:",squash"`
}

type allInOneInstanceDetails struct {
	dbmsInstanceDetails `json:",squash"`
	DatabaseName        string `json:"database"`
}

type secureAllInOneInstanceDetails struct {
	secureDBMSInstanceDetails `json:",squash"`
}
