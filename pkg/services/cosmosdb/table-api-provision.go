package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
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
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {

	pp := instance.ProvisioningParameters
	dt := instance.Details.(*cosmosdbInstanceDetails)
	p := t.buildGoTemplateParams(pp, dt, "GlobalDocumentDB")
	p["capability"] = "EnableTable"
	tags := getTags(pp)
	tags["defaultExperience"] = "Table"
	fqdn, pk, err := t.cosmosAccountManager.deployARMTemplate(
		pp,
		dt,
		p,
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	dt.FullyQualifiedDomainName = fqdn
	dt.PrimaryKey = service.SecureString(pk)
	dt.ConnectionString = service.SecureString(
		fmt.Sprintf(
			"DefaultEndpointsProtocol=https;AccountName=%s;"+
				"AccountKey=%s;TableEndpoint=%s",
			dt.DatabaseAccountName,
			dt.FullyQualifiedDomainName,
			dt.PrimaryKey,
		),
	)
	return dt, err

}
