package worker

import (
	"context"
	"math"
	"math/rand"
	"time"
)

type RetryPolicy struct {
	Attempts int
	Base     time.Duration
	Max      time.Duration
}

func (p RetryPolicy) Delay(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	d := float64(p.Base) * math.Pow(2, float64(attempt))
	if p.Max > 0 && time.Duration(d) > p.Max {
		return p.Max
	}
	j := 1 + rand.Float64()*0.2
	return time.Duration(d * j)
}
func Retry(ctx context.Context, p RetryPolicy, fn func(context.Context) error) error {
	var e error
	for i := 0; i < p.Attempts; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if e = fn(ctx); e == nil {
			return nil
		}
		timer := time.NewTimer(p.Delay(i))
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
	}
	return e
}
