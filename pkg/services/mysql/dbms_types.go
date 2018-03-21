package mysql

import "github.com/Azure/open-service-broker-azure/pkg/service"

type dbmsProvisioningParameters struct {
	SSLEnforcement string         `json:"sslEnforcement"`
	FirewallRules  []firewallRule `json:"firewallRules"`
}

type firewallRule struct {
	Name    string `json:"name"`
	StartIP string `json:"startIPAddress"`
	EndIP   string `json:"endIPAddress"`
}

// GetDBMSProvisionParametersSchema returns the parameter schema for
// DBMS only service instances
func GetDBMSProvisionParametersSchema() *service.ParametersSchema {
	p := GetDBMSCommonProvisionParametersSchema()

	p.Properties["alias"] = service.Parameter{
		Type:        "string",
		Description: "Alias to use when provisioning databases on this DBMS",
	}
	p.Required = []string{"alias"}
	return p
}

type dbmsInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	ServerName               string `json:"server"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	EnforceSSL               bool   `json:"enforceSSL"`
}

type secureDBMSInstanceDetails struct {
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}

func (d *dbmsManager) SplitProvisioningParameters(
	cpp service.CombinedProvisioningParameters,
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
) {
	pp := dbmsProvisioningParameters{}
	if err := service.GetStructFromMap(cpp, &pp); err != nil {
		return nil, nil, err
	}
	ppMap, err := service.GetMapFromStruct(pp)
	return ppMap, nil, err
}

func (d *dbmsManager) SplitBindingParameters(
	params service.CombinedBindingParameters,
) (service.BindingParameters, service.SecureBindingParameters, error) {
	return nil, nil, nil
}
