package cache

import "container/list"

type ILRU interface {
	Add(string, Value)
	RemoveOldest()
	Get(string) (Value, bool)
	Cap() int64
	Len() int
}

var _ ILRU = new(LRU)

// LRU 缓存淘汰策略 - 按值字节数做为淘汰容量
// hashmap + linklist
type LRU struct {
	OnEvicted func(key string, value Value) // 移除元素时触发时间

	maxCap     int64
	nowCap     int64
	cacheMap   map[string]*list.Element
	linkedList *list.List
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func NewLRU(maxCap int64, onEvicted func(key string, value Value)) *LRU {
	return &LRU{
		OnEvicted:  onEvicted,
		maxCap:     maxCap,
		cacheMap:   make(map[string]*list.Element),
		linkedList: list.New(),
	}
}

func (lru *LRU) Add(key string, value Value) {
	if ele, ok := lru.cacheMap[key]; ok {
		lru.linkedList.MoveToFront(ele)
		entry := ele.Value.(*entry)
		lru.nowCap += int64(value.Len()) - int64(entry.value.Len())
		entry.value = value
	} else {
		lru.cacheMap[key] = lru.linkedList.PushFront(&entry{key: key, value: value})
		lru.nowCap += int64(len(key) + value.Len())
	}
	for lru.maxCap != 0 && lru.maxCap < lru.nowCap {
		lru.RemoveOldest()
	}
}

func (lru *LRU) RemoveOldest() {
	ele := lru.linkedList.Back()
	if ele != nil {
		lru.linkedList.Remove(ele)
		entry := ele.Value.(*entry)
		delete(lru.cacheMap, entry.key)
		lru.nowCap -= int64(len(entry.key) + entry.value.Len())
		if lru.OnEvicted != nil {
			lru.OnEvicted(entry.key, entry.value)
		}
	}
}

func (lru *LRU) Get(key string) (Value, bool) {
	if ele, ok := lru.cacheMap[key]; ok {
		lru.linkedList.MoveToFront(ele)
		entry := ele.Value.(*entry)
		return entry.value, true
	}
	return nil, false
}

func (lru *LRU) Cap() int64 {
	return lru.nowCap
}

func (lru *LRU) Len() int {
	return len(lru.cacheMap)
}
