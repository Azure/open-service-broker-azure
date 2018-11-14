package mssqldr

import "github.com/Azure/open-service-broker-azure/pkg/service"

func buildBasicPlan(
	id string, // nolint: unparam
	// TODO remove the comment
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
	id string, // nolint: unparam
	// TODO remove the comment
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
	id string, // nolint: unparam
	// TODO remove the comment
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
	id string, // nolint: unparam
	// TODO remove the comment
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
	id string, // nolint: unparam
	// TODO remove the comment
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
					DisplayName:      "Azure SQL 12.0-- disaster recovery DBMS Pair registered",
					ImageURL:         "https://azure.microsoft.com/svghandler/sql-database/?width=200",
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
					DisplayName:      "Azure SQL 12.0-- disaster recovery Database Pair",
					ImageURL:         "https://azure.microsoft.com/svghandler/sql-database/?width=200",
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
					DisplayName:      "Azure SQL 12.0-- disaster recovery Database Pair registered",
					ImageURL:         "https://azure.microsoft.com/svghandler/sql-database/?width=200",
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
	}), nil
}
