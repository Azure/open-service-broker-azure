package postgresql

import "github.com/Azure/open-service-broker-azure/pkg/service"

func createBasicPlan(
	planID string,
	includeDBParams bool,
) *service.PlanProperties {
	provisionSchema := planSchema{
		defaultFirewallRules: []firewallRule{
			{
				Name:    "AllowAzure",
				StartIP: "0.0.0.0",
				EndIP:   "0.0.0.0",
			},
		},
		allowedSSLEnforcement:   []string{enabledParamString, disabledParamString},
		defaultSSLEnforcement:   enabledParamString,
		allowedHardware:         []string{gen4ParamString, gen5ParamString},
		defaultHardware:         gen5ParamString,
		allowedCores:            []int64{1, 2},
		defaultCores:            1,
		maxStorage:              1024,
		minStorage:              5,
		defaultStorage:          10,
		allowedBackupRedundancy: []string{"local"},
		defaultBackupRedundancy: "local",
		minBackupRetention:      7,
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
			"provisionSchema": provisionSchema,
			"tier":            "Basic",
		},
		Metadata: &service.ServicePlanMetadata{
			DisplayName: "Basic Tier",
			Bullets:     []string{"Up to 2 vCores", "Variable I/O performance"},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: generateDBMSPlanSchema(
					provisionSchema,
					includeDBParams,
				),
			},
		},
	}
}

func createGPPlan(
	planID string,
	includeDBParams bool,
) *service.PlanProperties {

	provisionSchema := planSchema{
		allowedSSLEnforcement:   []string{enabledParamString, disabledParamString},
		defaultSSLEnforcement:   enabledParamString,
		allowedHardware:         []string{gen4ParamString, gen5ParamString},
		defaultHardware:         gen5ParamString,
		allowedCores:            []int64{2, 4, 8, 16, 32},
		defaultCores:            2,
		maxStorage:              2048,
		minStorage:              5,
		defaultStorage:          10,
		allowedBackupRedundancy: []string{"local", "geo"},
		defaultBackupRedundancy: "local",
		minBackupRetention:      7,
		maxBackupRetention:      35,
		defaultBackupRetention:  7,
		tier: "GP",
	}
	extendedPlanData := map[string]interface{}{
		"provisionSchema": provisionSchema,
		"tier":            "GeneralPurpose",
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
				ProvisioningParametersSchema: generateDBMSPlanSchema(
					provisionSchema,
					includeDBParams,
				),
			},
		},
	}
}

func createMemoryOptimizedPlan(
	planID string,
	includeDBParams bool,
) *service.PlanProperties {

	provisionSchema := planSchema{
		allowedSSLEnforcement:   []string{enabledParamString, disabledParamString},
		defaultSSLEnforcement:   enabledParamString,
		allowedHardware:         []string{gen5ParamString},
		defaultHardware:         gen5ParamString,
		allowedCores:            []int64{2, 4, 8, 16},
		defaultCores:            2,
		maxStorage:              2048,
		minStorage:              5,
		defaultStorage:          10,
		allowedBackupRedundancy: []string{"local", "geo"},
		defaultBackupRedundancy: "local",
		minBackupRetention:      7,
		maxBackupRetention:      35,
		defaultBackupRetention:  7,
		tier: "MO",
	}
	extendedPlanData := map[string]interface{}{
		"provisionSchema": provisionSchema,
		"tier":            "MemoryOptimized",
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
				ProvisioningParametersSchema: generateDBMSPlanSchema(
					provisionSchema,
					includeDBParams,
				),
			},
		},
	}
}

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		// all-in-one
		service.NewService(
			&service.ServiceProperties{
				ID:          "b43b4bba-5741-4d98-a10b-17dc5cee0175",
				Name:        "azure-postgresql-9-6",
				Description: "Azure Database for PostgreSQL 9.6-- DBMS and single database (preview)",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL 9.6 (preview)",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- DBMS and single database (preview)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/postgresql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "PostgreSQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "9.6",
				},
			},
			m.allInOneManager,
			service.NewPlan(createBasicPlan("09b398f8-f3c1-49ae-b726-459444e22460", true)),
			service.NewPlan(createGPPlan("5807fb83-8065-4d91-a1f7-b4437657cd77", true)),
			service.NewPlan(createMemoryOptimizedPlan("90f27532-0286-42e5-8e23-c3bb37191368", true)),
		),
		// dbms only
		service.NewService(
			&service.ServiceProperties{
				ID:             "d3f74b44-79bc-4d1e-bf7d-c247c2b851f9",
				Name:           "azure-postgresql-9-6-dbms",
				Description:    "Azure Database for PostgreSQL 9.6-- DBMS only (preview)",
				ChildServiceID: "25434f16-d762-41c7-bbdd-8045d7f74ca",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL 9.6-- DBMS Only (preview)",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- DBMS only (preview)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/postgresql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: false,
				Tags:     []string{"Azure", "PostgreSQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "9.6",
				},
			},
			m.dbmsManager,
			service.NewPlan(createBasicPlan("73191861-04b3-4d0b-a29b-429eb15a83d4", false)),
			service.NewPlan(createGPPlan("4c6932e8-30ec-4af9-83d2-6e27286dbab3", false)),
			service.NewPlan(createMemoryOptimizedPlan("057e64ea-41b5-4ed7-bf99-4867a332cfb7", false)),
		),
		// database only
		service.NewService(
			&service.ServiceProperties{
				ID:              "25434f16-d762-41c7-bbdd-8045d7f74ca6",
				Name:            "azure-postgresql-9-6-database",
				Description:     "Azure Database for PostgreSQL 9.6-- database only (preview)",
				ParentServiceID: "d3f74b44-79bc-4d1e-bf7d-c247c2b851f9",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL 9.6-- Database Only (preview)",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- database only (preview)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/postgresql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "PostgreSQL", "Database"},
				Extended: map[string]interface{}{
					"version": "9.6",
				},
			},
			m.databaseManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "df6f5ef1-e602-406b-ba73-09c107d1e31b",
				Name:        "database",
				Description: "A new database added to an existing DBMS (preview)",
				Free:        false,
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Azure Database for PostgreSQL-- Database Only",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.databaseManager.getProvisionParametersSchema(),
					},
				},
			}),
		),
	}), nil
}
