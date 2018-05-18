package postgresql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/lib/pq" // Postgres SQL driver
)

func validateStorageUpdate(
	pp dbmsProvisioningParameters,
	up dbmsUpdatingParameters,
) error {
	if up.Storage != nil &&
		pp.Storage != nil &&
		*up.Storage < *pp.Storage {
		return service.NewValidationError(
			"storage",
			fmt.Sprintf(`invalid value: "%d". cannot reduce storage`, *up.Storage),
		)
	}
	return nil
}
