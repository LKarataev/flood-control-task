package main

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	rate     int
	burst    int
	interval time.Duration
	client *redis.Client
}

type Config struct {
	Address    string
	RateLimit  int
	BurstLimit int
	Interval   time.Duration
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
		rate:     config.RateLimit,
		burst:    config.BurstLimit,
		interval: config.Interval,
		client:   client,
	}
}

func (l *Limiter) Check(ctx context.Context, userID int64) (bool, error) {
	key := fmt.Sprintf("%d", userID)

    exists, err := l.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

    if exists == 0 {
		now := time.Now().Truncate(l.interval).Unix()
		_, err := l.client.LPush(ctx, key, now, float64(l.burst-1)).Result()
		if err != nil {
			return false, err
		}
		return true, nil
	}

	resp, err := l.client.LRange(ctx, key, 0, 1).Result()
	if err != nil {
		return false, err
	}

	last, _ := strconv.ParseInt(resp[0], 10, 64)
	tokens, _ := strconv.ParseFloat(resp[1], 64)

	since := time.Since(time.Unix(last, 0)).Truncate(l.interval)
	quota := since.Seconds() / l.interval.Seconds() * float64(l.rate)

	tokens = math.Min(tokens+quota, float64(l.burst))

	if tokens < 1.0 {
		return false, nil
	}

	tokens -= 1.0

	now := time.Now().Truncate(l.interval).Unix()

	_, err = l.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.LSet(ctx, key, 0, now)
		pipe.LSet(ctx, key, 1, tokens)
		return nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
