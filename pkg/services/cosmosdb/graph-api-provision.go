package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (g *graphAccountManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep(
			"preProvision", g.preProvision),
		service.NewProvisioningStep("deployARMTemplate", g.deployARMTemplate),
	)
}

func (g *graphAccountManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	dt := instance.Details.(*cosmosdbInstanceDetails)
	p, err := g.buildGoTemplateParams(
		pp,
		dt,
		"GlobalDocumentDB",
	)
	if err != nil {
		return nil, fmt.Errorf("error building arm params: %s", err)
	}
	p["capability"] = "EnableGremlin"
	tags := getTags(pp)
	tags["defaultExperience"] = "Graph"
	fqdn, pk, err := g.cosmosAccountManager.deployARMTemplate(
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
		fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s;",
			dt.FullyQualifiedDomainName,
			dt.PrimaryKey,
		),
	)
	return dt, err
}
