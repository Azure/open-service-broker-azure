package hdinsight

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "c0fba6e1-4ce4-4d93-b751-c8c5e337739c",
				Name:        "azure-hdinsight",
				Description: "Azure HDInisght Service",
				Bindable:    true,
				Tags:        []string{"Azure", "HDInsight", "Hadoop", "Spark", "Hbase"},
			},
			m.serviceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "fab80e5a-54c8-45e3-a466-f390de04e592",
				Name:        "Hadoop",
				Description: "Apache Hadoop: Uses HDFS, YARN resource management, and a simple MapReduce programming model to process and analyze batch data in parallel. Head Nodes: D12v2 * 2, Worker nodes: D4v2 * 4 by default, Zookeeper nodes: A1 * 3.", // nolint: lll
				Free:        false,
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "9815959a-35d2-4bf7-b467-3e77c03dcc3e",
				Name:        "Spark",
				Description: "Apache Spark: A parallel processing framework that supports in-memory processing to boost the performance of big-data analysis applications, Spark works for SQL, streaming data, and machine learning. Head Nodes: D12v2 * 2, Worker nodes: D4v2 * 4 by default.", // nolint: lll
				Free:        false,
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "609c6d56-851e-41cf-8a71-2dde705cf5a5",
				Name:        "HBase",
				Description: "Apache HBase: A NoSQL database built on Hadoop that provides random access and strong consistency for large amounts of unstructured and semi-structured data - potentially billions of rows times millions of columns. Head Nodes: D12v2 * 2, Worker nodes: D4v2 * 4 by default, Zookeeper nodes: A3 * 3.", // nolint: lll
				Free:        false,
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "ebd2dcf7-c586-42b7-8eeb-06e5641a34aa",
				Name:        "Storm",
				Description: "Apache Storm: A distributed, real-time computation system for processing large streams of data fast. Storm is offered as a managed cluster in HDInsight. Head Nodes: A3 * 2, Worker nodes: D3v2 * 4 by default, Zookeeper nodes: A3 * 3.", // nolint: lll
				Free:        false,
			}),
			// nolint: lll
			// Cluster type introduction link: https://docs.microsoft.com/en-us/azure/hdinsight/hdinsight-hadoop-introduction#overview
			// No RServer/Kafka plan. It seems that there is no suitable scenario for
			// an app to use them. We can add the plan if someone shares his scenario
			// to us.
		),
	}), nil
}
