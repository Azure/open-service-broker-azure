package redis

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
)

// runHeartFn defines functions used to implement a beating heart
type runHeartFn func(ctx context.Context) error

// heartbeatFn defines functions used to implement a single heartbeat
type heartbeatFn func() error

func (w *worker) defaultRunHeart(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	for {
		if err := w.heartbeat(); err != nil {
			return err
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			log.Debug("context canceled; async worker heartbeat stopping")
			return ctx.Err()
		}
	}
}

func (w *worker) defaultHeartbeat() error {
	key := getHeartbeatKey(w.id)
	err := w.redisClient.Set(key, aliveIndicator, time.Second*60).Err()
	if err != nil {
		return fmt.Errorf(
			"error sending heartbeat for worker %s: %s",
			w.id,
			err,
		)
	}
	return nil
}

func getHeartbeatKey(workerID string) string {
	return fmt.Sprintf("heartbeats:%s", workerID)
}
