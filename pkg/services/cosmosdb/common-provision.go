package cosmosdb

import (
	"errors"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

func generateDatabaseName(location string) string {
	databaseName := uuid.NewV4().String()
	// CosmosDB currently limits database name to 50 characters,
	// which includes location and a - character. Check if we will
	// exceed this and generate a shorter random identifier if needed.
	effectiveNameLength := len(location) + len(databaseName)
	if effectiveNameLength > 49 {
		nameLength := 49 - len(location)
		databaseName = generate.NewIdentifierOfLength(nameLength)
		logFields := log.Fields{
			"name":   databaseName,
			"length": len(databaseName),
		}
		log.WithFields(logFields).Debug("returning fallback database name")
	}
	return databaseName
}

func preProvision(
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*cosmosdbInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *cosmosdbInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.DatabaseAccountName = generateDatabaseName(instance.Location)
	return dt, instance.SecureDetails, nil
}
