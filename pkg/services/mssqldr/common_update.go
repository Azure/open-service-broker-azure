package mssqldr

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func updateDatabaseARMTemplate(
	armDeployer *arm.Deployer,
	armDeploymentName string,
	resourceGroup string,
	location string,
	serverName string,
	databaseName string,
	up service.ProvisioningParameters,
	pd planDetails,
	tags map[string]string,
) error {
	goTemplateParams, err := buildDatabaseGoTemplateParameters(
		databaseName,
		up,
		pd,
	)
	if err != nil {
		return err
	}
	goTemplateParams["location"] = location
	goTemplateParams["serverName"] = serverName
	_, err = (*armDeployer).Update(
		armDeploymentName,
		resourceGroup,
		location,
		databaseARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm params
		tags,
	)
	return err
}
