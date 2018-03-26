package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

type dbmsProvisioningParams struct {
	SSLEnforcement string         `json:"sslEnforcement"`
	FirewallRules  []firewallRule `json:"firewallRules"`
}

func (
	d *dbmsManager,
) getProvisionParametersSchema() map[string]*service.ParameterSchema {
	p := getDBMSCommonProvisionParamSchema()

	p["alias"] = &service.ParameterSchema{
		Type:        "string",
		Description: "Alias to use when provisioning databases on this DBMS",
		Required:    true,
	}
	return p
}

type firewallRule struct {
	Name    string `json:"name"`
	StartIP string `json:"startIPAddress"`
	EndIP   string `json:"endIPAddress"`
}

type dbmsInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	ServerName               string `json:"server"`
	AdministratorLogin       string `json:"administratorLogin"`
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
	pp := dbmsProvisioningParams{}
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
