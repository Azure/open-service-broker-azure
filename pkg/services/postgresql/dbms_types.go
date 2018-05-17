package postgresql

type dbmsProvisioningParameters struct {
	SSLEnforcement   string         `json:"sslEnforcement"`
	FirewallRules    []firewallRule `json:"firewallRules"`
	Cores            *int64         `json:"cores"`
	Storage          *int64         `json:"storage"`
	HardwareFamily   string         `json:"hardwareFamily"`
	BackupRetention  *int64         `json:"backupRetention"`
	BackupRedundancy string         `json:"backupRedundancy"`
}

type dbmsUpdatingParameters struct {
	SSLEnforcement  string         `json:"sslEnforcement"`
	FirewallRules   []firewallRule `json:"firewallRules"`
	Cores           *int64         `json:"cores"`
	Storage         *int64         `json:"storage"`
	BackupRetention *int64         `json:"backupRetention"`
}

type firewallRule struct {
	Name    string `json:"name"`
	StartIP string `json:"startIPAddress"`
	EndIP   string `json:"endIPAddress"`
}

type dbmsInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	ServerName               string `json:"server"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
}

type secureDBMSInstanceDetails struct {
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}
