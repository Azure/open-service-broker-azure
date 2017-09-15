package mssql

import (
	"encoding/json"
	"os"
)

// GetConfig parses configuration details needed for connecting to Azure SQL
// Servers from environment variables and returns a Config object that
// encapsulates those details
func GetConfig() (Config, error) {
	c := Config{}
	scMap := make(map[string]ServerConfig)

	scString := os.Getenv("AZURE_SQL_SERVERS")
	scArray := []ServerConfig{}
	if scString != "" {
		if err := json.Unmarshal([]byte(scString), &scArray); err != nil {
			return c, err
		}
		for _, s := range scArray {
			scMap[s.ServerName] = s
		}
	}

	c.Servers = scMap
	return c, nil
}
