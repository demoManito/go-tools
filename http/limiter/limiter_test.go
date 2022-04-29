package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"

	ttime "github.com/demoManito/go-tools/time"
)

func TestLimiter_Wait(t *testing.T) {
	require := require.New(t)

	ctx := context.Background()
	now := ttime.TimeNow()
	limiter := New(rate.NewLimiter(Per(2, time.Second), 2))

	_ = limiter.Wait(ctx)
	t.Log(ttime.TimeNow().Sub(now).Seconds())
	now = ttime.TimeNow()

	_ = limiter.Wait(ctx)
	t.Log(ttime.TimeNow().Sub(now).Seconds())
	now = ttime.TimeNow()

	_ = limiter.Wait(ctx)
	require.LessOrEqual(1/2*time.Second.Seconds(), ttime.TimeNow().Sub(now).Seconds())
	now = ttime.TimeNow()

	_ = limiter.Wait(ctx)
	require.LessOrEqual(1/2*time.Second.Seconds(), ttime.TimeNow().Sub(now).Seconds())
	now = ttime.TimeNow()

	_ = limiter.Wait(ctx)
	require.LessOrEqual(1/2*time.Second.Seconds(), ttime.TimeNow().Sub(now).Seconds())
}
