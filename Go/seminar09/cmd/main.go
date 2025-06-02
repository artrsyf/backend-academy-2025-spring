package main

import (
	"context"
	"log"
	"math/rand/v2"
	"time"

	"github.com/go-co-op/gocron"
)

func Repeat(fn func(context.Context) bool) func(context.Context) {
	return func(ctx context.Context) {
		repeat := true

		for repeat {
			select {
			case <-ctx.Done():
				return
			default:
				repeat = fn(ctx)
			}
		}
	}
}

func main() {
	ctx := context.Background()

	s := gocron.NewScheduler(time.UTC)

	s.Every(20*time.Second).Do(Repeat(func(ctx context.Context) bool {
		log.Println("Processor 1")
		return false
	}), ctx)

	s.Every(10*time.Second).Do(Repeat(func(ctx context.Context) bool {
		log.Println("Processor 2")
		if rand.N(2) == 0 {
			return false
		}

		return true
	}), ctx)

	s.Cron("0 0 * * *").Do(func() {
		log.Println("Daily cleanup at midnight")
	})

	s.StartBlocking()
}
