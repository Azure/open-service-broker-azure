package api

import (
	"log"
	"net/http"

	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/echo"
)

// TODO: Delete this later-- this is just to aid in hacking
func (s *server) test(w http.ResponseWriter, r *http.Request) {
	instance := &service.Instance{
		InstanceID: "88fd1ae6-817d-4138-ae9a-15bda06fbea0",
		ServiceID:  "470b4bb6-8603-432d-aa34-d2ee74d7966c",
		PlanID:     "39ce8f26-d87d-4fb7-b06b-56f48215e308",
		Status:     service.InstanceStateProvisioning,
	}
	module := echo.New()
	err := instance.SetProvisioningParameters(module.GetEmptyProvisioningParameters())
	if err != nil {
		log.Println("error encoding empty provisioning parameters")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = instance.SetProvisioningResult(module.GetEmptyProvisioningResult())
	if err != nil {
		log.Println("error encoding empty provisioning result")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.store.WriteInstance(instance)
	if err != nil {
		log.Println("error storing new instance")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.asyncEngine.Provision(instance)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
