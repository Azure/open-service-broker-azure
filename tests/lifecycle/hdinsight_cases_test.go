// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	hd "github.com/Azure/azure-service-broker/pkg/azure/hdinsight"
	"github.com/Azure/azure-service-broker/pkg/services/hdinsight"
)

func getHDInsightCases(
	armDeployer arm.Deployer,
) ([]moduleLifecycleTestCase, error) {
	hdinsightManager, err := hd.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{ // hadoop
			module:    hdinsight.New(armDeployer, hdinsightManager),
			serviceID: "c0fba6e1-4ce4-4d93-b751-c8c5e337739c",
			planID:    "fab80e5a-54c8-45e3-a466-f390de04e592",
			provisioningParameters: &hdinsight.ProvisioningParameters{
				Location:               "southcentralus",
				ResourceGroup:          newTestResourceGroupName(),
				ClusterWorkerNodeCount: 2,
			},
			bindingParameters: &hdinsight.BindingParameters{},
		},
		{ // spark
			module:    hdinsight.New(armDeployer, hdinsightManager),
			serviceID: "c0fba6e1-4ce4-4d93-b751-c8c5e337739c",
			planID:    "9815959a-35d2-4bf7-b467-3e77c03dcc3e",
			provisioningParameters: &hdinsight.ProvisioningParameters{
				Location:               "southcentralus",
				ResourceGroup:          newTestResourceGroupName(),
				ClusterWorkerNodeCount: 2,
			},
			bindingParameters: &hdinsight.BindingParameters{},
		},
		{ // hbase
			module:    hdinsight.New(armDeployer, hdinsightManager),
			serviceID: "c0fba6e1-4ce4-4d93-b751-c8c5e337739c",
			planID:    "609c6d56-851e-41cf-8a71-2dde705cf5a5",
			provisioningParameters: &hdinsight.ProvisioningParameters{
				Location:               "southcentralus",
				ResourceGroup:          newTestResourceGroupName(),
				ClusterWorkerNodeCount: 2,
			},
			bindingParameters: &hdinsight.BindingParameters{},
		},
		{ // storm
			module:    hdinsight.New(armDeployer, hdinsightManager),
			serviceID: "c0fba6e1-4ce4-4d93-b751-c8c5e337739c",
			planID:    "ebd2dcf7-c586-42b7-8eeb-06e5641a34aa",
			provisioningParameters: &hdinsight.ProvisioningParameters{
				Location:               "southcentralus",
				ResourceGroup:          newTestResourceGroupName(),
				ClusterWorkerNodeCount: 2,
			},
			bindingParameters: &hdinsight.BindingParameters{},
		},
	}, nil
}
