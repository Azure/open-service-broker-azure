// +build experimental

package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlAllInOneManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", s.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteDatabase", s.deleteDatabase,
		),
		service.NewDeprovisioningStep(
			"deleteCosmosDBServer", s.deleteCosmosDBAccount,
		),
	)
}

func (s *sqlAllInOneManager) deleteDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := &sqlAllInOneInstanceDetails{}
	err := service.GetStructFromMap(instance.Details, &dt)
	if err != nil {
		fmt.Printf("Failed to get DT Struct from Map %s", err)
		return nil, nil, err
	}

	sdt := &cosmosdbSecureInstanceDetails{}
	err = service.GetStructFromMap(instance.SecureDetails, &sdt)
	if err != nil {
		fmt.Printf("Failed to get SDT Struct from Map %s", err)
		return nil, nil, err
	}
	err = deleteDatabase(dt.DatabaseAccountName, dt.DatabaseName, sdt.PrimaryKey)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}
