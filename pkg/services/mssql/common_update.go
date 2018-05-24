package mssql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func buildDatabaseUpdateGoTemplateParameters(
	instance service.Instance,
) (map[string]interface{}, error) {
	dt := databaseInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}
	tierDetails, err := buildUpdateTierDetailsGoTemplateParameters(instance)
	if err != nil {
		return nil, err
	}
	p := map[string]interface{}{}
	p["databaseName"] = dt.DatabaseName
	for key, value := range tierDetails {
		p[key] = value
	}
	return p, nil
}

func buildUpdateTierDetailsGoTemplateParameters(
	instance service.Instance,
) (map[string]interface{}, error) {
	var p map[string]interface{}
	var err error
	dt, ok := instance.Plan.GetProperties().Extended["tierDetails"]
	if ok {
		details := dt.(planDetails)
		p, err = details.getTierUpdateParameters(instance)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func validateStorageUpdate(
	pp service.ProvisioningParameters,
	up service.ProvisioningParameters,
) error {
	existingStorage := pp.GetInt64("storage")
	newStorge := up.GetInt64("storage")
	if newStorge < existingStorage {
		return service.NewValidationError(
			"storage",
			fmt.Sprintf(
				`invalid value: cannot reduce storage from %d to %d`,
				existingStorage,
				newStorge,
			),
		)
	}
	return nil
}
