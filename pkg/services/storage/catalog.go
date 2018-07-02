package storage

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

const kindKey = "kind"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			service.ServiceProperties{
				ID:          "2e2fc314-37b6-4587-8127-8f9ee8b33fea",
				Name:        "azure-storage",
				Description: "Azure Storage (Experimental)",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure Storage",
					ImageURL: "https://azure.microsoft.com/svghandler/storage/" +
						"?width=200",
					LongDescription: "Offload the heavy lifting of datacenter management" +
						" (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/storage/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Storage"},
			},
			m.serviceManager,
			service.NewPlan(service.PlanProperties{
				ID:   "6ddf6b41-fb60-4b70-af99-8ecc4896b3cf",
				Name: "general-purpose-storage-account",
				Description: "Azure general-purpose storage account; create your " +
					"own containers, files, and tables within this account",
				Free:      false,
				Stability: service.StabilityExperimental,
				Extended: map[string]interface{}{
					kindKey: storageKindGeneralPurposeStorageAcccount,
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "General Purpose Storage Account",
					Bullets: []string{"Azure general-purpose storage account",
						"Create your own containers, files, and tables within this account",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "800a17e1-f20a-463d-a290-20516052f647",
				Name: "blob-storage-account",
				Description: "Specialized Azure storage account for storing block " +
					"blobs and append blobs; create your own blob containers within " +
					"this account",
				Free:      false,
				Stability: service.StabilityExperimental,
				Extended: map[string]interface{}{
					kindKey: storageKindBlobStorageAccount,
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Blob Storage Account",
					Bullets: []string{"Specialized Azure storage account for storing " +
						"block blobs and append blobs",
						"Create your own containers, files, and tables within this account",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "189d3b8f-8307-4b3f-8c74-03d069237f70",
				Name: "blob-container",
				Description: "A specialized Azure storage account for storing block " +
					"blobs and append blobs; automatically provisions a blob container " +
					" within the account",
				Free:      false,
				Stability: service.StabilityExperimental,
				Extended: map[string]interface{}{
					kindKey: storageKindBlobContainer,
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Blob Container",
					Bullets: []string{"A specialized Azure storage account for storing " +
						"block blobs and append blobs",
						"Automatically provisions a blob container within the account",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
		),
	}), nil
}
