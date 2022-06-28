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

func TestClock_LessAndGreatert(t *testing.T) {
	assert := assert.New(t)

	c1 := ParseClock("17:00:00")
	c2 := ParseClock("22:00:00")

	assert.True(c1.Less(c2))
	assert.False(c2.Less(c1))

	assert.True(c2.Greater(c1))
	assert.False(c1.Greater(c2))
}
