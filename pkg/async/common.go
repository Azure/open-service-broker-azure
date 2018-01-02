package async

import "fmt"

const mainWorkQueueName = "work"

func getWorkerQueueName(workerID string) string {
	return fmt.Sprintf("worker-queues:%s", workerID)
}
