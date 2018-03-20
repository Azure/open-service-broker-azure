package fake

import "github.com/Azure/open-service-broker-azure/pkg/service"

type provisioningParameters struct {
	SomeParameter string `json:"someParameter"`
}

type bindingParameters struct {
	SomeParameter string `json:"someParameter"`
}

// SplitProvisioningParameters splits a map of provisioning parameters into
// two separate maps, with one containing non-sensitive provisioning parameters
// and the other containing sensitive provisioning parameters.
func (s *ServiceManager) SplitProvisioningParameters(
	cpp service.CombinedProvisioningParameters,
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
) {
	pp := provisioningParameters{}
	if err := service.GetStructFromMap(cpp, &pp); err != nil {
		return nil, nil, err
	}
	ppMap, err := service.GetMapFromStruct(pp)
	if err != nil {
		return nil, nil, err
	}
	return ppMap, nil, nil
}

// SplitBindingParameters splits a map of binding parameters into two separate
// maps, with one containing non-sensitive binding parameters and the other
// containing sensitive binding parameters.
func (s *ServiceManager) SplitBindingParameters(
	cbp service.CombinedBindingParameters,
) (
	service.BindingParameters,
	service.SecureBindingParameters,
	error,
) {
	bp := bindingParameters{}
	err := service.GetStructFromMap(cbp, &bp)
	if err != nil {
		return nil, nil, err
	}
	bpMap, err := service.GetMapFromStruct(bp)
	if err != nil {
		return nil, nil, err
	}
	return bpMap, nil, nil
}
