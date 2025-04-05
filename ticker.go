package ticker

import (
	"context"
	"sync/atomic"
	"time"
)

func Ticker(d time.Duration, f func()) (context.CancelFunc, context.Context) {
	var now = time.Now()
	var next = now.Truncate(d).Add(d)
	time.Sleep(next.Sub(now))
	var ticker = time.NewTicker(d)

	var count int64
	var ctx, stop = context.WithCancel(context.Background())
	var after, finish = context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				for count > 0 {
					time.Sleep(time.Millisecond)
				}
				finish()
				return
			case <-ticker.C:
				go func() {
					atomic.AddInt64(&count, 1)
					defer atomic.AddInt64(&count, -1)
					f()
				}()
			}
		}
	}()
	return stop, after
}
