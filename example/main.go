package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/LKarataev/flood-control-task/internal/limiter"
)

const callFreq = 400 * time.Millisecond

type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}

func main() {
	config := parseConfig()
	ctx := context.Background()
	var userID int64 = 555

	var fc FloodControl
	fc = limiter.New(config)

	for i := 0; i < 20; i++ {
		ok, err := fc.Check(ctx, userID)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if ok {
			fmt.Printf("Call [%d] - Passed\n", i)
		} else {
			fmt.Printf("Call [%d] - Failed\n", i)
		}
		time.Sleep(callFreq)
	}
}

func parseConfig() Config {
	var cfg Config
	var secInterval int
	flag.StringVar(&cfg.Address, "redis", "localhost:6379", `redis address (default: "localhost:6379")`)
	flag.IntVar(&secInterval, "interval", 5, "Interval in seconds (N)")
	flag.IntVar(&cfg.RateLimit, "limit", 5, "Maximum calls in N seconds (K)")
	flag.Parse()

	cfg.Interval = time.Duration(secInterval) * time.Second
	cfg.BurstLimit = int(cfg.RateLimit)
	return cfg
}
