package mssqldr

import "github.com/Azure/open-service-broker-azure/pkg/service"

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
	}), nil
}
