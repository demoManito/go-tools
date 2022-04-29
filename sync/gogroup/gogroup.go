package gogroup

import (
	"context"
	"sync"
)

// GoGroup
// A Group is a collection of goroutines working on subtasks that are part of the same overall task.
// A zero Group is valid and does not cancel on error.
type GoGroup struct {
	cancel func()

	wg      *sync.WaitGroup
	errOnce *sync.Once

	err error
}

// WithContext new GoGroup and new cancel ctx
func WithContext(ctx context.Context) (*GoGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	goGroup := &GoGroup{
		cancel:  cancel,
		wg:      new(sync.WaitGroup),
		errOnce: new(sync.Once),
	}
	return goGroup, ctx
}

// Wait wg wait and ctx cancel
func (gg *GoGroup) Wait() error {
	gg.wg.Wait()
	if gg.cancel != nil {
		gg.cancel()
	}
	return gg.err
}

// Go go run goroutine
// any execution of error will cancel ctx, and error will be returned to you in Wait()
func (gg *GoGroup) Go(f func() error) {
	gg.wg.Add(1)

	go func() {
		defer gg.wg.Done()
		if err := f(); err != nil {
			gg.errOnce.Do(func() {
				gg.err = err
				if gg.cancel != nil {
					gg.cancel()
				}
			})
		}
	}()
}
