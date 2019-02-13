package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

func buildBasicPlan(
	id string,
	includesDBMS bool,
	forExistingInstance bool,
) service.PlanProperties {

	planDetails := dtuPlanDetails{
		storageInGB: 2,
		defaultDTUs: 5,
		tierName:    "Basic",
		skuMap: map[int64]string{
			5: "Basic",
		},
		includeDBMS: includesDBMS,
	}

	planProperties := service.PlanProperties{
		ID:          id,
		Name:        "basic",
		Description: "Basic Tier, 5 DTUs, 2GB Storage, 7 days point-in-time restore",
		Free:        false,
		Stability:   service.StabilityStable,
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
	if forExistingInstance {
		planProperties.Schemas.ServiceInstances.ProvisioningParametersSchema =
			planDetails.getProvisionSchemaForExistingIntance()
	}

	return planProperties
}

func buildStandardPlan(
	id string,
	includesDBMS bool,
	forExistingInstance bool,
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
		includeDBMS: includesDBMS,
	}

	planProperties := service.PlanProperties{
		ID:   id,
		Name: "standard",
		Description: "Standard Tier, Up to 3000 DTUs, 250GB Storage, " +
			"35 days point-in-time restore",
		Free:      false,
		Stability: service.StabilityStable,
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
	if forExistingInstance {
		planProperties.Schemas.ServiceInstances.ProvisioningParametersSchema =
			planDetails.getProvisionSchemaForExistingIntance()
	}

	return planProperties
}

func buildPremiumPlan(
	id string,
	includesDBMS bool,
	forExistingInstance bool,
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
		includeDBMS: includesDBMS,
	}

	planProperties := service.PlanProperties{
		ID:   id,
		Name: "premium",
		Description: "Premium Tier, Up to 4000 DTUs, 500GB Storage, " +
			"35 days point-in-time restore",
		Free:      false,
		Stability: service.StabilityStable,
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
	if forExistingInstance {
		planProperties.Schemas.ServiceInstances.ProvisioningParametersSchema =
			planDetails.getProvisionSchemaForExistingIntance()
	}

	return planProperties
}

func buildGeneralPurposePlan(
	id string,
	includesDBMS bool,
	forExistingInstance bool,
) service.PlanProperties {
	planDetails := vCorePlanDetails{
		tierName:      "GeneralPurpose",
		tierShortName: "GP",
		includeDBMS:   includesDBMS,
	}
	planProperties := service.PlanProperties{
		ID:          id,
		Name:        "general-purpose",
		Description: "Up to 80 vCores, 440 GB memory and 1 TB of storage (preview)",
		Free:        false,
		Stability:   service.StabilityPreview,
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
	if forExistingInstance {
		planProperties.Schemas.ServiceInstances.ProvisioningParametersSchema =
			planDetails.getProvisionSchemaForExistingIntance()
	}

	return planProperties
}

func buildBusinessCriticalPlan(
	id string,
	includesDBMS bool,
	forExistingInstance bool,
) service.PlanProperties {
	planDetails := vCorePlanDetails{
		tierName:      "BusinessCritical",
		tierShortName: "BC",
		includeDBMS:   includesDBMS,
	}
	planProperties := service.PlanProperties{
		ID:   id,
		Name: "business-critical",
		Description: "Up to 80 vCores, 440 GB memory and 1 TB of storage. " +
			"Local SSD, highest resilience to failures. (preview)",
		Free:      false,
		Stability: service.StabilityPreview,
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
	if forExistingInstance {
		planProperties.Schemas.ServiceInstances.ProvisioningParametersSchema =
			planDetails.getProvisionSchemaForExistingIntance()
	}

	return planProperties
}

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {

	return service.NewCatalog([]service.Service{
		// all-in-one (dbms and database) service
		service.NewService(
			service.ServiceProperties{
				ID:          "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
				Name:        "azure-sql-12-0",
				Description: "Azure SQL Database 12.0-- DBMS and single database",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure SQL Database 12.0",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/SQLDatabase.svg?sanitize=true",
					LongDescription:  "Azure SQL Database 12.0-- DBMS and single database",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "SQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.allInOneServiceManager,
			service.NewPlan(
				buildBasicPlan(
					"3819fdfa-0aaa-11e6-86f4-000d3a002ed5",
					true,
					false,
				),
			),
			service.NewPlan(
				buildStandardPlan(
					"2497b7f3-341b-4ac6-82fb-d4a48c005e19",
					true,
					false,
				),
			),
			service.NewPlan(
				buildPremiumPlan(
					"f9a3cc8e-a6e2-474d-b032-9837ea3dfcaa",
					true,
					false,
				),
			),
			service.NewPlan(
				buildGeneralPurposePlan(
					"c77e86af-f050-4457-a2ff-2b48451888f3",
					true,
					false,
				),
			),
			service.NewPlan(
				buildBusinessCriticalPlan(
					"ebc3ae35-91bc-480c-807b-e798c1ca8c4e",
					true,
					false,
				),
			),
		),
		// dbms only service
		service.NewService(
			service.ServiceProperties{
				ID:             "a7454e0e-be2c-46ac-b55f-8c4278117525",
				Name:           "azure-sql-12-0-dbms",
				Description:    "Azure SQL 12.0-- DBMS only",
				ChildServiceID: "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure SQL 12.0-- DBMS Only",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/SQLDatabase.svg?sanitize=true",
					LongDescription:  "Azure SQL 12.0-- DBMS only",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: false,
				Tags:     []string{"Azure", "SQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.dbmsManager,
			service.NewPlan(service.PlanProperties{
				ID:          "24f0f42e-1ab3-474e-a5ca-b943b2c48eee",
				Name:        "dbms",
				Description: "Azure SQL Server-- DBMS only",
				Free:        false,
				Stability:   service.StabilityPreview,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Azure SQL Server-- DBMS Only",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.dbmsManager.getProvisionParametersSchema(),
						UpdatingParametersSchema:     m.dbmsManager.getUpdatingParametersSchema(),
					},
				},
			}),
		),
		// database only service
		service.NewService(
			service.ServiceProperties{
				ID:              "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				Name:            "azure-sql-12-0-database",
				Description:     "Azure SQL 12.0-- database only",
				Bindable:        true,
				ParentServiceID: "a7454e0e-be2c-46ac-b55f-8c4278117525", // more parents in fact
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure SQL 12.0-- Database Only",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/SQLDatabase.svg?sanitize=true",
					LongDescription:  "Azure SQL 12.0-- database only",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Tags: []string{"Azure", "SQL", "Database"},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.databaseManager,
			service.NewPlan(
				buildBasicPlan(
					"8fa8d759-c142-45dd-ae38-b93482ddc04a",
					false,
					false,
				),
			),
			service.NewPlan(
				buildStandardPlan(
					"9d36b6b3-b5f3-4907-a713-5cc13b785409",
					false,
					false,
				),
			),
			service.NewPlan(
				buildPremiumPlan(
					"220e922a-a5b2-43e4-9388-fe45a32bbf31",
					false,
					false,
				),
			),
			service.NewPlan(
				buildGeneralPurposePlan(
					"da591616-77a1-4df8-a493-6c119649bc6b",
					false,
					false,
				),
			),
			service.NewPlan(
				buildBusinessCriticalPlan(
					"b05c25d2-1d63-4d09-a50a-e68c2710a069",
					false,
					false,
				),
			),
		),
		// dbms only registered service
		service.NewService(
			service.ServiceProperties{
				ID:             "c9bd94e7-5b7d-4b20-96e6-c5678f99d997",
				Name:           "azure-sql-12-0-dbms-registered",
				Description:    "Azure SQL 12.0-- DBMS only registered",
				ChildServiceID: "2bbc160c-e279-4757-a6b6-4c0a4822d0aa", // database-from-existing is also a valid child
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure SQL 12.0-- DBMS Only registered",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/SQLDatabase.svg?sanitize=true",
					LongDescription:  "Azure SQL 12.0-- DBMS only registered",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: false,
				Tags:     []string{"Azure", "SQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.dbmsRegisteredManager,
			service.NewPlan(service.PlanProperties{
				ID:          "4e95e962-0142-4117-b212-bcc7aec7f6c2",
				Name:        "dbms",
				Description: "Azure SQL Server-- DBMS only",
				Free:        false,
				Stability:   service.StabilityPreview,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Azure SQL Server-- DBMS Only",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.dbmsRegisteredManager.getProvisionParametersSchema(),
						UpdatingParametersSchema:     m.dbmsRegisteredManager.getUpdatingParametersSchema(),
					},
				},
			}),
		),
		// database only from existing service
		service.NewService(
			service.ServiceProperties{
				ID:              "b0b2a2f7-9b5e-4692-8b94-24fe2f6a9a8e",
				Name:            "azure-sql-12-0-database-from-existing",
				Description:     "Azure SQL 12.0-- database only from existing",
				Bindable:        true,
				ParentServiceID: "a7454e0e-be2c-46ac-b55f-8c4278117525", // dbms-registered is also a valid parent
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure SQL 12.0-- Database Only from existing",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/SQLDatabase.svg?sanitize=true",
					LongDescription:  "Azure SQL 12.0-- database only from existing",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Tags: []string{"Azure", "SQL", "Database", service.MigrationTag},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.databaseManagerForExistingInstance,
			service.NewPlan(
				buildBasicPlan(
					"e5804586-625a-4f67-996f-ca19a14711cc",
					false,
					true,
				),
			),
			service.NewPlan(
				buildStandardPlan(
					"ee01d17b-37ee-46b7-bfc7-39faf3230d02",
					false,
					true,
				),
			),
			service.NewPlan(
				buildPremiumPlan(
					"19f31593-6727-451e-95cf-3e64a90bd968",
					false,
					true,
				),
			),
			service.NewPlan(
				buildGeneralPurposePlan(
					"e8a788a0-2968-43ec-b8bb-5ecf2ce90ade",
					false,
					true,
				),
			),
			service.NewPlan(
				buildBusinessCriticalPlan(
					"9a26cd40-af08-4e05-bb8e-a521c3d3b60e",
					false,
					true,
				),
			),
		),
	}), nil
}
