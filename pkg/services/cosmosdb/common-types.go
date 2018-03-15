package cosmosdb

type databaseKind string

type cosmosdbInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	DatabaseAccountName      string `json:"name"`
	DatabaseKind             string `json:"kind"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
}

type cosmosdbSecureInstanceDetails struct {
	ConnectionString string `json:"connectionString"`
	PrimaryKey       string `json:"primaryKey"`
}
