package time

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseClock(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	clock := ParseClock("")
	assert.Nil(clock)

	clock = ParseClock("17:01:10")
	require.NotNil(clock)
	assert.EqualValues(17, clock.Hour)
	assert.EqualValues(1, clock.Min)
	assert.EqualValues(10, clock.Sec)

	require.PanicsWithValue("unsupported clock format: 01:10", func() {
		ParseClock("01:10")
	})

	require.Panics(func() {
		ParseClock("tt:01:10")
	})
}
