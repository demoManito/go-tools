package lru

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testLRU struct {
	val int
}

func (t *testLRU) Len() int {
	return t.val
}

func TestLRU_Add(t *testing.T) {
	require := require.New(t)

	lru := NewLRU(10, func(key string, value Value) {
		t.Logf("key: %s; value: %s", key, value)
	})
	lru.Add("1", &testLRU{val: 5})
	lru.Add("2", &testLRU{val: 4})
	lru.Add("3", &testLRU{val: 1})
	lru.Add("4", &testLRU{val: 2}) // 这个会在最上方
	require.EqualValues(3, lru.Len())
}

func TestLRU_RemoveOldest(t *testing.T) {
	require := require.New(t)

	lru := NewLRU(13, func(key string, value Value) {
		t.Logf("key: %s; value: %s", key, value)
	})
	lru.Add("1", &testLRU{val: 5})
	lru.Add("2", &testLRU{val: 4})
	lru.Add("3", &testLRU{val: 1})
	require.EqualValues(3, lru.Len())
	require.EqualValues(13, lru.Cap()) // 每个 KEY 占用 1 个长度

	lru.RemoveOldest()
	require.EqualValues(2, lru.Len())
	require.EqualValues(13-5-1, lru.Cap())
}

func TestLRU_Get(t *testing.T) {
	require := require.New(t)

	lru := NewLRU(13, func(key string, value Value) {
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
	require.EqualValues(13-4-1, lru.Cap())
}
