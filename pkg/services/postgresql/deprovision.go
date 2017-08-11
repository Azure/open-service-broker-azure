package postgresql

import (
	"context"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
)

func (m *module) GetDeprovisioner(
	string,
	string,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", m.deleteARMDeployment),
		// krancour: This next step is a workaround because, currently, deleting
		// the ARM deployment is NOT deleting the PostgreSQL server. This seems to
		// be a problem not with ARM, but with the Postgres RP.
		service.NewDeprovisioningStep(
			"deletePostgreSQLServer",
			m.deletePostgreSQLServer,
		),
	)
}

func (m *module) deleteARMDeployment(
	ctx context.Context, // nolint: unparam
	provisioningContext interface{},
) (interface{}, error) {
	pc, ok := provisioningContext.(*postgresqlProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as postgresqlProvisioningContext",
		)
	}
	if err := m.armDeployer.Delete(
		pc.ARMDeploymentName,
		pc.ResourceGroupName,
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return pc, nil
}

func (m *module) deletePostgreSQLServer(
	ctx context.Context, // nolint: unparam
	provisioningContext interface{},
) (interface{}, error) {
	pc, ok := provisioningContext.(*postgresqlProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as postgresqlProvisioningContext",
		)
	}
	if err := m.postgresqlManager.DeleteServer(
		pc.ServerName,
		pc.ResourceGroupName,
	); err != nil {
		return nil, fmt.Errorf("error deleting postgresql server: %s", err)
	}
	return pc, nil
}
