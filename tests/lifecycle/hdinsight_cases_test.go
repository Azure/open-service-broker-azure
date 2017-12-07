// +build !unit

package lifecycle

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	hd "github.com/Azure/open-service-broker-azure/pkg/azure/hdinsight"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/hdinsight"
)

func getHDInsightCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]moduleLifecycleTestCase, error) {
	hdinsightManager, err := hd.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{ // hadoop
			module:      hdinsight.New(armDeployer, hdinsightManager),
			description: "Hadoop",
			serviceID:   "c0fba6e1-4ce4-4d93-b751-c8c5e337739c",
			planID:      "fab80e5a-54c8-45e3-a466-f390de04e592",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &hdinsight.ProvisioningParameters{
				ClusterWorkerNodeCount: 1,
			},
			bindingParameters: &hdinsight.BindingParameters{},
		},
		{ // spark
			module:      hdinsight.New(armDeployer, hdinsightManager),
			description: "Spark",
			serviceID:   "c0fba6e1-4ce4-4d93-b751-c8c5e337739c",
			planID:      "9815959a-35d2-4bf7-b467-3e77c03dcc3e",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "eastus",
			},
			provisioningParameters: &hdinsight.ProvisioningParameters{
				ClusterWorkerNodeCount: 1,
			},
			bindingParameters: &hdinsight.BindingParameters{},
		},
		{ // hbase
			module:      hdinsight.New(armDeployer, hdinsightManager),
			description: "HBase",
			serviceID:   "c0fba6e1-4ce4-4d93-b751-c8c5e337739c",
			planID:      "609c6d56-851e-41cf-8a71-2dde705cf5a5",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "eastus2",
			},
			provisioningParameters: &hdinsight.ProvisioningParameters{
				ClusterWorkerNodeCount: 1,
			},
			bindingParameters: &hdinsight.BindingParameters{},
		},
		{ // storm
			module:      hdinsight.New(armDeployer, hdinsightManager),
			description: "Storm",
			serviceID:   "c0fba6e1-4ce4-4d93-b751-c8c5e337739c",
			planID:      "ebd2dcf7-c586-42b7-8eeb-06e5641a34aa",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "westus",
			},
			provisioningParameters: &hdinsight.ProvisioningParameters{
				ClusterWorkerNodeCount: 1,
			},
			bindingParameters: &hdinsight.BindingParameters{},
		},
		{ // kafka
			module:      hdinsight.New(armDeployer, hdinsightManager),
			description: "Kfaka",
			serviceID:   "c0fba6e1-4ce4-4d93-b751-c8c5e337739c",
			planID:      "c5f8277b-0cb1-4cfe-863d-03054493368a",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "northcentralus",
			},
			provisioningParameters: &hdinsight.ProvisioningParameters{
				ClusterWorkerNodeCount: 1,
			},
			bindingParameters: &hdinsight.BindingParameters{},
		},
	}, nil
}
