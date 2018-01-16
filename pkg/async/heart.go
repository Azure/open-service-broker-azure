package async

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

// Heart is an interface to be implemented by components that can send worker
// heartbeats
type Heart interface {
	// Beat sends a single heartbeat
	Beat() error
	// Run sends heartbeats at regular intervals.  It blocks until a fatal error
	// is encountered or the context passed to it has been canceled. Run always
	// returns a non-nil error.
	Run(context.Context) error
}

// heart is a Redis-based implementation of the Heart interface
type heart struct {
	workerID    string
	frequency   time.Duration
	ttl         time.Duration
	redisClient *redis.Client
	// This allows tests to inject an alternative implementation of this function
	beat func() error
}

// newHeart returns a new Redis-based implementation of the Heart interface
func newHeart(
	workerID string,
	frequency time.Duration,
	redisClient *redis.Client,
) Heart {
	h := &heart{
		workerID:    workerID,
		frequency:   frequency,
		ttl:         frequency * 2,
		redisClient: redisClient,
	}
	h.beat = h.defaultBeat
	return h
}

// Beat sends a single heartbeat
func (h *heart) Beat() error {
	if err := h.beat(); err != nil {
		return &errHeartbeat{workerID: h.workerID, err: err}
	}
	return nil
}

// Run sends heartbeats at regular intervals.  It blocks until a fatal error is
// encountered or the context passed to it has been canceled. Run always returns
// a non-nil error.
func (h *heart) Run(ctx context.Context) error {
	ticker := time.NewTicker(h.frequency)
	defer ticker.Stop()
	for {
		if err := h.Beat(); err != nil {
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

// This is the default function for sending a heartbeat. It can be overridden
// to facilitate testing.
func (h *heart) defaultBeat() error {
	key := getHeartbeatKey(h.workerID)
	err := h.redisClient.Set(key, aliveIndicator, h.ttl).Err()
	if err != nil {
		return fmt.Errorf(
			"error sending heartbeat for worker %s: %s",
			h.workerID,
			err,
		)
	}
	return nil
}

func getHeartbeatKey(workerID string) string {
	return fmt.Sprintf("heartbeats:%s", workerID)
}
