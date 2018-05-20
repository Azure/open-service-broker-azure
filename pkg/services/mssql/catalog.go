package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

func buildGeneralPurposePlan(
	id string,
	includesDBMS bool,
) *service.PlanProperties {
	gpDetails := vCorePlanDetails{
		tier:          "GeneralPurpose",
		tierShortName: "GP",
		includesDBMS:  includesDBMS,
	}
	return &service.PlanProperties{
		ID:          id,
		Name:        "general-purpose",
		Description: "Scalable compute and storage options for budget-oriented applications",
		Free:        false,
		Extended: map[string]interface{}{
			"tierDetails": gpDetails,
		},
		Metadata: &service.ServicePlanMetadata{
			DisplayName: "General Purpose",
			Bullets: []string{
				"Up to 80 vCores",
				"Up to 440 GB memory",
				"$187.62 / vCore",
				"7 days point-in-time restore",
				"Currently In Preview",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: gpDetails.getProvisionSchema(),
			},
		},
	}
}

func buildBusinessCriticalPlan(
	id string,
	includesDBMS bool,
) *service.PlanProperties {
	bcDetails := vCorePlanDetails{
		tier:          "BusinessCritical",
		tierShortName: "BC",
		includesDBMS:  includesDBMS,
	}
	return &service.PlanProperties{
		ID:          id,
		Name:        "business-critical",
		Description: "For applications with high transaction rate and highly resilient to failure",
		Free:        false,
		Extended: map[string]interface{}{
			"tierDetails": bcDetails,
		},
		Metadata: &service.ServicePlanMetadata{
			DisplayName: "Basic Tier",
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
				ProvisioningParametersSchema: bcDetails.getProvisionSchema(),
			},
		},
	}
}

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {

	return service.NewCatalog([]service.Service{
		// all-in-one (dbms and database) service
		service.NewService(
			&service.ServiceProperties{
				ID:          "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
				Name:        "azure-sql-12-0",
				Description: "Azure SQL Database 12.0-- DBMS and single database (preview)",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "Azure SQL Database 12.0 (preview)",
					ImageURL:         "https://azure.microsoft.com/svghandler/sql-database/?width=200",
					LongDescription:  "Azure SQL Database 12.0-- DBMS and single database (preview)",
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
			service.NewPlan(&service.PlanProperties{
				ID:          "3819fdfa-0aaa-11e6-86f4-000d3a002ed5",
				Name:        "basic",
				Description: "Basic Tier, 5 DTUs, 2GB, 7 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Basic",
						sku:        "Basic",
						maxStorage: 2147483648,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets: []string{
						"5 DTUs",
						"2GB",
						"7 days point-in-time restore",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.allInOneServiceManager.getProvisionParametersSchema(),
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "2497b7f3-341b-4ac6-82fb-d4a48c005e19",
				Name:        "standard-s0",
				Description: "Standard Tier, 10 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Standard",
						sku:        "S0",
						maxStorage: 268435456000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets: []string{
						"10 DTUs",
						"250GB",
						"35 days point-in-time restore",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.allInOneServiceManager.getProvisionParametersSchema(),
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "17725188-76a2-4d6c-8e86-49f146766eeb",
				Name:        "standard-s1",
				Description: "StandardS1 Tier, 20 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Standard",
						sku:        "S1",
						maxStorage: 268435456000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "StandardS1 Tier",
					Bullets: []string{
						"20 DTUs",
						"250GB",
						"35 days point-in-time restore",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.allInOneServiceManager.getProvisionParametersSchema(),
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "a5537f8e-d816-4b0e-9546-a13811944bdd",
				Name:        "standard-s2",
				Description: "StandardS2 Tier, 50 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Standard",
						sku:        "S2",
						maxStorage: 268435456000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "StandardS2 Tier",
					Bullets: []string{
						"50 DTUs",
						"250GB",
						"35 days point-in-time restore",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.allInOneServiceManager.getProvisionParametersSchema(),
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "26cf84bf-f700-4e65-8048-cbfa9c319d5f",
				Name:        "standard-s3",
				Description: "StandardS3 Tier, 100 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Standard",
						sku:        "S3",
						maxStorage: 268435456000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "StandardS3 Tier",
					Bullets: []string{
						"100 DTUs",
						"250GB",
						"35 days point-in-time restore",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.allInOneServiceManager.getProvisionParametersSchema(),
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "f9a3cc8e-a6e2-474d-b032-9837ea3dfcaa",
				Name:        "premium-p1",
				Description: "PremiumP1 Tier, 125 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Premium",
						sku:        "P1",
						maxStorage: 536870912000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "PremiumP1 Tier",
					Bullets: []string{
						"125 DTUs",
						"500GB",
						"35 days point-in-time restore",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.allInOneServiceManager.getProvisionParametersSchema(),
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "2bbbcc59-a0e0-4153-841b-2833cb417d43",
				Name:        "premium-p2",
				Description: "PremiumP2 Tier, 250 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Premium",
						sku:        "P2",
						maxStorage: 536870912000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "PremiumP2 Tier",
					Bullets: []string{
						"250 DTUs",
						"500GB",
						"35 days point-in-time restore",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.allInOneServiceManager.getProvisionParametersSchema(),
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "85d54d69-55ee-4fe8-a207-66bc96ecf9e7",
				Name:        "premium-p4",
				Description: "PremiumP4 Tier, 500 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Premium",
						sku:        "P4",
						maxStorage: 536870912000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "PremiumP4 Tier",
					Bullets: []string{
						"500 DTUs",
						"500GB",
						"35 days point-in-time restore",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.allInOneServiceManager.getProvisionParametersSchema(),
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "af3dc76f-5b31-4cad-8adc-a9e756640a57",
				Name:        "premium-p6",
				Description: "PremiumP6 Tier, 1000 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Premium",
						sku:        "P6",
						maxStorage: 536870912000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "PremiumP6 Tier",
					Bullets: []string{
						"1000 DTUs",
						"500GB",
						"35 days point-in-time restore",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.allInOneServiceManager.getProvisionParametersSchema(),
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "408f5f35-5f5e-48f3-98cf-9e10c1abc4e5",
				Name:        "premium-p11",
				Description: "PremiumP11 Tier, 1750 DTUs, 1024GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Premium",
						sku:        "P11",
						maxStorage: 1099511627776,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "PremiumP11 Tier",
					Bullets: []string{
						"1024GB",
						"35 days point-in-time restore",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.allInOneServiceManager.getProvisionParametersSchema(),
					},
				},
			}),
			service.NewPlan(
				buildGeneralPurposePlan(
					"c77e86af-f050-4457-a2ff-2b48451888f3",
					true,
				),
			),
			service.NewPlan(
				buildBusinessCriticalPlan(
					"ebc3ae35-91bc-480c-807b-e798c1ca8c4e",
					true,
				),
			),
		),
		// dbms only service
		service.NewService(
			&service.ServiceProperties{
				ID:             "a7454e0e-be2c-46ac-b55f-8c4278117525",
				Name:           "azure-sql-12-0-dbms",
				Description:    "Azure SQL 12.0-- DBMS only (preview)",
				ChildServiceID: "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "Azure SQL 12.0-- DBMS Only (preview)",
					ImageURL:         "https://azure.microsoft.com/svghandler/sql-database/?width=200",
					LongDescription:  "Azure SQL 12.0-- DBMS only (preview)",
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
			service.NewPlan(&service.PlanProperties{
				ID:          "24f0f42e-1ab3-474e-a5ca-b943b2c48eee",
				Name:        "dbms",
				Description: "Azure SQL Server-- DBMS only",
				Free:        false,
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Azure SQL Server-- DBMS Only",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.dbmsManager.getProvisionParametersSchema(),
					},
				},
			}),
		),
		// database only service
		service.NewService(
			&service.ServiceProperties{
				ID:              "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				Name:            "azure-sql-12-0-database",
				Description:     "Azure SQL 12.0-- database only (preview)",
				Bindable:        true,
				ParentServiceID: "a7454e0e-be2c-46ac-b55f-8c4278117525",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "Azure SQL 12.0-- Database Only (preview)",
					ImageURL:         "https://azure.microsoft.com/svghandler/sql-database/?width=200",
					LongDescription:  "Azure SQL 12.0-- database only (preview)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Tags: []string{"Azure", "SQL", "Database"},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.databaseManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "8fa8d759-c142-45dd-ae38-b93482ddc04a",
				Name:        "basic",
				Description: "Basic Tier, 5 DTUs, 2GB, 7 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Basic",
						sku:        "Basic",
						maxStorage: 2147483648,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets: []string{
						"5 DTUs",
						"2GB",
						"7 days point-in-time restore",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "9d36b6b3-b5f3-4907-a713-5cc13b785409",
				Name:        "standard-s0",
				Description: "Standard Tier, 10 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Standard",
						sku:        "S0",
						maxStorage: 268435456000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets: []string{
						"10 DTUs",
						"250GB",
						"35 days point-in-time restore",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "01c397f8-c999-4e86-bcc2-654cd8cae5fd",
				Name:        "standard-s1",
				Description: "StandardS1 Tier, 20 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Standard",
						sku:        "S1",
						maxStorage: 268435456000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "StandardS1 Tier",
					Bullets: []string{
						"20 DTUs",
						"250GB",
						"35 days point-in-time restore",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "9cd114a0-8356-4247-9b71-2e685e5a29f3",
				Name:        "standard-s2",
				Description: "StandardS2 Tier, 50 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Standard",
						sku:        "S2",
						maxStorage: 268435456000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "StandardS2 Tier",
					Bullets: []string{
						"50 DTUs",
						"250GB",
						"35 days point-in-time restore",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "624828a9-c73c-4d35-bc9d-ea41cfc75853",
				Name:        "standard-s3",
				Description: "StandardS3 Tier, 100 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Standard",
						sku:        "S3",
						maxStorage: 268435456000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "StandardS3 Tier",
					Bullets: []string{
						"100 DTUs",
						"250GB",
						"35 days point-in-time restore",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "220e922a-a5b2-43e4-9388-fe45a32bbf31",
				Name:        "premium-p1",
				Description: "PremiumP1 Tier, 125 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Premium",
						sku:        "P1",
						maxStorage: 536870912000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "PremiumP1 Tier",
					Bullets: []string{
						"125 DTUs",
						"500GB",
						"35 days point-in-time restore",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "e7eb13df-1fda-4492-b218-00dd0db1c85d",
				Name:        "premium-p2",
				Description: "PremiumP2 Tier, 250 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Premium",
						sku:        "P2",
						maxStorage: 536870912000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "PremiumP2 Tier",
					Bullets: []string{
						"250 DTUs",
						"500GB",
						"35 days point-in-time restore",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "feb25d68-2b52-41b5-a249-28a747bc2c2e",
				Name:        "premium-p4",
				Description: "PremiumP4 Tier, 500 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Premium",
						sku:        "P4",
						maxStorage: 536870912000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "PremiumP4 Tier",
					Bullets: []string{
						"500 DTUs",
						"500GB",
						"35 days point-in-time restore",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "19487202-dc8a-4930-bbad-7bbf1486dbca",
				Name:        "premium-p6",
				Description: "PremiumP6 Tier, 1000 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Premium",
						sku:        "P6",
						maxStorage: 536870912000,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "PremiumP6 Tier",
					Bullets: []string{
						"1000 DTUs",
						"500GB",
						"35 days point-in-time restore",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "a561c45a-33c8-412e-9315-411c1d7035da",
				Name:        "premium-p11",
				Description: "PremiumP11 Tier, 1750 DTUs, 1024GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"tierDetails": legacyPlanDetails{
						tier:       "Premium",
						sku:        "P11",
						maxStorage: 1099511627776,
					},
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "PremiumP11 Tier",
					Bullets: []string{
						"1024GB",
						"35 days point-in-time restore",
					},
				},
			}),
			service.NewPlan(
				buildGeneralPurposePlan(
					"da591616-77a1-4df8-a493-6c119649bc6b",
					false,
				),
			),
			service.NewPlan(
				buildBusinessCriticalPlan(
					"b05c25d2-1d63-4d09-a50a-e68c2710a069",
					false,
				),
			),
		),
	}), nil
}
