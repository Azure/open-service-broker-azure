package mssqldr

import "github.com/Azure/open-service-broker-azure/pkg/service"

func buildBasicPlan(
	id string,
) service.PlanProperties {

	planDetails := dtuPlanDetails{
		storageInGB: 2,
		defaultDTUs: 5,
		tierName:    "Basic",
		skuMap: map[int64]string{
			5: "Basic",
		},
	}

	planProperties := service.PlanProperties{
		ID:          id,
		Name:        "basic",
		Description: "Basic Tier, 5 DTUs, 2GB Storage, 7 days point-in-time restore",
		Free:        false,
		Stability:   service.StabilityExperimental,
		Extended: map[string]interface{}{
			"tierDetails": planDetails,
		},
		Metadata: service.ServicePlanMetadata{
			DisplayName: "Basic Tier",
			Bullets: []string{
				"5 DTUs",
				"Includes 2GB Storage",
				"7 days point-in-time restore",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: planDetails.getProvisionSchema(),
				UpdatingParametersSchema:     planDetails.getUpdateSchema(),
			},
		},
	}

	return planProperties
}

func buildStandardPlan(
	id string,
) service.PlanProperties {
	planDetails := dtuPlanDetails{
		storageInGB: 250,
		allowedDTUs: []int64{
			10, 20, 50, 100, 200, 400, 800, 1600, 3000,
		},
		defaultDTUs: 10,
		tierName:    "Standard",
		skuMap: map[int64]string{
			10:   "S0",
			20:   "S1",
			50:   "S2",
			100:  "S3",
			200:  "S4",
			400:  "S6",
			800:  "S7",
			1600: "S9",
			3000: "S12",
		},
	}

	planProperties := service.PlanProperties{
		ID:   id,
		Name: "standard",
		Description: "Standard Tier, Up to 3000 DTUs, 250GB Storage, " +
			"35 days point-in-time restore",
		Free:      false,
		Stability: service.StabilityExperimental,
		Extended: map[string]interface{}{
			"tierDetails": planDetails,
		},
		Metadata: service.ServicePlanMetadata{
			DisplayName: "Standard Tier",
			Bullets: []string{
				"10-3000 DTUs",
				"250GB",
				"35 days point-in-time restore",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: planDetails.getProvisionSchema(),
				UpdatingParametersSchema:     planDetails.getUpdateSchema(),
			},
		},
	}

	return planProperties
}

func buildPremiumPlan(
	id string,
) service.PlanProperties {
	planDetails := dtuPlanDetails{
		storageInGB: 500,
		allowedDTUs: []int64{
			125, 250, 500, 1000, 1750, 4000,
		},
		defaultDTUs: 125,
		tierName:    "Premium",
		skuMap: map[int64]string{
			125:  "P1",
			250:  "P2",
			500:  "P4",
			1000: "P6",
			1750: "P11",
			4000: "P15",
		},
	}

	planProperties := service.PlanProperties{
		ID:   id,
		Name: "premium",
		Description: "Premium Tier, Up to 4000 DTUs, 500GB Storage, " +
			"35 days point-in-time restore",
		Free:      false,
		Stability: service.StabilityExperimental,
		Extended: map[string]interface{}{
			"tierDetails": planDetails,
		},
		Metadata: service.ServicePlanMetadata{
			DisplayName: "Premium Tier",
			Bullets: []string{
				"Up to 4000 DTUs",
				"Includes 500GB Storage",
				"35 days point-in-time restore",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: planDetails.getProvisionSchema(),
				UpdatingParametersSchema:     planDetails.getUpdateSchema(),
			},
		},
	}

	return planProperties
}

func buildGeneralPurposePlan(
	id string,
) service.PlanProperties {
	planDetails := vCorePlanDetails{
		tierName:      "GeneralPurpose",
		tierShortName: "GP",
	}
	planProperties := service.PlanProperties{
		ID:          id,
		Name:        "general-purpose",
		Description: "Up to 80 vCores, 440 GB memory and 1 TB of storage (preview)",
		Free:        false,
		Stability:   service.StabilityExperimental,
		Extended: map[string]interface{}{
			"tierDetails": planDetails,
		},
		Metadata: service.ServicePlanMetadata{
			DisplayName: "General Purpose (preview)",
			Bullets: []string{
				"Scalable compute and storage options for budget-oriented applications",
				"Up to 80 vCores",
				"Up to 440 GB memory",
				"$187.62 / vCore",
				"7 days point-in-time restore",
				"Currently In Preview",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: planDetails.getProvisionSchema(),
				UpdatingParametersSchema:     planDetails.getUpdateSchema(),
			},
		},
	}

	return planProperties
}

