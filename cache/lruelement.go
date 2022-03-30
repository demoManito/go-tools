package cache

import "container/list"

var _ ILRU = new(LRUElement)

// LRUElement 缓存淘汰策略 - 按缓存的 KEY 数量做为淘汰容量
type LRUElement struct {
	LRU
}

func NewLRUElement(maxCap int64, onEvicted func(key string, value Value)) *LRUElement {
	return &LRUElement{
		LRU: LRU{
			OnEvicted:  onEvicted,
			maxCap:     maxCap,
			linkedList: list.New(),
			cacheMap:   make(map[string]*list.Element),
		},
	}
}

func (lru *LRUElement) Add(key string, value Value) {
	if ele, ok := lru.cacheMap[key]; ok {
		lru.linkedList.MoveToFront(ele)
		entry := ele.Value.(*entry)
		entry.value = value
	} else {
		lru.cacheMap[key] = lru.linkedList.PushFront(&entry{key: key, value: value})
	}
	for lru.maxCap != 0 && lru.maxCap < int64(len(lru.cacheMap)) {
		lru.RemoveOldest()
	}
}

func (lru *LRUElement) RemoveOldest() {
	ele := lru.linkedList.Back()
	if ele != nil {
		lru.linkedList.Remove(ele)
		entry := ele.Value.(*entry)
		delete(lru.cacheMap, entry.key)
		if lru.OnEvicted != nil {
			lru.OnEvicted(entry.key, entry.value)
		}
	}
}

// Get
func (lru *LRUElement) Get(key string) (Value, bool) {
	if ele, ok := lru.cacheMap[key]; ok {
		lru.linkedList.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return nil, false
}

func (lru *LRUElement) Cap() int64 {
	return int64(lru.Len())
}

func (lru *LRUElement) Len() int {
	return len(lru.cacheMap)
}
