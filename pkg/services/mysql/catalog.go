package mysql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func createBasicPlan(
	planID string,
) *service.PlanProperties {
	provisionSchema := planSchema{
		allowedHardware:         []string{"", "gen4", "gen5"},
		defaultHardware:         "gen5",
		validCores:              []int{1, 2},
		defaultCores:            1,
		maxStorage:              1024,
		minStorage:              5,
		defaultStorage:          10,
		allowedBackupRedundancy: []string{"local"},
		minBackupRetention:      7,
		maxBackupRetention:      35,
		defaultBackupRetention:  7,
		tier: "B",
	}

	return &service.PlanProperties{
		ID:          planID,
		Name:        "basic",
		Description: "Basic Tier",
		Free:        false,
		Extended: map[string]interface{}{
			"provisionSchema": provisionSchema,
			"tier":            "Basic",
		},
		Metadata: &service.ServicePlanMetadata{
			DisplayName: "Basic Tier",
			Bullets:     []string{"Up to 2 vCores", "Variable I/O performance"},
		},
		ProvisionParamsSchema: generateDBMSPlanSchema(provisionSchema),
	}
}

func createGPPlan(
	planID string,
) *service.PlanProperties {

	provisionSchema := planSchema{
		allowedHardware:         []string{"", "gen4", "gen5"},
		defaultHardware:         "gen5",
		validCores:              []int{2, 4, 8, 16, 32},
		defaultCores:            2,
		maxStorage:              2048,
		minStorage:              5,
		defaultStorage:          10,
		allowedBackupRedundancy: []string{"local", "geo"},
		minBackupRetention:      7,
		maxBackupRetention:      35,
		defaultBackupRetention:  7,
		tier: "GP",
	}
	extendedPlanData := map[string]interface{}{
		"provisionSchema": provisionSchema,
		"tier":            "GeneralPurpose",
	}

	return &service.PlanProperties{
		ID:          planID,
		Name:        "general-purpose",
		Description: "General Purpose",
		Free:        false,
		Extended:    extendedPlanData,
		Metadata: &service.ServicePlanMetadata{
			DisplayName: "General Purpose Tier",
			Bullets: []string{
				"Up to 32 vCores",
				"Predictable I/O Performance",
				"Local or Geo-Redundant Backups",
			},
		},
		ProvisionParamsSchema: generateDBMSPlanSchema(provisionSchema),
	}
}

func createMemoryOptimizedPlan(
	planID string,
) *service.PlanProperties {

	provisionSchema := planSchema{
		allowedHardware:         []string{"", "gen5"},
		defaultHardware:         "gen5",
		validCores:              []int{2, 4, 8, 16},
		defaultCores:            2,
		maxStorage:              2048,
		minStorage:              5,
		defaultStorage:          10,
		allowedBackupRedundancy: []string{"local", "geo"},
		minBackupRetention:      7,
		maxBackupRetention:      35,
		defaultBackupRetention:  7,
		tier: "MO",
	}
	extendedPlanData := map[string]interface{}{
		"provisionSchema": provisionSchema,
		"tier":            "MemoryOptimized",
	}

	return &service.PlanProperties{
		ID:          planID,
		Name:        "memory-optimized",
		Description: "Memory Optimized",
		Free:        false,
		Extended:    extendedPlanData,
		Metadata: &service.ServicePlanMetadata{
			DisplayName: "Memory Optimized Tier",
			Bullets: []string{
				"Up to 16 memory optimized vCores",
				"Predictable I/O Performance",
				"Local or Geo-Redundant Backups",
			},
		},
		ProvisionParamsSchema: generateDBMSPlanSchema(provisionSchema),
	}
}

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
				Extended: map[string]interface{}{
					"version": "5.7",
				},
			},
			m.allInOneServiceManager,
			service.NewPlan(createBasicPlan("1b093840-8e02-4e28-9aba-fa716757ec38")),
			service.NewPlan(createGPPlan("eae202c3-521c-46d1-a047-872dacf781fd")),
			service.NewPlan(createMemoryOptimizedPlan("129f06f6-cbf2-416e-a235-0fa6e081a07a")),
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
				Extended: map[string]interface{}{
					"version": "5.7",
				},
			},
			m.dbmsManager,
			service.NewPlan(createBasicPlan("20938530-cb42-48b2-93dc-ea0d3003a89f")),
			service.NewPlan(createGPPlan("3a00b95f-6acf-4bf9-8b01-52fe03a2d607")),
			service.NewPlan(createMemoryOptimizedPlan("b242a78f-9946-406a-af67-813c56341960")),
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
