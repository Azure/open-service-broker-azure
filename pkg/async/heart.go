package async

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

const aliveIndicator = "alive"

type errHeartbeat struct {
	workerID string
	err      error
}

func (e *errHeartbeat) Error() string {
	return fmt.Sprintf(
		`error sending heartbeat for worker "%s": %s`,
		e.workerID,
		e.err,
	)
}

// Heart is an interface to be implemented by components that can send worker
// heartbeats
type Heart interface {
	// Beat sends a single heartbeat
	Beat() error
	// Start sends heartbeats at regular intervals
	Start(context.Context) error
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
	err := h.beat()
	if err != nil {
		return &errHeartbeat{workerID: h.workerID, err: err}
	}
	return nil
}

// Start sends heartbeats at regular intervals
func (h *heart) Start(ctx context.Context) error {
	ticker := time.NewTicker(h.frequency)
	for {
		err := h.Beat()
		if err != nil {
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

// This is the default function for sending a heartbeart. It can be overriden
// to facilitate testing.
func (h *heart) defaultBeat() error {
	statusCmd := h.redisClient.Set(h.workerID, aliveIndicator, h.ttl)
	if statusCmd.Err() != nil {
		return fmt.Errorf(
			"error sending heartbeat for worker %s: %s",
			h.workerID,
			statusCmd.Err(),
		)
	}
	return nil
}
