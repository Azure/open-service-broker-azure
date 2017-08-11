package azure

var locations = []string{
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

// IsValidLocation returns a bool indicating whether the provided location is a
// valid one
func IsValidLocation(location string) bool {
	for _, l := range locations {
		if location == l {
			return true
		}
	}
	return false
}
