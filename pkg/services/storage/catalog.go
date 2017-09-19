package storage

import (
	"github.com/Azure/azure-service-broker/pkg/service"
)

const kindKey = "kind"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "2e2fc314-37b6-4587-8127-8f9ee8b33fea",
				Name:        "azure-storage",
				Description: "Azure Storage",
				Bindable:    true,
				Tags:        []string{"Azure", "Storage"},
			},
			service.NewPlan(&service.PlanProperties{
				ID:   "6ddf6b41-fb60-4b70-af99-8ecc4896b3cf",
				Name: "general-purpose-storage-account",
				Description: "Azure general-purpose storage account; create your " +
					"own containers, files, and tables within this account",
				Free: false,
				Extended: map[string]interface{}{
					kindKey: storageKindGeneralPurposeStorageAcccount,
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:   "800a17e1-f20a-463d-a290-20516052f647",
				Name: "blob-storage-account",
				Description: "Specialized Azure storage account for storing block " +
					"blobs and append blobs; create your own blob containers within " +
					"this account",
				Free: false,
				Extended: map[string]interface{}{
					kindKey: storageKindBlobStorageAccount,
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:   "189d3b8f-8307-4b3f-8c74-03d069237f70",
				Name: "blob-container",
				Description: "A specialized Azure storage account for storing block " +
					"blobs and append blobs; automatically provisions a blob container " +
					" within the account",
				Free: false,
				Extended: map[string]interface{}{
					kindKey: storageKindBlobContainer,
				},
			}),
		),
	}), nil
}
