package postgresql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func validateStorageUpdate(instance service.Instance) error {
	oldStorageAmt := instance.ProvisioningParameters.GetInt64("storage")
	newStorageAmt := instance.UpdatingParameters.GetInt64("storage")
	fmt.Printf("----> old storage amount: %d", oldStorageAmt)
	fmt.Printf("----> new storage amount: %d", newStorageAmt)
	if oldStorageAmt > newStorageAmt {
		return service.NewValidationError(
			"storage",
			fmt.Sprintf(
				`invalid value: cannot reduce storage from %d to %d`,
				oldStorageAmt,
				newStorageAmt,
			),
		)
	}
	return nil
}
