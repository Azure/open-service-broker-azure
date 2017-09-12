package hdinsight

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	"github.com/Azure/azure-service-broker/pkg/azure/hdinsight"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/hdinsight/armtemplate"
)

type module struct {
	armDeployer      arm.Deployer
	hdinsightManager hdinsight.Manager
}

var armTemplateBytes = map[string][]byte{
	"Hadoop": armtemplate.Hadoop(),
	"HBase":  armtemplate.HBase(),
	"Spark":  armtemplate.Spark(),
	"Storm":  armtemplate.Storm(),
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning HDInsight cluster and an additional
// storage account
func New(
	armDeployer arm.Deployer,
	hdinsightManager hdinsight.Manager,
) service.Module {
	return &module{
		armDeployer:      armDeployer,
		hdinsightManager: hdinsightManager,
	}
}

func (m *module) GetName() string {
	return "hdinsight"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityAlpha
}
