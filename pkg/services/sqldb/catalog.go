package sqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
				Name:        "azure-sqldb",
				Description: "Azure SQL Database (Experimental)",
				Bindable:    true,
				Tags:        []string{"Azure", "SQL", "Database", "VM"},
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
			}),
		),
		service.NewService(
			&service.ServiceProperties{
				ID:          "2787cd60-8184-4b80-aa45-f507fa5a6ff4",
				Name:        "azure-sqldb-server-only",
				Description: "Azure SQL Server VM (Experimental)",
				Bindable:    false,
				Tags:        []string{"Azure", "SQL", "Server", "VM"},
			},
			m.vmOnlyServiceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "0f4baa94-92cb-4222-9d7e-600c394ec50d",
				Name:        "sql-server",
				Description: "Azure SQL Server - VM Only, No Database",
				Free:        false,
			}),
		),
		//db only
		service.NewService(
			&service.ServiceProperties{
				ID:          "d79e4073-9e30-4cf8-8485-f4284d3cf904",
				Name:        "azure-sqldb-db-only",
				Description: "Azure SQL Database Only (Experimental)",
				Bindable:    true,
				Tags:        []string{"Azure", "SQL", "Database"},
			},
			m.dbOnlyServiceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "c653f05c-b16b-4e3f-9b19-32356c6a4df1",
				Name:        "basic",
				Description: "Basic Tier, 5 DTUs, 2GB, 7 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Basic",
					"requestedServiceObjectiveName": "Basic",
					"maxSizeBytes":                  "2147483648",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "9636b94f-150c-44c3-9965-a37f4ac9a8ef",
				Name:        "standard-s0",
				Description: "Standard Tier, 10 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Standard",
					"requestedServiceObjectiveName": "S0",
					"maxSizeBytes":                  "268435456000",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "3c3119ea-c5d9-4185-84f6-0e0ca68927fd",
				Name:        "standard-s1",
				Description: "StandardS1 Tier, 20 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Standard",
					"requestedServiceObjectiveName": "S1",
					"maxSizeBytes":                  "268435456000",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "dee68ba9-d958-49a1-a16a-bdd8d95f1d4c",
				Name:        "standard-s2",
				Description: "StandardS2 Tier, 50 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Standard",
					"requestedServiceObjectiveName": "S2",
					"maxSizeBytes":                  "268435456000",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "091fba49-986c-4359-93b6-340239eb2de6",
				Name:        "standard-s3",
				Description: "StandardS3 Tier, 100 DTUs, 250GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Standard",
					"requestedServiceObjectiveName": "S3",
					"maxSizeBytes":                  "268435456000",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "3d6b5dff-f892-4e35-b5e1-ed91a7fd1543",
				Name:        "premium-p1",
				Description: "PremiumP1 Tier, 125 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P1",
					"maxSizeBytes":                  "536870912000",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "6add0078-b045-4d18-8e07-7f8bac4699bc",
				Name:        "premium-p2",
				Description: "PremiumP2 Tier, 250 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P2",
					"maxSizeBytes":                  "536870912000",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "9d6ac038-67ea-404d-b5da-c60631c12a1c",
				Name:        "premium-p4",
				Description: "PremiumP4 Tier, 500 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P4",
					"maxSizeBytes":                  "536870912000",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "55e788a6-dfff-443f-9cf4-35bc88051cdd",
				Name:        "premium-p6",
				Description: "PremiumP6 Tier, 1000 DTUs, 500GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P6",
					"maxSizeBytes":                  "536870912000",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "5c536485-933d-47da-957b-3422237eb46c",
				Name:        "premium-p11",
				Description: "PremiumP11 Tier, 1750 DTUs, 1024GB, 35 days point-in-time restore",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "Premium",
					"requestedServiceObjectiveName": "P11",
					"maxSizeBytes":                  "1099511627776",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "27d3fd70-3b76-450c-afab-607a41e275d2",
				Name:        "data-warehouse-100",
				Description: "DataWarehouse100 Tier, 100 DWUs, 1024GB",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "DataWarehouse",
					"requestedServiceObjectiveName": "DW100",
					"maxSizeBytes":                  "1099511627776",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "9a0c6c42-e4a5-4d82-beee-be3c32eb203b",
				Name:        "data-warehouse-1200",
				Description: "DataWarehouse1200 Tier, 1200 DWUs, 1024GB",
				Free:        false,
				Extended: map[string]interface{}{
					"edition":                       "DataWarehouse",
					"requestedServiceObjectiveName": "DW1200",
					"maxSizeBytes":                  "1099511627776",
				},
			}),
		),
	}), nil
}
