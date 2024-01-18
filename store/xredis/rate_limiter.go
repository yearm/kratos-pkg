package xredis

import (
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

// RedisLimiter ...
type RedisLimiter struct {
	redisClient *redis.Client
	key         string
	period      time.Duration
	maxCount    int64
}

// NewRedisLimiter ...
func NewRedisLimiter(redisDB RedisDB, key string, period time.Duration, maxCount int64) *RedisLimiter {
	return &RedisLimiter{
		redisClient: redisDB,
		key:         key,
		period:      period,
		maxCount:    maxCount,
	}
}

// IsActionAllow ...
func (r *RedisLimiter) IsActionAllow() (bool, error) {
	pipe := r.redisClient.Pipeline()
	now := time.Now()
	beforeTime := now.Add(-r.period).UnixNano()
	pipe.ZAdd(r.key, &redis.Z{
		Score:  float64(now.UnixNano()),
		Member: now.UnixNano(),
	})
	pipe.ZRemRangeByScore(r.key, "0", strconv.FormatInt(beforeTime, 10))
	result := pipe.ZCard(r.key)
	pipe.Expire(r.key, r.period)
	if _, err := pipe.Exec(); err != nil {
		return false, err
	}
	return result.Val() <= r.maxCount, nil
}
