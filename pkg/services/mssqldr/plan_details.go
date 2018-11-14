package mssqldr

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type planDetails interface {
	getProvisionSchema() service.InputParametersSchema
	getTierProvisionParameters(
		pp service.ProvisioningParameters,
	) (map[string]interface{}, error)
	getUpdateSchema() service.InputParametersSchema
	validateUpdateParameters(service.Instance) error
}

type dtuPlanDetails struct {
	tierName    string
	skuMap      map[int64]string
	allowedDTUs []int64
	defaultDTUs int64
	storageInGB int64
}

func (d dtuPlanDetails) validateUpdateParameters(service.Instance) error {
	return nil // no op
}

func (d dtuPlanDetails) getUpdateSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{},
	}
	// Basic tier is constrained to just 5 DTUs, so don't present this as an
	// option
	if len(d.allowedDTUs) > 0 {
		ips.PropertySchemas["dtus"] = &service.IntPropertySchema{
			Title:         "DTUs",
			AllowedValues: d.allowedDTUs,
			DefaultValue:  ptr.ToInt64(d.defaultDTUs),
			Description: "DTUs are a bundled measure of compute, " +
				"storage, and IO resources.",
		}
	}
	return ips
}

func (d dtuPlanDetails) getProvisionSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		RequiredProperties: []string{
			"failoverGroup",
			"database",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"failoverGroup": &service.StringPropertySchema{
				Title:       "Failover Group",
				Description: "The name of the failover group",
			},
			"database": &service.StringPropertySchema{
				Title:       "Database",
				Description: "The name of the database",
			},
		},
	}
	// Basic tier is constrained to just 5 DTUs, so don't present this as an
	// option
	if len(d.allowedDTUs) > 0 {
		ips.PropertySchemas["dtus"] = &service.IntPropertySchema{
			Title:         "DTUs",
			AllowedValues: d.allowedDTUs,
			DefaultValue:  ptr.ToInt64(d.defaultDTUs),
			Description: "DTUs are a bundled measure of compute, " +
				"storage, and IO resources.",
		}
	}
	return ips
}

func (d dtuPlanDetails) getTierProvisionParameters(
	pp service.ProvisioningParameters,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["sku"] = d.getSKU(pp)
	p["tier"] = d.tierName
	// ARM template needs bytes
	p["maxSizeBytes"] = pp.GetInt64("storage") * 1024 * 1024 * 1024
	return p, nil
}

func (d dtuPlanDetails) getSKU(pp service.ProvisioningParameters) string {
	// Basic tier is constrained to just 5 DTUs, if this is the basic tier, there
	// is no dtus param. We can infer this is the case if the tier details don't
	// tell us there's a choice.
	if len(d.allowedDTUs) == 0 {
		return d.skuMap[d.defaultDTUs]
	}
	return d.skuMap[pp.GetInt64("dtus")]
}

type vCorePlanDetails struct {
	tierName      string
	tierShortName string
}

func (v vCorePlanDetails) validateUpdateParameters(
	instance service.Instance,
) error {
	return validateStorageUpdate(
		*instance.ProvisioningParameters,
		*instance.UpdatingParameters,
	)
}

func (v vCorePlanDetails) getUpdateSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{},
	}
	ips.PropertySchemas["cores"] = &service.IntPropertySchema{
		Title:         "Cores",
		AllowedValues: []int64{2, 4, 8, 16, 24, 32, 48, 80},
		DefaultValue:  ptr.ToInt64(2),
		Description:   "A virtual core represents the logical CPU",
	}
	ips.PropertySchemas["storage"] = &service.IntPropertySchema{
		Title:        "Storage",
		MinValue:     ptr.ToInt64(5),
		MaxValue:     ptr.ToInt64(1024),
		DefaultValue: ptr.ToInt64(10),
		Description:  "The maximum data storage capacity (in GB)",
	}
	return ips
}

func (v vCorePlanDetails) getProvisionSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		RequiredProperties: []string{
			"failoverGroup",
			"database",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"failoverGroup": &service.StringPropertySchema{
				Title:       "FailoverGroup",
				Description: "The name of the failover group",
			},
			"database": &service.StringPropertySchema{
				Title:       "Database",
				Description: "The name of the database",
			},
			"cores": &service.IntPropertySchema{
				Title:         "Cores",
				AllowedValues: []int64{2, 4, 8, 16, 24, 32, 48, 80},
				DefaultValue:  ptr.ToInt64(2),
				Description:   "A virtual core represents the logical CPU",
			},
			"storage": &service.IntPropertySchema{
				Title:        "Storage",
				MinValue:     ptr.ToInt64(5),
				MaxValue:     ptr.ToInt64(1024),
				DefaultValue: ptr.ToInt64(10),
				Description:  "The maximum data storage capacity (in GB)",
			},
		},
	}
	return ips
}

