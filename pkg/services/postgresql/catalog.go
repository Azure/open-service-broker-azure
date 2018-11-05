package postgresql

import "github.com/Azure/open-service-broker-azure/pkg/service"

func createBasicPlan(
	planID string,
	includeDBParams bool,
	stability service.Stability,
) service.PlanProperties {
	td := tierDetails{
		tierName:                "Basic",
		tierShortName:           "B",
		allowedCores:            []int64{1, 2},
		defaultCores:            1,
		maxStorage:              1024,
		allowedBackupRedundancy: []string{"local"},
	}

	return service.PlanProperties{
		ID:   planID,
		Name: "basic",
		Description: "Basic Tier-- For workloads that require light compute and " +
			"I/O performance.",
		Free:      false,
		Stability: stability,
		Extended: map[string]interface{}{
			"tierDetails": td,
		},
		Metadata: service.ServicePlanMetadata{
			DisplayName: "Basic Tier",
			Bullets:     []string{"Up to 2 vCores", "Variable I/O performance"},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: generateProvisioningParamsSchema(
					td,
					includeDBParams,
				),
				UpdatingParametersSchema: generateUpdatingParamsSchema(td),
			},
		},
	}
}

func createGPPlan(
	planID string,
	includeDBParams bool,
	stability service.Stability,
) service.PlanProperties {

	td := tierDetails{
		tierName:                "GeneralPurpose",
		tierShortName:           "GP",
		allowedCores:            []int64{2, 4, 8, 16, 32},
		defaultCores:            2,
		maxStorage:              2048,
		allowedBackupRedundancy: []string{"local", "geo"},
	}

	extendedPlanData := map[string]interface{}{
		"tierDetails": td,
	}

	return service.PlanProperties{
		ID:   planID,
		Name: "general-purpose",
		Description: "General Purpose Tier-- For most business workloads that " +
			"require balanced compute and memory with scalable I/O throughput.",
		Free:      false,
		Stability: stability,
		Extended:  extendedPlanData,
		Metadata: service.ServicePlanMetadata{
			DisplayName: "General Purpose Tier",
			Bullets: []string{
				"Up to 32 vCores",
				"Predictable I/O Performance",
				"Local or Geo-Redundant Backups",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: generateProvisioningParamsSchema(
					td,
					includeDBParams,
				),
				UpdatingParametersSchema: generateUpdatingParamsSchema(td),
			},
		},
	}
}

