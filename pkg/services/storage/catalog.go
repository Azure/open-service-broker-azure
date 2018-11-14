package storage

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

const serviceGeneralPurposeV2 = "azure-storage-general-purpose-v2-storage-account"
const serviceGeneralPurposeV1 = "azure-storage-general-purpose-v1-storage-account"
const serviceBlobAllInOne = "azure-storage-blob-storage-account-and-container"
const serviceBlobAccount = "azure-storage-blob-storage-account"

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {

	return service.NewCatalog([]service.Service{
		service.NewService(
			service.ServiceProperties{
				ID:   "9a3e28fe-8c02-49da-9b35-1b054eb06c95",
				Name: serviceGeneralPurposeV2,
				Description: "Azure general purpose v2 storage account; create your " +
					"own containers, files, and tables within this account",
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
			m.generalPurposeV2Manager,
			service.NewPlan(service.PlanProperties{
				ID:   "bc4f766a-c372-479c-b0b4-bd9d0546b3ef",
				Name: "account",
				Description: "Azure general purpose v2 storage account; create your " +
					"own containers, files, and tables within this account",
				Free:      false,
				Stability: service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "General Purpose V2 Storage Account",
					Bullets: []string{"Azure general-purpose v2 storage account",
						"Create your own containers, files, and tables within this account",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(serviceGeneralPurposeV2),
						UpdatingParametersSchema:     generateUpdatingParamsSchema(serviceGeneralPurposeV2),
					},
				},
			}),
		),
		service.NewService(
			service.ServiceProperties{
				ID:   "d10ea062-b627-41e8-a240-543b60030694",
				Name: serviceGeneralPurposeV1,
				Description: "Azure general purpose v1 storage account; create your " +
					"own containers, files, and tables within this account",
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
			m.generalPurposeV1Manager,
			service.NewPlan(service.PlanProperties{
				ID:   "9364d013-3690-4ce5-b0a2-b43d9b970b02",
				Name: "account",
				Description: "General-purpose v1 accounts provide access to all " +
					"Azure Storage services, but may not have the latest features" +
					"or the lowest per gigabyte pricing",
				Free:      false,
				Stability: service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "General Purpose V1 Storage Account",
					Bullets: []string{"Azure general-purpose v1 storage account",
						"Create your own containers, files, and tables within this account",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(serviceGeneralPurposeV1),
						UpdatingParametersSchema:     generateUpdatingParamsSchema(serviceGeneralPurposeV1),
					},
				},
			}),
		),
		service.NewService(
			service.ServiceProperties{
				ID:   "1a5b4582-29a3-48c5-9cac-511fd8c52756",
				Name: serviceBlobAccount,
				Description: "Specialized Azure storage account for storing block " +
					"blobs and append blobs",
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
			m.blobAccountManager,
			service.NewPlan(service.PlanProperties{
				ID:   "98ae02ec-da21-4b09-b5e0-e2f9583d565c",
				Name: "account",
				Description: "Specialized Azure storage account for storing block " +
					"blobs and append blobs; create your own blob containers within " +
					"this account",
				Free:      false,
				Stability: service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Blob Storage Account",
					Bullets: []string{"Specialized Azure storage account for storing " +
						"block blobs and append blobs",
						"Create your own containers, files, and tables within this account",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(serviceBlobAccount),
						UpdatingParametersSchema:     generateUpdatingParamsSchema(serviceBlobAccount),
					},
				},
			}),
		),
		service.NewService(
			service.ServiceProperties{
				ID:   "d799916e-3faf-4bdf-a48b-bf5012a2d38c",
				Name: serviceBlobAllInOne,
				Description: "A specialized Azure storage account for storing block " +
					"blobs and append blobs; automatically provisions a blob container " +
					" within the account",
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
			m.blobAllInOneManager,
			service.NewPlan(service.PlanProperties{
				ID:   "6c3b587d-0f88-4112-982a-dbe541f30669",
				Name: "all-in-one",
				Description: "A specialized Azure storage account for storing block " +
					"blobs and append blobs; automatically provisions a blob container " +
					" within the account",
				Free:      false,
				Stability: service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Blob Container",
					Bullets: []string{"A specialized Azure storage account for storing " +
						"block blobs and append blobs",
						"Automatically provisions a blob container within the account",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(serviceBlobAllInOne),
						UpdatingParametersSchema:     generateUpdatingParamsSchema(serviceBlobAllInOne),
					},
				},
			}),
		),
	}), nil
}
