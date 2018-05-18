package mysql

import (
	"context"
	"fmt"

	mysqlSDK "github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-04-30-preview/mysql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	uuid "github.com/satori/go.uuid"
)

const (
	enabledParamString  = "enabled"
	disabledParamString = "disabled"
	enabledARMString    = "Enabled"
	disabledARMString   = "Disabled"
)

func buildGoTemplateParameters(
	instance service.Instance,
) (map[string]interface{}, error) {

	plan := instance.Plan

	dt := dbmsInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}
	sdt := secureAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, err
	}
	pp := dbmsProvisioningParameters{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, &pp); err != nil {
		return nil, err
	}

	td := plan.GetProperties().Extended["tierDetails"].(tierDetails)

	p := map[string]interface{}{}
	p["sku"] = td.getSku(pp)
	p["tier"] = plan.GetProperties().Extended["tier"]
	p["cores"] = td.getCores(pp)
	p["storage"] = td.getStorage(pp) * 1024 //storage is in MB to arm :/
	p["backupRetention"] = td.getBackupRetention(pp)
	p["hardwareFamily"] = td.getHardwareFamily(pp)
	if td.isGeoRedundentBackup(pp) {
		p["geoRedundantBackup"] = enabledARMString
	}
	p["serverName"] = dt.ServerName
	p["administratorLoginPassword"] = sdt.AdministratorLoginPassword
	if td.isSSLRequired(pp) {
		p["sslEnforcement"] = enabledARMString
	} else {
		p["sslEnforcement"] = disabledARMString
	}
	p["version"] = instance.Service.GetProperties().Extended["version"]
	p["firewallRules"] = getFirewallRules(pp)

	return p, nil
}

func getAvailableServerName(
	ctx context.Context,
	checkNameAvailabilityClient mysqlSDK.CheckNameAvailabilityClient,
) (string, error) {
	for {
		serverName := uuid.NewV4().String()
		nameAvailability, err := checkNameAvailabilityClient.Execute(
			ctx,
			mysqlSDK.NameAvailabilityRequest{
				Name: &serverName,
			},
		)
		if err != nil {
			return "", fmt.Errorf(
				"error determining server name availability: %s",
				err,
			)
		}
		if *nameAvailability.NameAvailable {
			return serverName, nil
		}
	}
}
