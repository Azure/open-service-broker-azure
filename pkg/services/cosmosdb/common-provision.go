package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

func generateAccountName(location string) string {
	databaseAccountName := uuid.NewV4().String()
	// CosmosDB currently limits database account names to 50 characters,
	// which includes location and a - character. Check if we will
	// exceed this and generate a shorter random identifier if needed.
	effectiveNameLength := len(location) + len(databaseAccountName)
	if effectiveNameLength > 49 {
		nameLength := 49 - len(location)
		databaseAccountName = generate.NewIdentifierOfLength(nameLength)
		logFields := log.Fields{
			"name":   databaseAccountName,
			"length": len(databaseAccountName),
		}
		log.WithFields(logFields).Debug(
			"returning fallback database account name",
		)
	}
	return databaseAccountName
}

func preProvision(
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := cosmosdbInstanceDetails{
		ARMDeploymentName:   uuid.NewV4().String(),
		DatabaseAccountName: generateAccountName(instance.Location),
	}
	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, nil, err
}
