package gogroup

import (
	"context"
	"errors"
	"testing"
)

func TestGo(t *testing.T) {
	gg, ctx := WithContext(context.Background())

	for i := 0; i < 100; i++ {
		func(i int) {
			gg.Go(func() error {
				if i == 10 {
					t.Log("error")
					return errors.New("test error")
				}

				select {
				case <-ctx.Done():
					t.Log("ctx done")
					return nil
				default:
					t.Logf("ctx err: %s, index: %d", ctx.Err(), i)
				}
				return nil
			})
		}(i)
	}
	if err := gg.Wait(); err != nil {
		t.Log(err)
	}
}
