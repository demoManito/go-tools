package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniq(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// slice
	arr, err := Uniq([]int{1, 2, 3, 4, 4, 4, 5})
	require.NoError(err)
	assert.NotNil(arr)
	assert.Equal([]interface{}{1, 2, 3, 4, 5}, arr)

	arr, err = Uniq([]int{1, 2, 3, 4, 5})
	require.NoError(err)
	assert.NotNil(arr)
	assert.Equal([]interface{}{1, 2, 3, 4, 5}, arr)

	arr, err = Uniq([]interface{}{})
	require.NoError(err)
	assert.Empty(arr)

	_, err = Uniq(1)
	require.EqualError(err, "cannot resolve type")

	// string
	arr, err = Uniq([]string{"1", "1", "2", "2", "3", "3"})
	require.NoError(err)
	assert.NotNil(arr)
	assert.Equal([]interface{}{"1", "2", "3"}, arr)

	// interface
	arr, err = Uniq([]interface{}{"1", "1", "2", "2", "3", "3"})
	require.NoError(err)
	assert.NotNil(arr)
	assert.Equal([]interface{}{"1", "2", "3"}, arr)

	// struct
	type test struct {
		A int `json:"a"`
		B int `json:"b"`
	}
	arr, err = Uniq([]test{{A: 1, B: 2}, {A: 1, B: 3}, {A: 2, B: 2}}, "json", "a")
	require.NoError(err)
	assert.Len(arr, 2)

	arr, err = Uniq([]*test{{A: 1, B: 2}, {A: 1, B: 3}, {A: 2, B: 2}}, "json", "a")
	require.NoError(err)
	assert.Len(arr, 2)

	arr, err = Uniq([]*test{{A: 1, B: 2}, {A: 1, B: 3}, {A: 2, B: 2}}, "json", "b")
	require.NoError(err)
	assert.Len(arr, 2)

	arr, err = Uniq([]*test{{A: 1, B: 2}, {A: 1, B: 3}, {A: 2, B: 2}}, "json", "c")
	require.NoError(err)
	require.Len(arr, 3)

	arr, err = Uniq([]*test{{A: 1, B: 2}, {A: 1, B: 3}, {A: 2, B: 2}}, "yaml", "a")
	require.NoError(err)
	require.Len(arr, 3)
}

func TestUniqItem(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	arr, err := UniqItem([]int{1, 2, 3, 4, 4, 4, 5})
	require.NoError(err)
	assert.NotNil(arr)
	assert.Equal([]interface{}{4}, arr)

	arr, err = UniqItem([]int{1, 2, 3, 4, 5})
	require.NoError(err)
	assert.NotNil(arr)
	assert.Equal([]interface{}{}, arr)

	arr, err = UniqItem([]interface{}{})
	require.NoError(err)
	assert.Empty(arr)

	_, err = UniqItem(1)
	require.EqualError(err, "cannot resolve type")
}

func TestUniqIndex(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	ms, err := UniqIndex([]int{1, 2, 3, 4, 4, 4, 5})
	require.NoError(err)
	assert.NotEmpty(ms)
	assert.Equal([]interface{}{4, 5}, ms)

	ms, err = UniqIndex([]int{1, 2, 3, 4, 5})
	require.NoError(err)
	assert.Empty(ms)
	assert.Equal([]interface{}{}, ms)

	ms, err = UniqIndex([]interface{}{})
	require.NoError(err)
	assert.Empty(ms)

	_, err = UniqIndex(1)
	require.EqualError(err, "cannot resolve type")
}
