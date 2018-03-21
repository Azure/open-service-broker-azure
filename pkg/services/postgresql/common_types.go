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

// GetDBMSCommonProvisionParametersSchema generates a common schema for both
// the DBMS-only and All In One service
func GetDBMSCommonProvisionParametersSchema() *service.ParametersSchema {
	p := service.GetCommonProvisionParametersSchema()

	p.Properties["sslEnforcement"] = service.Parameter{
		Type: "string",
		Description: "Specifies whether the server requires the use of TLS " +
			"when connecting. Can be 'enabled', 'disabled' or ''. " +
			"Left unspecified, SSL will be enforced",
	}

	firewallRuleSchema := map[string]interface{}{}
	firewallRuleSchema["name"] = service.Parameter{
		Type:        "string",
		Description: "Name of firewall rule",
	}

	firewallRuleSchema["startIPAddress"] = service.Parameter{
		Type:        "string",
		Description: "Start of firewall rule range",
	}

	firewallRuleSchema["endIPAddress"] = service.Parameter{
		Type:        "string",
		Description: "End of firewall rule range",
	}

	p.Properties["firewallRules"] = service.Parameter{
		Type: "array",
		Description: "Firewall rules to apply to instance. " +
			"If left unspecified, defaults to only Azure IPs",
		Items: service.Parameter{
			Type:       "object",
			Properties: firewallRuleSchema,
		},
	}
	return p
}
