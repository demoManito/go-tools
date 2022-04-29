package limiter

import (
	"context"
	"sort"
	"time"

	"golang.org/x/time/rate"
)

// Per 返回一个事件所需等待的最细粒度
func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

// ILimiter
type ILimiter interface {
	Wait(ctx context.Context) error
	WaitN(ctx context.Context, n int) error
	Limit() rate.Limit
}

// Limiter
type Limiter struct {
	limiters []ILimiter
}

// New limiters 兼容不同细度的限流器
func New(limiters ...ILimiter) *Limiter {
	sort.Slice(limiters, func(i, j int) bool {
		return limiters[i].Limit() < limiters[i].Limit()
	})
	return &Limiter{limiters: limiters}
}

// Wait consume a token, is shorthand for WaitN(ctx, 1).
func (l *Limiter) Wait(ctx context.Context) error {
	return l.WaitN(ctx, 1)
}

// WaitN consume n token
func (l *Limiter) WaitN(ctx context.Context, n int) error {
	for _, limiter := range l.limiters {
		if err := limiter.WaitN(ctx, n); err != nil {
			return err
		}
	}
	return nil
}

// Limit limiter max limit
func (l *Limiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}
