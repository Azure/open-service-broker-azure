package mysql

import "github.com/Azure/open-service-broker-azure/pkg/service"

type bindingDetails struct {
	LoginName string               `json:"loginName"`
	Password  service.SecureString `json:"password"`
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
