package lru

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLRUElement_Add(t *testing.T) {
	require := require.New(t)

	lru := NewLRUElement(3, func(key string, value Value) {
		t.Logf("key: %s; value: %s", key, value)
	})
	lru.Add("1", &testLRU{val: 5})
	lru.Add("2", &testLRU{val: 4})
	lru.Add("3", &testLRU{val: 1})
	lru.Add("4", &testLRU{val: 2}) // 这个会在最上方
	require.EqualValues(3, lru.Len())
	require.EqualValues(3, lru.Cap())
}

func TestLRUElement_RemoveOldest(t *testing.T) {
	require := require.New(t)

	lru := NewLRUElement(13, func(key string, value Value) {
		t.Logf("key: %s; value: %s", key, value)
	})
	lru.Add("1", &testLRU{val: 5})
	lru.Add("2", &testLRU{val: 4})
	lru.Add("3", &testLRU{val: 1})
	require.EqualValues(3, lru.Len())
	require.EqualValues(3, lru.Cap())

	lru.RemoveOldest()
	require.EqualValues(2, lru.Len())
	require.EqualValues(2, lru.Cap())
}

func TestLRUElement_Get(t *testing.T) {
	require := require.New(t)

	lru := NewLRUElement(13, func(key string, value Value) {
		t.Logf("key: %s; value: %s", key, value)
	})
	lru.Add("1", &testLRU{val: 5})
	lru.Add("2", &testLRU{val: 4})
	lru.Add("3", &testLRU{val: 1})

	val, ok := lru.Get("1")
	require.True(ok)
	require.EqualValues(5, val.Len())

	// 此时 key=1 被拉倒了链表最前方，删除最老的元素时，应该删除的是 key=2
	lru.RemoveOldest()
	require.EqualValues(2, lru.Len())
	require.EqualValues(2, lru.Cap())
}
