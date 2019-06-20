package azure

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

var azurePublicCloudLocations = []string{
	"australiaeast",
	"australiasoutheast",
	"brazilsouth",
	"canadacentral",
	"canadaeast",
	"centralindia",
	"centralus",
	"eastasia",
	"eastus",
	"eastus2",
	"japaneast",
	"japanwest",
	"koreacentral",
	"koreasouth",
	"northcentralus",
	"northeurope",
	"southcentralus",
	"southeastasia",
	"southindia",
	"uksouth",
	"ukwest",
	"westcentralus",
	"westeurope",
	"westindia",
	"westus",
	"westus2",
}

var azureChinaCloudLocations = []string{
	"chinanorth2",
	"chinanorth",
	"chinaeast2",
	"chinaeast",
}

// IsValidLocation returns a bool indicating whether the provided location is a
// valid one
func IsValidLocation(location string) bool {
	envrionmentName := GetEnvrionmentName()
	var locations []string
	switch envrionmentName {
	case "AzureChinaCloud":
		locations = azureChinaCloudLocations
	case "AzurePublicCloud":
		locations = azurePublicCloudLocations
	// We shouldn't run into default case, but instead
	// of raising a panic, we use public cloud locations
	// and the error can be reported when provisioning the resource.
	default:
		locations = azurePublicCloudLocations
	}
	for _, l := range locations {
		if location == l {
			return true
		}
	}
	return false
}

// LocationValidator is a custom schema validator that validates a specified
// location is a real Azure region.
func LocationValidator(context, val string) error {
	if !IsValidLocation(val) {
		return service.NewValidationError(
			context,
			fmt.Sprintf(`invalid location: "%s"`, val),
		)
	}
	return nil
}
