package postgresql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

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
	URI         string   `json:"uri"`
	SSLRequired bool     `json:"sslRequired"`
	Tags        []string `json:"tags"`
}

func getDBMSCommonProvisionParamSchema() map[string]service.ParameterSchema {
	p := map[string]service.ParameterSchema{}

	p["firewallRules"] = &service.ArrayParameterSchema{
		Description: "Firewall rules to apply to instance. " +
			"If left unspecified, defaults to only Azure IPs",
		ItemsSchema: &service.ObjectParameterSchema{
			Description: "Individual Firewall Rule",
			Properties: map[string]service.ParameterSchema{
				"name": &service.SimpleParameterSchema{
					Type:        "string",
					Description: "Name of firewall rule",
					Required:    true,
				},
				"startIPAddress": &service.SimpleParameterSchema{
					Type:        "string",
					Description: "Start of firewall rule range",
					Required:    true,
				},
				"endIPAddress": &service.SimpleParameterSchema{
					Type:        "string",
					Description: "End of firewall rule range",
					Required:    true,
				},
			},
		},
	}

	p["sslEnforcement"] = &service.SimpleParameterSchema{
		Type: "string",
		Description: "Specifies whether the server requires the use of TLS" +
			" when connecting. Left unspecified, SSL will be enforced",
		AllowedValues: []string{"", "enabled", "disabled"},
		Default:       "",
	}
	return p
}
