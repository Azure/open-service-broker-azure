package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *sqlAllInOneManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
		service.NewProvisioningStep("createDatabase", s.createDatabase),
	)
}

func (s *sqlAllInOneManager) preProvision(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {

	cdt, err := s.cosmosAccountManager.preProvision(ctx, instance)
	if err != nil {
		return nil, err
	}
	aid := &sqlAllInOneInstanceDetails{
		cosmosdbInstanceDetails: *cdt.(*cosmosdbInstanceDetails),
		DatabaseName:            uuid.NewV4().String(),
	}
	return aid, nil
}

func (s *sqlAllInOneManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {

	pp := instance.ProvisioningParameters
	dt := instance.Details.(*sqlAllInOneInstanceDetails)
	p := s.cosmosAccountManager.buildGoTemplateParams(
		pp,
		&dt.cosmosdbInstanceDetails,
		"GlobalDocumentDB",
	)
	tags := getTags(pp)
	tags["defaultExperience"] = "DocumentDB"
	fqdn, pk, err := s.cosmosAccountManager.deployARMTemplate(
		pp,
		&dt.cosmosdbInstanceDetails,
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

func (s *sqlAllInOneManager) createDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*sqlAllInOneInstanceDetails)
	err := createDatabase(
		dt.DatabaseAccountName,
		dt.DatabaseName,
		string(dt.PrimaryKey),
	)
	if err != nil {
		return nil, err
	}
	return instance.Details, nil
}
