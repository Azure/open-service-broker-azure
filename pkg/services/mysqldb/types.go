package mysqldb

// ServerProvisioningParameters encapsulates non-sensitive MySQL-specific
// server provisioning options
type ServerProvisioningParameters struct {
	SSLEnforcement  string `json:"sslEnforcement"`
	FirewallIPStart string `json:"firewallStartIPAddress"`
	FirewallIPEnd   string `json:"firewallEndIPAddress"`
}

// SecureServerProvisioningParameters encapsulates sensitive MySQL-specific
// server provisioning options
type SecureServerProvisioningParameters struct{}

// UpdatingParameters encapsulates MySQL-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates MySQL-specific binding options
type BindingParameters struct {
}

type mysqlBindingDetails struct {
	LoginName string `json:"loginName"`
}

type mysqlSecureBindingDetails struct {
	Password string `json:"password"`
}

// Credentials encapsulates MySQL-specific coonection details and credentials.
type Credentials struct {
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	Database    string   `json:"database"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	SSLRequired bool     `json:"sslRequired"`
	URI         string   `json:"uri"`
	Tags        []string `json:"tags"`
}
