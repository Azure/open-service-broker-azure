package mysql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func createBasicPlan(
	planID string,
) *service.PlanProperties {
	td := tierDetails{
		allowedHardware:         []string{gen4ParamString, gen5ParamString},
		allowedCores:            []int64{1, 2},
		defaultCores:            1,
		maxStorage:              1024,
		allowedBackupRedundancy: []string{"local"},
		maxBackupRetention:      35,
		defaultBackupRetention:  7,
		tier: "B",
	}

	return &service.PlanProperties{
		ID:   planID,
		Name: "basic",
		Description: "Basic Tier-- For workloads that require light compute and " +
			"I/O performance. Examples include servers used for development or " +
			"testing or small-scale infrequently used applications.",
		Free: false,
		Extended: map[string]interface{}{
			"tierDetails": td,
			"tier":        "Basic",
		},
		Metadata: &service.ServicePlanMetadata{
			DisplayName: "Basic Tier",
			Bullets:     []string{"Up to 2 vCores", "Variable I/O performance"},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: generateDBMSPlanSchema(td),
			},
		},
	}
}

func createGPPlan(
	planID string,
) *service.PlanProperties {

	td := tierDetails{
		allowedHardware:         []string{gen4ParamString, gen5ParamString},
		allowedCores:            []int64{2, 4, 8, 16, 32},
		defaultCores:            2,
		maxStorage:              2048,
		allowedBackupRedundancy: []string{"local", "geo"},
		maxBackupRetention:      35,
		defaultBackupRetention:  7,
		tier: "GP",
	}
	extendedPlanData := map[string]interface{}{
		"tierDetails": td,
		"tier":        "GeneralPurpose",
	}

	return &service.PlanProperties{
		ID:   planID,
		Name: "general-purpose",
		Description: "General Purpose Tier-- For most business workloads that " +
			"require balanced compute and memory with scalable I/O throughput. " +
			"Examples include servers for hosting web and mobile apps and other " +
			"enterprise applications.",
		Free:     false,
		Extended: extendedPlanData,
		Metadata: &service.ServicePlanMetadata{
			DisplayName: "General Purpose Tier",
			Bullets: []string{
				"Up to 32 vCores",
				"Predictable I/O Performance",
				"Local or Geo-Redundant Backups",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: generateDBMSPlanSchema(td),
			},
		},
	}
}

func createMemoryOptimizedPlan(
	planID string,
) *service.PlanProperties {

	td := tierDetails{
		allowedHardware:         []string{gen5ParamString},
		allowedCores:            []int64{2, 4, 8, 16},
		defaultCores:            2,
		maxStorage:              2048,
		allowedBackupRedundancy: []string{"local", "geo"},
		maxBackupRetention:      35,
		defaultBackupRetention:  7,
		tier: "MO",
	}
	extendedPlanData := map[string]interface{}{
		"tierDetails": td,
		"tier":        "MemoryOptimized",
	}

	return &service.PlanProperties{
		ID:   planID,
		Name: "memory-optimized",
		Description: "Memory Optimized Tier-- For high-performance database " +
			"workloads that require in-memory performance for faster transaction " +
			"processing and higher concurrency. Examples include servers for " +
			"processing real-time data and high-performance transactional or " +
			"analytical apps.",
		Free:     false,
		Extended: extendedPlanData,
		Metadata: &service.ServicePlanMetadata{
			DisplayName: "Memory Optimized Tier",
			Bullets: []string{
				"Up to 16 memory optimized vCores",
				"Predictable I/O Performance",
				"Local or Geo-Redundant Backups",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: generateDBMSPlanSchema(td),
			},
		},
	}
}

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "997b8372-8dac-40ac-ae65-758b4a5075a5",
				Name:        "azure-mysql-5-7",
				Description: "Azure Database for MySQL 5.7-- DBMS and single database (preview)",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "Azure Database for MySQL 5.7 (preview)",
					ImageURL:         "https://azure.microsoft.com/svghandler/mysql/?width=200",
					LongDescription:  "Azure Database for MySQL 5.7-- DBMS and single database (preview)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/mysql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "MySQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "5.7",
				},
			},
			m.allInOneServiceManager,
			service.NewPlan(createBasicPlan("1b093840-8e02-4e28-9aba-fa716757ec38")),
			service.NewPlan(createGPPlan("eae202c3-521c-46d1-a047-872dacf781fd")),
			service.NewPlan(createMemoryOptimizedPlan("129f06f6-cbf2-416e-a235-0fa6e081a07a")),
		),
		// dbms only service
		service.NewService(
			&service.ServiceProperties{
				ID:             "30e7b836-199d-4335-b83d-adc7d23a95c2",
				Name:           "azure-mysql-5-7-dbms",
				Description:    "Azure Database for MySQL 5.7-- DBMS only (preview)",
				ChildServiceID: "6704ae59-3eae-49e9-82b4-4cbcc00edf08",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "Azure Database for MySQL 5.7-- DBMS Only (preview)",
					ImageURL:         "https://azure.microsoft.com/svghandler/mysql/?width=200",
					LongDescription:  "Azure Database for MySQL 5.7-- DBMS only (preview)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/mysql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: false,
				Tags:     []string{"Azure", "MySQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "5.7",
				},
			},
			m.dbmsManager,
			service.NewPlan(createBasicPlan("20938530-cb42-48b2-93dc-ea0d3003a89f")),
			service.NewPlan(createGPPlan("3a00b95f-6acf-4bf9-8b01-52fe03a2d607")),
			service.NewPlan(createMemoryOptimizedPlan("b242a78f-9946-406a-af67-813c56341960")),
		),
		// database only service
		service.NewService(
			&service.ServiceProperties{
				ID:              "6704ae59-3eae-49e9-82b4-4cbcc00edf08",
				Name:            "azure-mysql-5-7-database",
				Description:     "Azure Database for MySQL 5.7-- database only (preview)",
				ParentServiceID: "30e7b836-199d-4335-b83d-adc7d23a95c2",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "Azure Database for MySQL 5.7-- Database Only (preview)",
					ImageURL:         "https://azure.microsoft.com/svghandler/mysql/?width=200",
					LongDescription:  "Azure Database for MySQL 5.7-- database only (preview)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/mysql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "MySQL", "Database"},
				Extended: map[string]interface{}{
					"version": "5.7",
				},
			},
			m.databaseManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "ec77bd04-2107-408e-8fde-8100c1ce1f46",
				Name:        "database",
				Description: "A new database added to an existing DBMS",
				Free:        false,
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Azure Database for MySQL-- Database Only",
				},
			}),
		),
	}), nil
}
