package iothub

import "github.com/Azure/open-service-broker-azure/pkg/service"

const (
	planS1 = "standard-s1"
	planS2 = "standard-s2"
	planS3 = "standard-s3"
	planB1 = "basic-b1"
	planB2 = "basic-b2"
	planB3 = "basic-b3"
	planF1 = "free"
)

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			service.ServiceProperties{
				ID:          "afd72c3b-6c2d-40f2-ad0d-d90467989be5",
				Name:        "azure-iot-hub",
				Description: "Azure IoT Hub (Experimental)",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure IoT Hub",
					ImageURL: "https://azure.microsoft.com/svghandler/iot-hub/" +
						"?width=200",
					LongDescription: "Securely connect, monitor, and manage billions of " +
						"devices to develop Internet of Things (IoT) applications. " +
						"IoT Hub is an open and flexible cloud platform as a service " +
						"that supports open-source SDKs and multiple protocols. (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/" +
						"iot-hub/",
					SupportURL: "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "IoT Hub", "IoT"},
			},
			m.iotHubManager,
			service.NewPlan(service.PlanProperties{
				ID:   "4d6c40dd-7525-4260-8e4d-f65818197c2b",
				Name: planF1,
				Description: "IoT hub Free Tier - max 8,000 " +
					"messages per day.",
				Free: true,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Free Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(planF1),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "bdff693c-39cb-4590-b4ce-d1a17fab5848",
				Name: planB1,
				Description: "IoT hub Basic B1 Tier - max 400,000 " +
					"messages per day.",
				Free: false,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Basic B1 Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(planB1),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "eaa7bebe-6b62-4471-b02a-d9f97094f894",
				Name: planB2,
				Description: "IoT hub Basic B2 Tier - max 6,000,000 " +
					"messages per day.",
				Free: false,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Basic B2 Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(planB2),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "d5b67100-d826-4ac3-bc30-6c01f4cc5c52",
				Name: planB3,
				Description: "IoT hub Basic B3 Tier - max 300,000,000 " +
					"messages per day.",
				Free: false,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Basic B3 Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(planB1),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "0dde7e80-1f32-470d-ba0b-9db4fe1826be",
				Name: planS1,
				Description: "IoT hub Standard S1 Tier - max 400,000 " +
					"messages per day.",
				Free: false,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Standard S1 Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(planS1),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "857a73b3-4a3a-44cd-b4fa-e74cab6bd4db",
				Name: planS2,
				Description: "IoT hub Standard S2 Tier - max 6,000,000 " +
					"messages per day.",
				Free: false,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Standard S2 Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(planS2),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "7f1264d8-4786-4121-be43-c1de31f1cb1e",
				Name: planS3,
				Description: "IoT hub Standard S3 Tier - max 300,000,000 " +
					"messages per day.",
				Free: false,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Standard S3 Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(planS3),
					},
				},
			}),
		),
	}), nil
}
