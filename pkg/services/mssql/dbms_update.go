package mssql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (d *dbmsManager) GetUpdater(service.Plan) (service.Updater, error) {
	// There isn't a need to do any "pre-provision here. just the update step"
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", d.updateARMTemplate),
		service.NewUpdatingStep("updateConnectionPolicy", d.updateConnectionPolicy),
	)
}

func (d *dbmsManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsInstanceDetails)
	pp := instance.ProvisioningParameters
	up := instance.UpdatingParameters
	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParams, err := buildDBMSGoTemplateParameters(
		dt,
		*up,
		version,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	goTemplateParams["location"] = pp.GetString("location")
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	_, err = d.armDeployer.Update(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
		pp.GetString("location"),
		dbmsARMTemplateBytes,
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

func (d *dbmsManager) updateConnectionPolicy(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	up := instance.UpdatingParameters
	connectionPolicy := up.GetString("connectionPolicy")
	var err error
	if connectionPolicy != "" {
		err = setConnectionPolicy(
			ctx,
			&d.serverConnectionPoliciesClient,
			pp.GetString("resourceGroup"),
			pp.GetString("server"),
			connectionPolicy,
		)
	}
	return instance.Details, err
}
