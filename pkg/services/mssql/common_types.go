package mssql

type bindingDetails struct {
	LoginName string `json:"loginName"`
}

type secureBindingDetails struct {
	Password string `json:"password"`
}

// Credentials encapsulates MSSQL-specific coonection details and credentials.
type credentials struct {
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Database string   `json:"database"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	URI      string   `json:"uri"`
	Tags     []string `json:"tags"`
	JDBC     string   `json:"jdbcUrl"`
	Encrypt  bool     `json:"encrypt"`
}
