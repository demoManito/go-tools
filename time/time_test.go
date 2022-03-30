package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeNowAndSetTimeNow(t *testing.T) {
	assert := assert.New(t)

	now := TimeNow()
	assert.LessOrEqual(now.Unix(), time.Now().Unix())

	SetTimeNow(func() time.Time {
		return time.Unix(1, 0)
	})
	now = TimeNow()
	assert.EqualValues(1, now.Unix())

	SetTimeNow(nil)
	now = TimeNow()
	assert.LessOrEqual(now.Unix(), time.Now().Unix())
}
