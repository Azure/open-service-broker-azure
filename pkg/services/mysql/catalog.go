package mysql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

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
				ProvisionParamsSchema: m.allInOneServiceManager.getProvisionParametersSchema(),
				Extended: map[string]interface{}{
					"version": "5.7",
				},
			},
			m.allInOneServiceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "427559f1-bf2a-45d3-8844-32374a3e58aa",
				Name:        "basic50",
				Description: "Basic Tier, 50 DTUs.",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLB50",
					"skuTier":        "Basic",
					"skuCapacityDTU": 50,
					"skuSizeMB":      51200,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets:     []string{"50 DTUs"},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "1a538e06-9bcc-4077-8480-966cbf85bf36",
				Name:        "basic100",
				Description: "Basic Tier, 100 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLB100",
					"skuTier":        "Basic",
					"skuCapacityDTU": 100,
					"skuSizeMB":      51200,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets:     []string{"100 DTUs"},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "edc2badc-d93b-4d9c-9d8e-da2f1c8c3e1c",
				Name:        "standard100",
				Description: "Standard Tier, 100 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLS100",
					"skuTier":        "Standard",
					"skuCapacityDTU": 100,
					"skuSizeMB":      128000,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets: []string{
						"100 DTUs",
						"Additional Storage",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "9995c891-48ba-46cc-8dae-83595c1f443f",
				Name:        "standard200",
				Description: "Standard Tier, 200 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLS200",
					"skuTier":        "Standard",
					"skuCapacityDTU": 200,
					"skuSizeMB":      128000,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets: []string{
						"200 DTUs",
						"Additional Storage",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "ae3cd3dd-9818-48c0-9cd0-62c3b130944e",
				Name:        "standard400",
				Description: "Standard Tier, 400 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLS400",
					"skuTier":        "Standard",
					"skuCapacityDTU": 400,
					"skuSizeMB":      128000,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets: []string{
						"400 DTUs",
						"Additional Storage",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "08e4b43a-36bc-447e-a81f-8202b13e339c",
				Name:        "standard800",
				Description: "Standard Tier, 800 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLS800",
					"skuTier":        "Standard",
					"skuCapacityDTU": 800,
					"skuSizeMB":      128000,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets: []string{
						"800 DTUs",
						"Additional Storage",
					},
				},
			}),
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
				ProvisionParamsSchema: m.dbmsManager.getProvisionParametersSchema(),
				Extended: map[string]interface{}{
					"version": "5.7",
				},
			},
			m.dbmsManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "3f65ebf9-ac1d-4e77-b9bf-918889a4482b",
				Name:        "basic50",
				Description: "Basic Tier, 50 DTUs.",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLB50",
					"skuTier":        "Basic",
					"skuCapacityDTU": 50,
					"skuSizeMB":      51200,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets:     []string{"50 DTUs"},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "9f71584c-8e97-46a7-b170-20c4273a64f9",
				Name:        "basic100",
				Description: "Basic Tier, 100 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLB100",
					"skuTier":        "Basic",
					"skuCapacityDTU": 100,
					"skuSizeMB":      51200,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets:     []string{"100 DTUs"},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "dac995d8-2618-4aa5-9f2b-0376914ed2f7",
				Name:        "standard100",
				Description: "Standard Tier, 100 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLS100",
					"skuTier":        "Standard",
					"skuCapacityDTU": 100,
					"skuSizeMB":      128000,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets: []string{
						"100 DTUs",
						"Additional Storage",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "1c7cf479-7dba-4ed4-a855-9ab032c40466",
				Name:        "standard200",
				Description: "Standard Tier, 200 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLS200",
					"skuTier":        "Standard",
					"skuCapacityDTU": 200,
					"skuSizeMB":      128000,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets: []string{
						"200 DTUs",
						"Additional Storage",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "d8565a53-1db0-4842-9e64-5a5df560b668",
				Name:        "standard400",
				Description: "Standard Tier, 400 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLS400",
					"skuTier":        "Standard",
					"skuCapacityDTU": 400,
					"skuSizeMB":      128000,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets: []string{
						"400 DTUs",
						"Additional Storage",
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "6765fa7b-6b0a-4560-960f-7425dac56d47",
				Name:        "standard800",
				Description: "Standard Tier, 800 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLS800",
					"skuTier":        "Standard",
					"skuCapacityDTU": 800,
					"skuSizeMB":      128000,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets: []string{
						"800 DTUs",
						"Additional Storage",
					},
				},
			}),
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
