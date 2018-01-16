package broker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async/model"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
)

func (b *broker) doCheckChildrenStatuses(
	ctx context.Context,
	args map[string]string,
) error {
	instanceID, ok := args["instanceID"]
	if !ok {
		return errors.New(`missing required argument "instanceID"`)
	}
	instance, ok, err := b.store.GetInstance(instanceID)
	if !ok {
		return b.handleDeprovisioningError(
			instanceID,
			"checkChildrenStatus",
			nil,
			"error loading persisted instance",
		)
	}
	if err != nil {
		return b.handleDeprovisioningError(
			instanceID,
			"checkChildrenStatus",
			err,
			"error loading persisted instance",
		)
	}
	childCount, err := b.store.GetInstanceChildCountByAlias(instance.Alias)
	if err != nil {
		log.WithFields(log.Fields{
			"step":       "checkChildrenStatus",
			"instanceID": instanceID,
			"error":      err,
		}).Error(
			"deprovisioning error: error determining child count",
		)
		return b.handleDeprovisioningError(
			instance,
			"checkChildrenStatus",
			err,
			"error determining child count",
		)
	}
	var task model.Task
	if childCount > 0 {
		//Put this task back into the queue
		task = model.NewDelayedTask(
			"checkChildrenStatus",
			map[string]string{
				"instanceID": instanceID,
			},
			time.Minute*1,
		)
		log.WithFields(log.Fields{
			"instanceID": instanceID,
			"provisionedChildren" : childCount,
		}).Debug("children not deprovisioned, will wait again")
	} else {
		svc, ok := b.catalog.GetService(instance.ServiceID)
		if !ok {
			// If we don't find the Service in the catalog, something is really wrong.
			// (It should exist, because an instance with this serviceID exists.)
			log.WithFields(log.Fields{
				"step":       "checkChildrenStatus",
				"instanceID": instanceID,
				"serviceID":  instance.ServiceID,
			}).Error(
				"deprovisioning error: no Service found for serviceID",
			)
			return b.handleDeprovisioningError(
				instance,
				"checkChildrenStatus",
				nil,
				"error: no Service found for serviceID",
			)
		}
		plan, ok := svc.GetPlan(instance.PlanID)
		if !ok {
			// If we don't find the Service in the catalog, something is really wrong.
			// (It should exist, because an instance with this serviceID exists.)
			log.WithFields(log.Fields{
				"instanceID": instanceID,
				"serviceID":  instance.ServiceID,
				"planID":     instance.PlanID,
			}).Error(
				"deprovisioning error: no Plan found for planID in Service",
			)
			return b.handleDeprovisioningError(
				instance,
				"checkChildrenStatus",
				nil,
				"error: no Plan found for planID",
			)
		}
		serviceManager := svc.GetServiceManager()
		var deprovisioner service.Deprovisioner
		deprovisioner, err = serviceManager.GetDeprovisioner(plan)
		if err != nil {
			log.WithFields(log.Fields{
				"instanceID": instanceID,
				"serviceID":  instance.ServiceID,
				"planID":     instance.PlanID,
				"error":      err,
			}).Error(
				"deprovisioning error: error retrieving deprovisioner for " +
					"service and plan",
			)
			return b.handleDeprovisioningError(
				instance,
				"checkChildrenStatus",
				err,
				"error retrieving deprovisioner for service and service",
			)
		}
		deprovisionFirstStep, ok := deprovisioner.GetFirstStepName()
		if !ok {
			log.WithFields(log.Fields{
				"instanceID": instanceID,
				"serviceID":  instance.ServiceID,
				"planID":     instance.PlanID,
			}).Error(
				"pre-deprovisioning error: no steps found for deprovisioning " +
					"service and plan",
			)
			return b.handleDeprovisioningError(
				instance,
				"checkChildrenStatus",
				nil,
				"error: no steps found for deprovisioning service ance plan",
			)
		}

		//Put the real deprovision task into the queue
		task = model.NewTask(
			"deprovisionStep",
			map[string]string{
				"stepName":   deprovisionFirstStep,
				"instanceID": instanceID,
			},
		)
		log.WithFields(log.Fields{
			"step":       "checkChildrenStatus",
			"instanceID": instanceID,
		}).Debug("children deprovisioned,  sending start deprovision task")
	}
	if err = b.asyncEngine.SubmitTask(task); err != nil {
		return b.handleDeprovisioningError(
			instance,
			"checkChildrenStatus",
			err,
			fmt.Sprintf(
				`error submitting task %s from checkChildrenStatuses`,
				task.GetJobName(),
			),
		)
	}

	return nil
}
