package postgresql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func validateStorageUpdate(
	pp dbmsProvisioningParameters,
	up dbmsUpdatingParameters,
) error {
	var existingStorage int64
	if pp.Storage == nil {
		existingStorage = 10
	} else {
		existingStorage = *pp.Storage
	}
	var newStorge int64
	if up.Storage == nil {
		newStorge = 10
	} else {
		newStorge = *up.Storage
	}
	if newStorge < existingStorage {
		return service.NewValidationError(
			"storage",
			fmt.Sprintf(`invalid value: "%d". cannot reduce storage`, *up.Storage),
		)
	}
	return nil
}
