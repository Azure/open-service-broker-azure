package async

import "testing"
import "context"
import "time"
import "github.com/stretchr/testify/assert"

func TestCleanerCleanBlocksUntilCleanErrors(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)
	c.clean = func() error {
		return errSome
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := c.Clean(ctx)
	assert.Equal(t, &errCleaning{err: errSome}, err)
}

func TestCleanerCleanBlocksUntilContextCanceled(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)
	c.clean = func() error {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := c.Clean(ctx)
	assert.Equal(t, ctx.Err(), err)
}

func TestCleanerClean(t *testing.T) {
	// TODO: Implement this
}

func TestCleanerCleanWorker(t *testing.T) {
	// TODO: Implement this
}
