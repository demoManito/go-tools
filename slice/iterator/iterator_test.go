package iterator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_NewIteratorAndHashNextAndNext(t *testing.T) {
	require := require.New(t)

	is := []int{1, 2, 3, 4}
	iterator, err := NewIterator(is)
	require.NoError(err)
	for iterator.HasNext() {
		t.Log(iterator.Next())
	}

	ss := "111"
	iterator, err = NewIterator(ss)
	require.EqualError(err, "data is not slice or array")
	require.Nil(iterator)
}
