package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/LKarataev/flood-control-task/limiter"
)

const callFreq = 1000 * time.Millisecond

type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}

func parseConfig() limiter.Config {
	var config limiter.Config
	var interval int
	flag.StringVar(&config.Address, "redis", "localhost:6379", `redis address`)
	flag.IntVar(&interval, "interval", 5, "Interval in seconds (N)")
	flag.Int64Var(&config.Limit, "limit", 5, "Maximum calls in N seconds (K)")
	flag.Parse()
	config.Interval = time.Duration(interval) * time.Second
	return config
}

func main() {
	config := parseConfig()
	ctx := context.Background()
	var userID int64 = 555

	var fc FloodControl = limiter.New(config)

	fmt.Println("Interval (N):", config.Interval)
	fmt.Println("Limit calls (K):", config.Limit)
	fmt.Println("Call frequency:", callFreq)
	for i := 0; i < 15; i++ {
		ok, err := fc.Check(ctx, userID)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if ok {
			fmt.Printf("Call [%d] -\x1b[32m Passed \x1b[0m\n", i)
		} else {
			fmt.Printf("Call [%d] -\x1b[33m Failed \x1b[0m\n", i)
		}
		time.Sleep(callFreq)
	}
	fmt.Println("End of example!")
}
