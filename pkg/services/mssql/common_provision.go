// +build experimental

package mssql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func buildDBMSGoTemplateParameters(
	instance service.Instance,
) (map[string]interface{}, error) {
	dt := dbmsInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}
	sdt := secureDBMSInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, err
	}
	pp := dbmsProvisioningParams{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, &pp); err != nil {
		return nil, err
	}
	p := map[string]interface{}{}
	p["serverName"] = dt.ServerName
	p["administratorLogin"] = dt.AdministratorLogin
	p["administratorLoginPassword"] = sdt.AdministratorLoginPassword
	p["version"] = instance.Service.GetProperties().Extended["version"]
	// Only include these if they are not empty.
	// ARM Deployer will fail if the values included are not
	// valid IPV4 addresses (i.e. empty string wil fail)
	if len(pp.FirewallRules) > 0 {
		p["firewallRules"] = pp.FirewallRules
	} else {
		// Build the azure default
		p["firewallRules"] = []firewallRule{
			{
				Name:    "AllowAzure",
				StartIP: "0.0.0.0",
				EndIP:   "0.0.0.0",
			},
		}
	}
	return p, nil
}

func buildDatabaseGoTemplateParameters(
	instance service.Instance,
) (map[string]interface{}, error) {
	dt := databaseInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}
	tierDetails, err := buildTierDetailsGoTemplateParameters(instance)
	if err != nil {
		return nil, err
	}
	p := map[string]interface{}{}
	p["databaseName"] = dt.DatabaseName
	for key, value := range tierDetails {
		p[key] = value
	}
	return p, nil
}

func buildTierDetailsGoTemplateParameters(
	instance service.Instance,
) (map[string]interface{}, error) {
	var p map[string]interface{}
	var err error
	dt, ok := instance.Plan.GetProperties().Extended["tierDetails"]
	if ok {
		details := dt.(planDetails)
		p, err = details.getTierProvisionParameters(instance)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}
