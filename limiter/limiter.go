package limiter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	limit    int64
	interval time.Duration
	client   *redis.Client
}

type Config struct {
	Limit    int64
	Interval time.Duration
	Address  string
}

func New(config Config) *Limiter {
	if config.Interval == 0 {
		config.Interval = time.Second
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: "",
		DB:       0,
	})
	return &Limiter{
		limit:    config.Limit,
		interval: config.Interval,
		client:   client,
	}
}

func (l *Limiter) Check(ctx context.Context, userID int64) (bool, error) {
	key := strconv.FormatInt(userID, 10)
	now := time.Now().UnixNano()
	window := l.interval.Nanoseconds()

	pipe := l.client.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-window))
	pipe.ZCard(ctx, key)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count, err := l.client.ZCard(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count >= l.limit {
		return false, nil
	}

	_, err = l.client.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: now,
	}).Result()
	if err != nil {
		return false, err
	}

	_, err = l.client.Expire(ctx, key, l.interval).Result()
	if err != nil {
		return false, err
	}

	return true, nil
}
