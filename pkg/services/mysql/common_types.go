// +build experimental

package mysql

type bindingDetails struct {
	LoginName string `json:"loginName"`
}

type secureBindingDetails struct {
	Password string `json:"password"`
}

type credentials struct {
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	Database    string   `json:"database"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	SSLRequired bool     `json:"sslRequired"`
	URI         string   `json:"uri"`
	Tags        []string `json:"tags"`
}
