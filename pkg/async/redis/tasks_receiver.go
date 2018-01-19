package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// receiveTasksFn defines functions used to receive tasks from one queue and
// dispatch them to another
type receiveTasksFn func(
	ctx context.Context,
	sourceQueueName string,
	destinationQueueName string,
	retCh chan []byte,
	errCh chan error,
)

// defaultReceive receives tasks from a source queue and dispatches them to a
// to both a destination queue and a return channel.
func (w *worker) defaultReceiveTasks(
	ctx context.Context,
	sourceQueueName string,
	destinationQueueName string,
	retCh chan []byte,
	errCh chan error,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		taskJSON, err := w.redisClient.BRPopLPush(
			sourceQueueName,
			destinationQueueName,
			time.Second*5,
		).Bytes()
		if err == redis.Nil {
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
		if err != nil {
			select {
			case errCh <- fmt.Errorf(
				`error receiving task from queue "%s": %s`,
				sourceQueueName,
				err,
			):
				continue
			case <-ctx.Done():
				return
			}
		}
		select {
		case retCh <- taskJSON:
		case <-ctx.Done():
			return
		}
	}
}
