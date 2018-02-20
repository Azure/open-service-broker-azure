package sqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		// all-in-one (server and db) service
		service.NewService(
			&service.ServiceProperties{
				ID:          "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
				Name:        "azure-sqldb",
				Description: "Azure SQL Database (Experimental)",
				Metadata: &service.ServiceMetadata{
					DisplayName: "Azure SQL Database",
					ImageUrl: "https://azure.microsoft.com/svghandler/sql-database/" +
						"?width=200",
					LongDescription: "The intelligent relational cloud database service" +
						" (Experimental)",
					DocumentationUrl: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportUrl:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "SQL", "Database"},
			},
			m.allInOneServiceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "3819fdfa-0aaa-11e6-86f4-000d3a002ed5",
				Name:        "basic",
				Description: "Basic Tier, 5 DTUs, 2GB, 7 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Basic",
					"requestedServiceObjectiveName": "Basic",
					"maxSizeBytes":                  "2147483648",
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
				ID:          "2497b7f3-341b-4ac6-82fb-d4a48c005e19",
				Name:        "standard-s0",
				Description: "Standard Tier, 10 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Standard",
					"requestedServiceObjectiveName": "S0",
					"maxSizeBytes":                  "268435456000",
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
				ID:          "17725188-76a2-4d6c-8e86-49f146766eeb",
				Name:        "standard-s1",
				Description: "StandardS1 Tier, 20 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Standard",
					"requestedServiceObjectiveName": "S1",
					"maxSizeBytes":                  "268435456000",
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
				ID:          "a5537f8e-d816-4b0e-9546-a13811944bdd",
				Name:        "standard-s2",
				Description: "StandardS2 Tier, 50 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Standard",
					"requestedServiceObjectiveName": "S2",
					"maxSizeBytes":                  "268435456000",
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
				ID:          "26cf84bf-f700-4e65-8048-cbfa9c319d5f",
				Name:        "standard-s3",
				Description: "StandardS3 Tier, 100 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Standard",
					"requestedServiceObjectiveName": "S3",
					"maxSizeBytes":                  "268435456000",
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
				ID:          "f9a3cc8e-a6e2-474d-b032-9837ea3dfcaa",
				Name:        "premium-p1",
				Description: "PremiumP1 Tier, 125 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P1",
					"maxSizeBytes":                  "536870912000",
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
				ID:          "2bbbcc59-a0e0-4153-841b-2833cb417d43",
				Name:        "premium-p2",
				Description: "PremiumP2 Tier, 250 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P2",
					"maxSizeBytes":                  "536870912000",
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
				ID:          "85d54d69-55ee-4fe8-a207-66bc96ecf9e7",
				Name:        "premium-p4",
				Description: "PremiumP4 Tier, 500 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P4",
					"maxSizeBytes":                  "536870912000",
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
				ID:          "af3dc76f-5b31-4cad-8adc-a9e756640a57",
				Name:        "premium-p6",
				Description: "PremiumP6 Tier, 1000 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P6",
					"maxSizeBytes":                  "536870912000",
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
				ID:          "408f5f35-5f5e-48f3-98cf-9e10c1abc4e5",
				Name:        "premium-p11",
				Description: "PremiumP11 Tier, 1750 DTUs, 1024GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P11",
					"maxSizeBytes":                  "1099511627776",
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "PremiumP11 Tier",
					Bullets: []string{
						"1024GB",
						"35 days point-in-time restore",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "b69af389-7af5-47bd-9ccf-c1ffdc2620d9",
				Name:        "data-warehouse-100",
				Description: "DataWarehouse100 Tier, 100 DWUs, 1024GB",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "DataWarehouse",
					"requestedServiceObjectiveName": "DW100",
					"maxSizeBytes":                  "1099511627776",
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "DataWarehouse100 Tier",
					Bullets: []string{
						"1024GB",
						"100 DWUs",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "470a869b-1b02-474b-b5e5-10ca0ea488df",
				Name:        "data-warehouse-1200",
				Description: "DataWarehouse1200 Tier, 1200 DWUs, 1024GB",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "DataWarehouse",
					"requestedServiceObjectiveName": "DW1200",
					"maxSizeBytes":                  "1099511627776",
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "DataWarehouse1200 Tier",
					Bullets: []string{
						"1024GB",
						"1200 DWUs",
					},
				},
			}),
		),
		// vm only service
		service.NewService(
			&service.ServiceProperties{
				ID:          "a7454e0e-be2c-46ac-b55f-8c4278117525",
				Name:        "azure-sqldb-dbms-only",
				Description: "Azure SQL Server DBMS (Experimental)",
				Metadata: &service.ServiceMetadata{
					DisplayName: "Azure SQL Server (DBMS Only)",
					ImageUrl: "https://azure.microsoft.com/svghandler/sql-database/" +
						"?width=200",
					LongDescription: "The intelligent relational cloud database service" +
						" (Experimental)",
					DocumentationUrl: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportUrl:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable:       false,
				Tags:           []string{"Azure", "SQL", "Server", "DBMS"},
				ChildServiceID: "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
			},
			m.dbmsOnlyManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "24f0f42e-1ab3-474e-a5ca-b943b2c48eee",
				Name:        "sqldb-dbms-only",
				Description: "Azure SQL Server DBMS Only",
				Free:        false,
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Azure SQL Server (DBMS Only)",
				},
			}),
		),
		// db only service
		service.NewService(
			&service.ServiceProperties{
				ID:              "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				Name:            "azure-sqldb-db-only",
				Description:     "Azure SQL Database Only (Experimental)",
				Bindable:        true,
				Tags:            []string{"Azure", "SQL", "Database"},
				ParentServiceID: "a7454e0e-be2c-46ac-b55f-8c4278117525",
				Metadata: &service.ServiceMetadata{
					DisplayName: "Azure SQL Server (Database Only)",
					ImageUrl: "https://azure.microsoft.com/svghandler/sql-database/" +
						"?width=200",
					LongDescription: "The intelligent relational cloud database service" +
						" (Experimental)",
					DocumentationUrl: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportUrl:       "https://azure.microsoft.com/en-us/support/",
				},
			},
			m.dbOnlyServiceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "8fa8d759-c142-45dd-ae38-b93482ddc04a",
				Name:        "basic",
				Description: "Basic Tier, 5 DTUs, 2GB, 7 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Basic",
					"requestedServiceObjectiveName": "Basic",
					"maxSizeBytes":                  "2147483648",
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
					"edition":                       "Standard",
					"requestedServiceObjectiveName": "S0",
					"maxSizeBytes":                  "268435456000",
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
					"edition":                       "Standard",
					"requestedServiceObjectiveName": "S1",
					"maxSizeBytes":                  "268435456000",
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
					"edition":                       "Standard",
					"requestedServiceObjectiveName": "S2",
					"maxSizeBytes":                  "268435456000",
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
					"edition":                       "Standard",
					"requestedServiceObjectiveName": "S3",
					"maxSizeBytes":                  "268435456000",
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
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P1",
					"maxSizeBytes":                  "536870912000",
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
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P2",
					"maxSizeBytes":                  "536870912000",
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
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P4",
					"maxSizeBytes":                  "536870912000",
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
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P6",
					"maxSizeBytes":                  "536870912000",
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
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P11",
					"maxSizeBytes":                  "1099511627776",
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "PremiumP11 Tier",
					Bullets: []string{
						"1024GB",
						"35 days point-in-time restore",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "7a466f47-f137-4b9c-a63d-c5cbe724b874",
				Name:        "data-warehouse-100",
				Description: "DataWarehouse100 Tier, 100 DWUs, 1024GB",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "DataWarehouse",
					"requestedServiceObjectiveName": "DW100",
					"maxSizeBytes":                  "1099511627776",
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "DataWarehouse100 Tier",
					Bullets: []string{
						"1024GB",
						"100 DWUs",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "2717d839-be32-4225-8685-47adf0e6ff15",
				Name:        "data-warehouse-1200",
				Description: "DataWarehouse1200 Tier, 1200 DWUs, 1024GB",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "DataWarehouse",
					"requestedServiceObjectiveName": "DW1200",
					"maxSizeBytes":                  "1099511627776",
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "DataWarehouse1200 Tier",
					Bullets: []string{
						"1024GB",
						"1200 DWUs",
					},
				},
			}),
		),
	}), nil
}
