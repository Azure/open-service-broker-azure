package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlAccountManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
	)
}

func (s *sqlAccountManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {

	pp := instance.ProvisioningParameters
	dt := instance.Details.(*cosmosdbInstanceDetails)
	p := s.buildGoTemplateParams(pp, dt, "GlobalDocumentDB")
	tags := getTags(pp)
	tags["defaultExperience"] = "DocumentDB"

	fqdn, pk, err := s.cosmosAccountManager.deployARMTemplate(
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
