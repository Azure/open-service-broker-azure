package mssql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	td := instance.Plan.GetProperties().Extended["tierDetails"]
	details := td.(planDetails)
	return details.validateUpdateParameters(instance)
}

func (a *allInOneManager) GetUpdater(service.Plan) (service.Updater, error) {
	// There isn't a need to do any "pre-provision here. just the update step"
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", a.updateARMTemplate),
		service.NewUpdatingStep("updateConnectionPolicy", a.updateConnectionPolicy),
	)
}

func (a *allInOneManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	pp := instance.ProvisioningParameters
	up := instance.UpdatingParameters
	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParams, err := buildDBMSGoTemplateParameters(
		&dt.dbmsInstanceDetails,
		*up,
		version,
	)
	if err != nil {
		return nil, err
	}
	pd := instance.Plan.GetProperties().Extended["tierDetails"].(planDetails)
	dbParams, err := buildDatabaseGoTemplateParameters(
		dt.DatabaseName,
		*up,
		pd,
	)
	if err != nil {
		return nil, err
	}
	for key, value := range dbParams {
		goTemplateParams[key] = value
	}
	goTemplateParams["location"] = pp.GetString("location")
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	_, err = a.armDeployer.Update(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
		pp.GetString("location"),
		allInOneARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm template params
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	// This shouldn't change the instance details, so just return
	// what was there already
	return instance.Details, err
}

func (a *allInOneManager) updateConnectionPolicy(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	pp := instance.ProvisioningParameters
	up := instance.UpdatingParameters
	connectionPolicy := up.GetString("connectionPolicy")
	var err error
	if connectionPolicy != "" {
		err = setConnectionPolicy(
			ctx,
			&a.serverConnectionPoliciesClient,
			pp.GetString("resourceGroup"),
			dt.ServerName,
			pp.GetString("location"),
			connectionPolicy,
		)
	}
	return instance.Details, err
}
