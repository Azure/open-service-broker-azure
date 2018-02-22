package postgresqldb

// BindingParameters encapsulates non-sensitive PostgreSQL-specific binding
// options
type BindingParameters struct {
}

// SecureBindingParameters encapsulates sensitive PostgreSQL-specific binding
// options
type SecureBindingParameters struct {
}

type postgresqlBindingDetails struct {
	LoginName string `json:"loginName"`
}

type postgresqlSecureBindingDetails struct {
	Password string `json:"password"`
}

// Credentials encapsulates PostgreSQL-specific coonection details and
// credentials.
type Credentials struct {
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	Database    string   `json:"database"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	URI         string   `json:"uri"`
	SSLRequired bool     `json:"sslRequired"`
	Tags        []string `json:"tags"`
}
