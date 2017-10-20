package hdinsight

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/azure/hdinsight"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/hdinsight/armtemplate"
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
	"Kafka":  armtemplate.Kafka(),
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

func (s *serviceManager) GetName() string {
	return "hdinsight"
}

func (s *serviceManager) GetStability() service.Stability {
	return service.StabilityAlpha
}
