package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInclude(t *testing.T) {
	assert := assert.New(t)

	is := []int{2, 2, 2, 2, 2, 21, 11, 1}
	assert.True(Include(len(is), func(i int) bool { return is[i] == 21 }))
	assert.False(Include(len(is), func(i int) bool { return is[i] == 22 }))

	ss := []string{"1", "2", "1", "2", "1", "2", "1", "2", "3"}
	assert.True(Include(len(ss), func(i int) bool { return ss[i] == "3" }))
	assert.False(Include(len(ss), func(i int) bool { return ss[i] == "4" }))

	// ... other type do not test
}

func TestFindIndex(t *testing.T) {
	assert := assert.New(t)

	list := []int64{1, 2, 3, 4, 5, 6, 1, 2}
	assert.Equal(1, FindIndex(len(list), func(i int) bool {
		return list[i] == 2
	}), "test return top element")

	assert.Equal(-1, FindIndex(len(list), func(i int) bool {
		return list[i] == 7
	}), "test not found return -1")

	assert.Equal(-1, FindIndex(3, func(i int) bool {
		return list[i] == 6
	}), "test find not max length")

	list = []int64{}
	assert.Equal(-1, FindIndex(len(list), func(i int) bool {
		return list[i] == 2
	}))
}