func createMemoryOptimizedPlan(
	planID string,
	includeDBParams bool,
	stability service.Stability,
) service.PlanProperties {

	td := tierDetails{
		tierName:                "MemoryOptimized",
		tierShortName:           "MO",
		allowedCores:            []int64{2, 4, 8, 16},
		defaultCores:            2,
		maxStorage:              2048,
		allowedBackupRedundancy: []string{"local", "geo"},
	}

	extendedPlanData := map[string]interface{}{
		"tierDetails": td,
	}

	return service.PlanProperties{
		ID:   planID,
		Name: "memory-optimized",
		Description: "Memory Optimized Tier-- For high-performance database " +
			"workloads that require in-memory performance for faster transaction " +
			"processing and higher concurrency.",
		Free:      false,
		Stability: stability,
		Extended:  extendedPlanData,
		Metadata: service.ServicePlanMetadata{
			DisplayName: "Memory Optimized Tier",
			Bullets: []string{
				"Up to 16 memory optimized vCores",
				"Predictable I/O Performance",
				"Local or Geo-Redundant Backups",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: generateProvisioningParamsSchema(
					td,
					includeDBParams,
				),
				UpdatingParametersSchema: generateUpdatingParamsSchema(td),
			},
		},
	}
}

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		// all-in-one
		service.NewService(
			service.ServiceProperties{
				ID:          "b43b4bba-5741-4d98-a10b-17dc5cee0175",
				Name:        "azure-postgresql-9-6",
				Description: "Azure Database for PostgreSQL 9.6-- DBMS and single database",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL 9.6",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- DBMS and single database",
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
			service.NewPlan(createBasicPlan("09b398f8-f3c1-49ae-b726-459444e22460", true, service.StabilityStable)),
			service.NewPlan(createGPPlan("5807fb83-8065-4d91-a1f7-b4437657cd77", true, service.StabilityStable)),
			service.NewPlan(createMemoryOptimizedPlan("90f27532-0286-42e5-8e23-c3bb37191368", true, service.StabilityStable)),
		),
		// dbms only
		service.NewService(
			service.ServiceProperties{
				ID:             "d3f74b44-79bc-4d1e-bf7d-c247c2b851f9",
				Name:           "azure-postgresql-9-6-dbms",
				Description:    "Azure Database for PostgreSQL 9.6-- DBMS only",
				ChildServiceID: "25434f16-d762-41c7-bbdd-8045d7f74ca",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL 9.6-- DBMS Only",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- DBMS only",
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
			service.NewPlan(createBasicPlan("73191861-04b3-4d0b-a29b-429eb15a83d4", false, service.StabilityStable)),
			service.NewPlan(createGPPlan("4c6932e8-30ec-4af9-83d2-6e27286dbab3", false, service.StabilityStable)),
			service.NewPlan(createMemoryOptimizedPlan("057e64ea-41b5-4ed7-bf99-4867a332cfb7", false, service.StabilityStable)),
		),
		// database only
		service.NewService(
			service.ServiceProperties{
				ID:              "25434f16-d762-41c7-bbdd-8045d7f74ca6",
				Name:            "azure-postgresql-9-6-database",
				Description:     "Azure Database for PostgreSQL 9.6-- database only",
				ParentServiceID: "d3f74b44-79bc-4d1e-bf7d-c247c2b851f9",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL 9.6-- Database Only",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- database only",
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
			service.NewPlan(service.PlanProperties{
				ID:          "df6f5ef1-e602-406b-ba73-09c107d1e31b",
				Name:        "database",
				Description: "A new database added to an existing DBMS",
				Free:        false,
				Stability:   service.StabilityStable,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Azure Database for PostgreSQL-- Database Only",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.databaseManager.getProvisionParametersSchema(),
					},
				},
			}),
		),
		// all-in-one for version 10
		service.NewService(
			service.ServiceProperties{
				ID:          "32d3b4e0-e68f-4e96-93d4-35fd380f0874",
				Name:        "azure-postgresql-10",
				Description: "Azure Database for PostgreSQL 10-- DBMS and single database",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL 10",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- DBMS and single database",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/postgresql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "PostgreSQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "10",
				},
			},
			m.allInOneManager,
			service.NewPlan(createBasicPlan("e5b147b1-9abd-4444-adf4-e8777b5219ec", true, service.StabilityPreview)),
			service.NewPlan(createGPPlan("0f209a55-f166-4530-b5ad-26f81c598616", true, service.StabilityPreview)),
			service.NewPlan(createMemoryOptimizedPlan("6caf83ec-5cc1-42a0-9b34-0d163d73064c", true, service.StabilityPreview)),
		),
		// dbms only for version 10
		service.NewService(
			service.ServiceProperties{
				ID:             "cabd3125-5a13-46ea-afad-a69582af9578",
				Name:           "azure-postgresql-10-dbms",
				Description:    "Azure Database for PostgreSQL 10-- DBMS only",
				ChildServiceID: "1fd01042-3b70-4612-ac19-9ced0b2a1525",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL 10-- DBMS Only",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- DBMS only",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/postgresql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: false,
				Tags:     []string{"Azure", "PostgreSQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "10",
				},
			},
			m.dbmsManager,
			service.NewPlan(createBasicPlan("5cc758d2-b530-479e-8af7-e66f2906463a", false, service.StabilityPreview)),
			service.NewPlan(createGPPlan("f5218659-72ba-4fd3-9567-afd52d871fee", false, service.StabilityPreview)),
			service.NewPlan(createMemoryOptimizedPlan("c985dcc8-a4cd-43ac-a912-10793caed46b", false, service.StabilityPreview)),
		),
		// database only for version 10
		service.NewService(
			service.ServiceProperties{
				ID:              "1fd01042-3b70-4612-ac19-9ced0b2a1525",
				Name:            "azure-postgresql-10-database",
				Description:     "Azure Database for PostgreSQL 10-- database only",
				ParentServiceID: "cabd3125-5a13-46ea-afad-a69582af9578",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL 10-- Database Only",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- database only",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/postgresql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "PostgreSQL", "Database"},
				Extended: map[string]interface{}{
					"version": "10",
				},
			},
			m.databaseManager,
			service.NewPlan(service.PlanProperties{
				ID:          "672f80d5-8c9e-488f-b9ce-41142c205e99",
				Name:        "database",
				Description: "A new database added to an existing DBMS",
				Free:        false,
				Stability:   service.StabilityPreview,
				Metadata: service.ServicePlanMetadata{
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
