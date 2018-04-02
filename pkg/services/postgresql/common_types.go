package postgresql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
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

	sslEnforcementSchema := service.NewParameterSchema(
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

	firewallRuleNameSchema := service.NewParameterSchema(
		"string",
		"Name of firewall rule",
	)
	firewallRuleNameSchema.SetRequired(true)
	firewallRuleSchema["name"] = firewallRuleNameSchema

	startIPSchema := service.NewParameterSchema(
		"string",
		"Start of firewall rule range",
	)
	startIPSchema.SetRequired(true)
	firewallRuleSchema["startIPAddress"] = startIPSchema

	endIPSchema := service.NewParameterSchema(
		"string",
		"End of firewall rule range",
	)
	endIPSchema.SetRequired(true)
	firewallRuleSchema["endIPAddress"] = endIPSchema

	firewallObject := service.NewParameterSchema(
		"object",
		"Individual Firewall Rule",
	)
	err := firewallObject.AddParameters(firewallRuleSchema)
	if err != nil {
		log.Errorf("error building firewallObject schema : %s", err)
	}
	firewallRulesSchema := service.NewParameterSchema(
		"arary",
		"Firewall rules to apply to instance. "+
			"If left unspecified, defaults to only Azure IPs",
	)
	err = firewallRulesSchema.SetItems(firewallObject)
	if err != nil {
		log.Errorf("error building firewallObject array schema : %s", err)
	}

	p["firewallRules"] = firewallRulesSchema
	return p
}
