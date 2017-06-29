package main

import (
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/echo"
)

func getModules() []service.Module {
	return []service.Module{
		echo.New(),
	}
}