func buildBusinessCriticalPlan(
	id string,
) service.PlanProperties {
	planDetails := vCorePlanDetails{
		tierName:      "BusinessCritical",
		tierShortName: "BC",
	}
	planProperties := service.PlanProperties{
		ID:   id,
		Name: "business-critical",
		Description: "Up to 80 vCores, 440 GB memory and 1 TB of storage. " +
			"Local SSD, highest resilience to failures. (preview)",
		Free:      false,
		Stability: service.StabilityExperimental,
		Extended: map[string]interface{}{
			"tierDetails": planDetails,
		},
		Metadata: service.ServicePlanMetadata{
			DisplayName: "Business Critical (preview)",
			Bullets: []string{
				"Up to 80 vCores",
				"Up to 440 GB memory",
				"$505.50 / vCore",
				"7 days point-in-time restore",
				"Currently In Preview",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: planDetails.getProvisionSchema(),
				UpdatingParametersSchema:     planDetails.getUpdateSchema(),
			},
		},
	}

	return planProperties
}

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {

	return service.NewCatalog([]service.Service{
		// dbms pair registered service
		service.NewService(
			service.ServiceProperties{
				ID:             "00ce53a3-d6c3-4c24-8cb2-3f48d3b161d8",
				Name:           "azure-sql-12-0-dr-dbms-pair-registered",
				Description:    "Azure SQL 12.0-- disaster recovery DBMS pair registered",
				ChildServiceID: "2eb94a7e-5a7c-46f9-b9d2-ff769f215845", // More children in fact
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure SQL 12.0-- disaster recovery DBMS Pair registered",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/SQLDatabase.svg",
					LongDescription:  "Azure SQL 12.0-- disaster recovery DBMS pair registered, as the primary server and the secondary server of failover groups",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: false,
				Tags: []string{
					"Azure",
					"SQL",
					"DBMS",
					"Server",
					service.DRTag,
					"Failover Group",
				},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.dbmsPairRegisteredManager,
			service.NewPlan(service.PlanProperties{
				ID:          "5683ca92-372b-49a6-b7cd-96a14645ec15",
				Name:        "dbms",
				Description: "Azure SQL Server-- DBMS",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Azure SQL Server-- DBMS",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.dbmsPairRegisteredManager.getProvisionParametersSchema(),
						UpdatingParametersSchema:     m.dbmsPairRegisteredManager.getUpdatingParametersSchema(),
					},
				},
			}),
		),
		// database pair service
		service.NewService(
			service.ServiceProperties{
				ID:              "2eb94a7e-5a7c-46f9-b9d2-ff769f215845",
				Name:            "azure-sql-12-0-dr-database-pair",
				Description:     "Azure SQL 12.0-- disaster recovery database pair",
				Bindable:        true,
				ParentServiceID: "00ce53a3-d6c3-4c24-8cb2-3f48d3b161d8",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure SQL 12.0-- disaster recovery Database Pair",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/SQLDatabase.svg",
					LongDescription:  "Azure SQL 12.0--disaster recovery database pair, create the primary database, the secondary database, and the failover group",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Tags: []string{
					"Azure",
					"SQL",
					"Database",
					service.DRTag,
					"Failover Group",
				},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.databasePairManager,
			service.NewPlan(
				buildBasicPlan(
					"07194a34-f2c8-4f01-aa39-84bbfc4cab73",
				),
			),
			service.NewPlan(
				buildStandardPlan(
					"edce3e74-69eb-4524-aabb-f2c4a7ee9398",
				),
			),
			service.NewPlan(
				buildPremiumPlan(
					"e67ab6df-20c4-4e82-86a2-de80278baa99",
				),
			),
			service.NewPlan(
				buildGeneralPurposePlan(
					"8cca0294-1eb2-46fd-9449-4c6a73cae3c0",
				),
			),
			service.NewPlan(
				buildBusinessCriticalPlan(
					"d19b19a5-0dec-4e1d-a443-3b5fba173010",
				),
			),
		),
		// database pair registered service
		service.NewService(
			service.ServiceProperties{
				ID:              "8480271a-f4c7-4232-b2b7-7f33391728f7",
				Name:            "azure-sql-12-0-dr-database-pair-registered",
				Description:     "Azure SQL 12.0-- disaster recovery database pair registered",
				Bindable:        true,
				ParentServiceID: "00ce53a3-d6c3-4c24-8cb2-3f48d3b161d8",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure SQL 12.0-- disaster recovery Database Pair registered",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/SQLDatabase.svg",
					LongDescription:  "Azure SQL 12.0-- disaster recovery database pair registered, the primary database, the secondary database, and the failover group are existing",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Tags: []string{
					"Azure",
					"SQL",
					"Database",
					service.DRTag,
					"Failover Group",
				},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.databasePairRegisteredManager,
			service.NewPlan(service.PlanProperties{
				ID:          "9e05f8b7-27ce-4fb4-b889-e7b2f8575df7",
				Name:        "database",
				Description: "Azure SQL Server-- database",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Azure SQL Server-- database",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.databasePairRegisteredManager.getProvisionParametersSchema(),
					},
				},
			}),
		),
		// database pair from existing primary service
		service.NewService(
			service.ServiceProperties{
				ID:              "505ae87a-5cd8-4aeb-b7ea-809dd249dc1f",
				Name:            "azure-sql-12-0-dr-database-pair-from-existing-primary",
				Description:     "Azure SQL 12.0-- disaster recovery database pair from existing primary",
				Bindable:        true,
				ParentServiceID: "00ce53a3-d6c3-4c24-8cb2-3f48d3b161d8",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure SQL 12.0-- disaster recovery Database Pair from existing primary",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/SQLDatabase.svg",
					LongDescription:  "Azure SQL 12.0-- disaster recovery database pair from existing primary database, create the secondary database and the failover group",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Tags: []string{
					"Azure",
					"SQL",
					"Database",
					service.DRTag,
					service.MigrationTag,
					"Failover Group",
				},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.databasePairManagerForExistingPrimary,
			service.NewPlan(
				buildBasicPlan(
					"8a65de90-6d8b-4ac6-8a4c-8edbe892d909",
				),
			),
			service.NewPlan(
				buildStandardPlan(
					"8ec86bea-42f6-4805-b3e9-506eaebbf9e0",
				),
			),
			service.NewPlan(
				buildPremiumPlan(
					"9c8c9dd1-fd0e-49a6-8178-7b3a21e5d4e0",
				),
			),
			service.NewPlan(
				buildGeneralPurposePlan(
					"5408c07e-8a08-4ff6-bd4a-967099bb3a1e",
				),
			),
			service.NewPlan(
				buildBusinessCriticalPlan(
					"82d21981-b8b8-4f06-81d7-cfb4011aedd7",
				),
			),
		),
		// database pair from existing service
		service.NewService(
			service.ServiceProperties{
				ID:              "e18a9861-5740-4e1a-9bd0-6f0fc3e4d12f",
				Name:            "azure-sql-12-0-dr-database-pair-from-existing",
				Description:     "Azure SQL 12.0-- disaster recovery database pair from existing",
				Bindable:        true,
				ParentServiceID: "00ce53a3-d6c3-4c24-8cb2-3f48d3b161d8",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure SQL 12.0-- disaster recovery Database Pair from existing",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/SQLDatabase.svg",
					LongDescription:  "Azure SQL 12.0-- disaster recovery database pair from existing",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Tags: []string{
					"Azure",
					"SQL",
					"Database",
					service.DRTag,
					service.MigrationTag,
					"Failover Group",
				},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.databasePairManagerForExistingPair,
			service.NewPlan(
				buildBasicPlan(
					"5ffdb255-8261-4841-a0f4-f1ec4ed9402c",
				),
			),
			service.NewPlan(
				buildStandardPlan(
					"af66e3e3-c500-4042-879e-5a6d47901d1c",
				),
			),
			service.NewPlan(
				buildPremiumPlan(
					"e3b13175-5686-443e-aee0-33f76d62ab55",
				),
			),
			service.NewPlan(
				buildGeneralPurposePlan(
					"94a2888f-dcec-4f6a-bea8-f7a7e4f6edc8",
				),
			),
			service.NewPlan(
				buildBusinessCriticalPlan(
					"0a7d3d85-147a-4bae-9c67-055e8404af1e",
				),
			),
		),
	}), nil
}
