// +build experimental

package cosmosdb

import (
	"context"
	"fmt"

	"open-service-broker-azure/pkg/service"
)

func (t *tableAccountManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", t.preProvision),
		service.NewProvisioningStep("deployARMTemplate", t.deployARMTemplate),
	)
}

func (t *tableAccountManager) deployARMTemplate(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {

	pp := &provisioningParameters{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, pp); err != nil {
		return nil, nil, err
	}

	dt := &cosmosdbInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}

	p, err := t.buildGoTemplateParams(instance, "GlobalDocumentDB")
	if err != nil {
		return nil, nil, err
	}
	p["capability"] = "EnableTable"
	if instance.Tags == nil {
		instance.Tags = make(map[string]string)
	}
	instance.Tags["defaultExperience"] = "Table"

	fqdn, sdt, err := t.cosmosAccountManager.deployARMTemplate(
		ctx,
		instance,
		p,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	dt.FullyQualifiedDomainName = fqdn
	sdt.ConnectionString = fmt.Sprintf(
		"DefaultEndpointsProtocol=https;AccountName=%s;"+
			"AccountKey=%s;TableEndpoint=%s",
		dt.DatabaseAccountName,
		dt.FullyQualifiedDomainName,
		sdt.PrimaryKey,
	)

	dtMap, err := service.GetMapFromStruct(dt)
	if err != nil {
		return nil, nil, err
	}
	sdtMap, err := service.GetMapFromStruct(sdt)
	return dtMap, sdtMap, err

}
