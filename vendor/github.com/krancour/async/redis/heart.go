package redis

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
)

// runHeartFn defines functions used to implement a beating heart
type runHeartFn func(ctx context.Context, interval time.Duration) error

// heartbeatFn defines functions used to implement a single heartbeat
type heartbeatFn func(ttl time.Duration) error

func (e *engine) defaultRunHeart(
	ctx context.Context,
	interval time.Duration,
) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if err := e.heartbeat(interval * 2); err != nil {
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

func (e *engine) defaultHeartbeat(ttl time.Duration) error {
	key := getHeartbeatKey(e.workerID)
	err := e.redisClient.Set(key, aliveIndicator, ttl).Err()
	if err != nil {
		return fmt.Errorf(
			"error sending heartbeat for worker %s: %s",
			e.workerID,
			err,
		)
	}
	return nil
}

func getHeartbeatKey(workerID string) string {
	return fmt.Sprintf("heartbeats:%s", workerID)
}
