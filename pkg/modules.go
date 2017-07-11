package main

import (
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/echo"
)

var modules = []service.Module{
	echo.New(),
}