func (v vCorePlanDetails) getTierProvisionParameters(
	pp service.ProvisioningParameters,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["sku"] = v.getSKU(pp)
	p["tier"] = v.tierName
	// ARM template needs bytes
	p["maxSizeBytes"] = pp.GetInt64("storage") * 1024 * 1024 * 1024
	return p, nil
}

func (v vCorePlanDetails) getSKU(pp service.ProvisioningParameters) string {
	return fmt.Sprintf(
		"%s_Gen5_%d",
		v.tierShortName,
		pp.GetInt64("cores"),
	)
}

func getDBMSPairRegisteredUpdateParamSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		SecureProperties: []string{
			"primaryAdministratorLoginPassword",
			"secondaryAdministratorLoginPassword",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"primaryAdministratorLogin": &service.StringPropertySchema{
				Title: "Primary Administrator Login",
				Description: "Specifies the administrator login name" +
					" of the primary existing server",
			},
			"primaryAdministratorLoginPassword": &service.StringPropertySchema{
				Title: "Primary Administrator Login Password",
				Description: "Specifies the administrator login name" +
					" of the existing primary server",
			},
			"secondaryAdministratorLogin": &service.StringPropertySchema{
				Title: "Secondary Administrator Login",
				Description: "Specifies the administrator login name" +
					" of the secondary existing server",
			},
			"secondaryAdministratorLoginPassword": &service.StringPropertySchema{
				Title: "Secondary Administrator Login Password",
				Description: "Specifies the administrator login name" +
					" of the existing secondary server",
			},
		},
	}
}

func getDBMSPairRegisteredProvisionParamSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		RequiredProperties: []string{
			"primaryResourceGroup",
			"primaryLocation",
			"primaryServer",
			"primaryAdministratorLogin",
			"primaryAdministratorLoginPassword",
			"secondaryResourceGroup",
			"secondaryLocation",
			"secondaryServer",
			"secondaryAdministratorLogin",
			"secondaryAdministratorLoginPassword",
		},
		SecureProperties: []string{
			"primaryAdministratorLoginPassword",
			"secondaryAdministratorLoginPassword",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"primaryResourceGroup": &service.StringPropertySchema{
				Title: "Primary Resource Group",
				Description: "Specifies the resource group of " +
					"the primary existing server",
			},
			"primaryLocation": &service.StringPropertySchema{
				Title: "Primary Location",
				Description: "Specifies the location of " +
					"the primary existing server",
				CustomPropertyValidator: azure.LocationValidator,
			},
			"primaryServer": &service.StringPropertySchema{
				Title:       "Primary Server",
				Description: "Specifies the name of the primary existing server",
			},
			"primaryAdministratorLogin": &service.StringPropertySchema{
				Title: "Primary Administrator Login",
				Description: "Specifies the administrator login name" +
					" of the primary existing server",
			},
			"primaryAdministratorLoginPassword": &service.StringPropertySchema{
				Title: "Primary Administrator Login Password",
				Description: "Specifies the administrator login password" +
					" of the primary existing server",
			},
			"secondaryResourceGroup": &service.StringPropertySchema{
				Title: "Secondary Resource Group",
				Description: "Specifies the resource group of " +
					"the secondary existing server",
			},
			"secondaryLocation": &service.StringPropertySchema{
				Title: "Secondary Location",
				Description: "Specifies the location of " +
					"the secondary existing server",
				CustomPropertyValidator: azure.LocationValidator,
			},
			"secondaryServer": &service.StringPropertySchema{
				Title:       "Secondary Server",
				Description: "Specifies the name of the secondary existing server",
			},
			"secondaryAdministratorLogin": &service.StringPropertySchema{
				Title: "Secondary Administrator Login",
				Description: "Specifies the administrator login name" +
					" of the secondary existing server",
			},
			"secondaryAdministratorLoginPassword": &service.StringPropertySchema{
				Title: "Secondary Administrator Login Password",
				Description: "Specifies the administrator login password" +
					" of the secondary existing server",
			},
			"tags": &service.ObjectPropertySchema{
				Title: "Tags",
				Description: "Tags to be applied to resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
			},
		},
	}
}

func getDatabasePairRegisteredProvisionParamSchema() service.InputParametersSchema { // nolint: lll
	return service.InputParametersSchema{
		RequiredProperties: []string{
			"failoverGroup",
			"database",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"failoverGroup": &service.StringPropertySchema{
				Title:       "Failover Group",
				Description: "The name of the failover group",
			},
			"database": &service.StringPropertySchema{
				Title:       "Database",
				Description: "The name of the database",
			},
		},
	}
}

func validateStorageUpdate(
	pp service.ProvisioningParameters,
	up service.ProvisioningParameters,
) error {
	existingStorage := pp.GetInt64("storage")
	newStorge := up.GetInt64("storage")
	if newStorge < existingStorage {
		return service.NewValidationError(
			"storage",
			fmt.Sprintf(
				`invalid value: cannot reduce storage from %d to %d`,
				existingStorage,
				newStorge,
			),
		)
	}
	return nil
}
