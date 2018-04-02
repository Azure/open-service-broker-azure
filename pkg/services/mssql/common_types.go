package mssql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

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

func getDBMSCommonProvisionParamSchema() map[string]service.ParameterSchema {
	p := map[string]service.ParameterSchema{}

	sslEnforcementSchema := service.NewSimpleParameterSchema(
		"string",
		"Specifies whether the server requires the use of TLS"+
			" when connecting. Left unspecified, SSL will be enforced",
	)
	sslEnforcementSchema.SetRequired(true)
	sslEnforcementSchema.SetAllowedValues(
		[]string{"", "enabled", "disabled"},
	)
	sslEnforcementSchema.SetDefault("")

	firewallRuleSchema := map[string]service.ParameterSchema{}

	firewallRuleNameSchema := service.NewSimpleParameterSchema(
		"string",
		"Name of firewall rule",
	)
	firewallRuleNameSchema.SetRequired(true)
	firewallRuleSchema["name"] = firewallRuleNameSchema

	startIPSchema := service.NewSimpleParameterSchema(
		"string",
		"Start of firewall rule range",
	)
	startIPSchema.SetRequired(true)
	firewallRuleSchema["startIPAddress"] = startIPSchema

	endIPSchema := service.NewSimpleParameterSchema(
		"string",
		"End of firewall rule range",
	)
	endIPSchema.SetRequired(true)
	firewallRuleSchema["endIPAddress"] = endIPSchema

	firewallObject := service.NewObjectParameterSchema(
		"Individual Firewall Rule",
		firewallRuleSchema,
	)

	firewallRulesSchema := service.NewArrayParameterSchema(
		"Firewall rules to apply to instance. "+
			"If left unspecified, defaults to only Azure IPs",
		firewallObject,
	)
	p["firewallRules"] = firewallRulesSchema
	return p
}
