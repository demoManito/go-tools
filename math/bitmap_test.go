package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitMap(t *testing.T) {
	assert := assert.New(t)

	bm8 := struct {
		sets BitMap8
	}{
		sets: BitMap8(0),
	}
	bm8.sets.Set(8)
	assert.True(bm8.sets.IsSet(8))
	bm8.sets.UnSet(8)
	assert.False(bm8.sets.IsSet(8))

	bm16 := struct {
		sets BitMap16
	}{
		sets: BitMap16(0),
	}
	bm16.sets.Set(16)
	assert.True(bm16.sets.IsSet(16))
	bm16.sets.UnSet(16)
	assert.False(bm16.sets.IsSet(16))

	// ... other type do not test
}
