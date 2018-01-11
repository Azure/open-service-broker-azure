package broker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async/model"
	log "github.com/Sirupsen/logrus"
)

func (b *broker) doWaitForChildrenStep(
	ctx context.Context,
	args map[string]string,
) error {

	_, cancel := context.WithCancel(ctx)
	defer cancel()

	//time.Sleep(60 * time.Second)
	deprovisionFirstStep, ok := args["deprovisionFirstStep"]
	if !ok {
		return errors.New(`missing required argument "deprovisionFirstStep"`)
	}
	instanceID, ok := args["instanceID"]
	if !ok {
		return errors.New(`missing required argument "instanceID"`)
	}
	instance, ok, err := b.store.GetInstance(instanceID)
	if !ok {
		return b.handleDeprovisioningError(
			instanceID,
			"waitForParent",
			nil,
			"error loading persisted instance",
		)
	}
	if err != nil {
		return b.handleDeprovisioningError(
			instanceID,
			"waitForParent",
			err,
			"error loading persisted instance",
		)
	}
	childCount, err := b.store.GetInstanceChildCountByAlias(instance.Alias)
	if err != nil {
		log.WithFields(log.Fields{
			"step":       "waitforParent",
			"instanceID": instanceID,
			"error":      err,
		}).Error(
			"deprovisioning error: error determining child count",
		)
		return b.handleDeprovisioningError(
			instance,
			"waitForChilren",
			err,
			"error determining child count",
		)
	}
	var task model.Task
	if childCount > 0 {
		//Put this task back into the queue
		task = model.NewDelayedTask(
			"waitForChildrenStep",
			map[string]string{
				"deprovisionFirstStep": deprovisionFirstStep,
				"instanceID":           instanceID,
			},
			time.Minute*5,
		)
		log.WithFields(log.Fields{
			"step":       "waitforChildrenStep",
			"instanceID": instanceID,
		}).Debug("children not deprovisioned, will wait again")
	} else {
		//Put the real deprovision task into the queue
		task = model.NewTask(
			"deprovisionStep",
			map[string]string{
				"stepName":   deprovisionFirstStep,
				"instanceID": instanceID,
			},
		)
		log.WithFields(log.Fields{
			"step":       "waitforChildrenStep",
			"instanceID": instanceID,
		}).Debug("children deprovisioned,  sending start deprovision task")
	}
	if err = b.asyncEngine.SubmitTask(task); err != nil {
		return b.handleDeprovisioningError(
			instance,
			"waitForChildrenStep",
			err,
			fmt.Sprintf(`error transitioning deprovision to: "%s"`,
				deprovisionFirstStep),
		)
	}

	return nil
}
