package mssqldb

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
				Tags:        []string{"Azure", "SQL", "Database"},
			},
			m.serviceManager,
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
	}), nil
}
