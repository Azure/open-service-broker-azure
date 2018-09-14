package rediscache

import "github.com/Azure/open-service-broker-azure/pkg/service"

const basic = "basic"
const standard = "standard"
const premium = "premium"

func (m *module) GetCatalog() (service.Catalog, error) {
	bpd := planDetail{
		planName:        basic,
		allowedCapacity: []int64{0, 1, 2, 3, 4, 5, 6},
	}
	spd := planDetail{
		planName:        standard,
		allowedCapacity: []int64{0, 1, 2, 3, 4, 5, 6},
	}
	ppd := planDetail{
		planName:          premium,
		allowedCapacity:   []int64{1, 2, 3, 4},
		allowedShardCount: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	return service.NewCatalog([]service.Service{
		service.NewService(
			service.ServiceProperties{
				ID:          "0346088a-d4b2-4478-aa32-f18e295ec1d9",
				Name:        "azure-rediscache",
				Description: "Azure Redis Cache (Experimental)",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure Redis Cache",
					ImageURL: "https://azure.microsoft.com/svghandler/redis-cache/" +
						"?width=200",
					LongDescription: "High throughput and consistent low-latency data " +
						"access to power fast, scalable Azure applications (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/redis-cache/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Redis", "Cache", "Database"},
			},
			m.serviceManager,
			service.NewPlan(service.PlanProperties{
				ID:          "362b3d1b-5b57-4289-80ad-4a15a760c29c",
				Name:        "basic",
				Description: "Basic Tier, 250MB Cache",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Extended: map[string]interface{}{
					"redisCacheSKU":    "Basic",
					"redisCacheFamily": "C",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets:     []string{"250MB Cache"},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: bpd.getProvisioningParamsSchema(),
						UpdatingParametersSchema:     bpd.getUpdatingParamsSchema(),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:          "4af8bbd1-962d-4e26-84f1-f72d1d959d87",
				Name:        "standard",
				Description: "Standard Tier, 1GB Cache",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Extended: map[string]interface{}{
					"redisCacheSKU":    "Standard",
					"redisCacheFamily": "C",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets:     []string{"1GB Cache"},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: spd.getProvisioningParamsSchema(),
						UpdatingParametersSchema:     spd.getUpdatingParamsSchema(),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:          "b1057a8f-9a01-423a-bc35-e168d5c04cf0",
				Name:        "premium",
				Description: "Premium Tier, 6GB Cache",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Extended: map[string]interface{}{
					"redisCacheSKU":    "Premium",
					"redisCacheFamily": "P",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Premium Tier",
					Bullets:     []string{"6GB Cache"},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: ppd.getProvisioningParamsSchema(),
						UpdatingParametersSchema:     ppd.getUpdatingParamsSchema(),
					},
				},
			}),
		),
	}), nil
}
