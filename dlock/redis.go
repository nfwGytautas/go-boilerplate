package dlock

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisLock an implementation of DLock interface using redis
type redisLock struct {
	rdb   redis.UniversalClient
	key   string
	ttl   time.Duration
	owner atomic.Bool
}

// NewRedisLock returns a new redis distributed lock
func NewRedisLock(rdb redis.UniversalClient, key string, ttl time.Duration) DLock {
	return &redisLock{
		rdb: rdb,
		key: key,
		ttl: ttl,
	}
}

func (rl *redisLock) Acquire(ctx context.Context) error {
	result, err := rl.rdb.SetArgs(ctx, rl.key, "locked", redis.SetArgs{
		Mode: "NX",
		TTL:  rl.ttl,
	}).Result()
	if err != nil {
		return err
	}

	if result != "OK" {
		return ErrLocked
	}

	rl.owner.Store(true)

	return nil
}

func (rl *redisLock) Release(ctx context.Context) error {
	if !rl.owner.Load() {
		return ErrNonOwnerRelease
	}

	defer rl.owner.Store(false)

	return rl.rdb.Del(ctx, rl.key).Err()
}
