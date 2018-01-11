package async

import "fmt"

const (
	mainActiveWorkQueueName  = "activeWork"
	mainDelayedWorkQueueName = "delayedWork"
)

func getWorkerActiveQueueName(workerID string) string {
	return fmt.Sprintf("worker-active-queues:%s", workerID)
}

func getWorkerDelayedQueueName(workerID string) string {
	return fmt.Sprintf("worker-delayed-queues:%s", workerID)
}
