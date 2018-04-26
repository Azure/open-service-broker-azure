package postgresql

import (
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

const (
	gen4TemplateString = "Gen4"
	gen5TemplateString = "Gen5"
	gen4ParamString    = "gen4"
	gen5ParamString    = "gen5"
)

type planSchema struct {
	defaultHardware         string
	allowedHardware         []string
	validCores              []int
	defaultCores            int
	tier                    string
	minStorage              int
	maxStorage              int
	defaultStorage          int
	allowedBackupRedundancy []string
	defaultBackupRedundancy string
	minBackupRetention      int
	maxBackupRetention      int
	defaultBackupRetention  int
}

func (p *planSchema) buildSku(
	pp dbmsProvisioningParameters,
) (string, error) {
	hardwareFamily, err := p.generateHardwareFamilyString(pp)
	if err != nil {
		return "", fmt.Errorf("error building sku: %s", err)
	}
	var cores int

	if pp.Cores == nil {
		cores = p.defaultCores
	} else {
		cores = *pp.Cores
	}
	//The name of the sku typically,
	// tier + family + cores, e.g. B_Gen4_1, GP_Gen5_8.
	sku := fmt.Sprintf("%s_%s_%d", p.tier, hardwareFamily, cores)
	return sku, nil
}

func (p *planSchema) generateHardwareFamilyString(
	pp dbmsProvisioningParameters,
) (string, error) {
	if pp.HardwareFamily == "" {
		return gen5TemplateString, nil
	} else if strings.ToLower(pp.HardwareFamily) == gen4ParamString {
		return gen4TemplateString, nil
	} else if strings.ToLower(pp.HardwareFamily) == gen5ParamString {
		return gen5TemplateString, nil
	}
	return "", fmt.Errorf("unknown hardware family")

}

func generateDBMSPlanSchema(
	schema planSchema,
) map[string]service.ParameterSchema {
	p := getDBMSCommonProvisionParamSchema()
	p["hardwareFamily"] = &service.SimpleParameterSchema{
		Type:          "string",
		Description:   "Specifies the compute generation to use for the DBMS",
		AllowedValues: schema.allowedHardware,
		Default:       schema.defaultHardware,
	}

	p["cores"] = &service.SimpleParameterSchema{
		Type: "number",
		Description: "Specifies vCores, which represent the logical " +
			"CPU of the underlying hardware",
		AllowedValues: schema.validCores,
		Default:       schema.defaultCores,
	}
	p["storage"] = &service.SimpleParameterSchema{
		Type:        "number",
		Description: "Specifies the storage in GBs",
		Default:     schema.defaultStorage,
	}
	p["backupRetention"] = &service.SimpleParameterSchema{
		Type:        "number",
		Description: "Specifies the number of days for backup retention",
		Default:     schema.minBackupRetention,
	}
	p["backupRedundancy"] = &service.SimpleParameterSchema{
		Type:          "string",
		Description:   "Specifies the backup redundancy",
		AllowedValues: schema.allowedBackupRedundancy,
		Default:       schema.defaultBackupRedundancy,
	}
	return p
}

func (p *planSchema) validateProvisionParameters(
	pp dbmsProvisioningParameters,
) error {
	//hardware family
	hardwareValid := false
	for _, v := range p.allowedHardware {
		if v == pp.HardwareFamily {
			hardwareValid = true
		}
	}
	if !hardwareValid {
		return service.NewValidationError(
			"hardwareFamily",
			fmt.Sprintf(`invalid value: "%s"`, pp.HardwareFamily))
	}

	//cores
	if pp.Cores != nil {
		coresValid := false
		coresParam := *pp.Cores
		for _, v := range p.validCores {
			if v == coresParam {
				coresValid = true
			}
		}
		if !coresValid {
			return service.NewValidationError(
				"cores",
				fmt.Sprintf(`invalid value : "%d"`, coresParam))
		}
	}
	//storage
	if pp.Storage != nil {
		storageParam := *pp.Storage
		if storageParam < p.minStorage || storageParam > p.maxStorage {
			return service.NewValidationError(
				"storage",
				fmt.Sprintf(`invalid value : "%d"`, storageParam))
		}
	}
	//backupRetation
	if pp.BackupRetention != nil {
		backupRetentionParam := *pp.BackupRetention
		if backupRetentionParam < p.minBackupRetention ||
			backupRetentionParam > p.maxBackupRetention {
			return service.NewValidationError(
				"backupRetention",
				fmt.Sprintf(`invalid value : "%d"`, backupRetentionParam))
		}
	}
	//backupRedundancy
	if pp.BackupRedundancy != "" {
		backupRedundancyValid := false
		for _, v := range p.allowedBackupRedundancy {
			if v == pp.BackupRedundancy {
				backupRedundancyValid = true
			}
		}
		if !backupRedundancyValid {
			return service.NewValidationError(
				"backupRedundancy",
				fmt.Sprintf(`invalid value : "%s"`, pp.BackupRedundancy))
		}
	}
	return nil
}

func (p *planSchema) getCores(pp dbmsProvisioningParameters) int {
	if pp.Cores != nil {
		return *pp.Cores
	}
	return p.defaultCores

}

func (p *planSchema) getStorage(pp dbmsProvisioningParameters) int {
	if pp.Storage != nil {
		return *pp.Storage
	}
	return p.defaultStorage

}

func (p *planSchema) getBackupRetention(pp dbmsProvisioningParameters) int {
	if pp.BackupRetention != nil {
		return *pp.BackupRetention
	}
	return p.defaultBackupRetention

}

func (p *planSchema) isGeoRedundentBackup(pp dbmsProvisioningParameters) bool {
	return pp.BackupRedundancy == "geo"
}

func (p *planSchema) getHardwareFamily(pp dbmsProvisioningParameters) string {
	if pp.HardwareFamily == "" {
		if p.defaultHardware == gen4ParamString {
			return gen4TemplateString
		}
		return gen5TemplateString
	} else if pp.HardwareFamily == gen4ParamString {
		return gen4TemplateString
	}
	return gen5TemplateString
}
