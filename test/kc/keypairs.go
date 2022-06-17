package kc

import (
	"sort"
)

type keyPairs struct {
	kpMap      map[string]int // map[column]status
	keys       []string       // 用户传递的 key
	structKeys []string       // 结构体中的 key
}

func (kp keyPairs) add() {
	for _, key := range kp.structKeys {
		kp.kpMap[key] = 1
	}
}

func (kp keyPairs) match() []string {
	for _, key := range kp.keys {
		kp.kpMap[key] |= 2
	}
	diffKeys := make([]string, 0, len(kp.kpMap))
	for key, status := range kp.kpMap {
		if status != 3 {
			diffKeys = append(diffKeys, key)
		}
	}
	sort.Strings(diffKeys)
	return diffKeys
}

// ExistsKeys 对比结构体对应 tag 字段和校验字段是否一致
func ExistsKeys(parser keyParser, option string, v interface{}, keys ...string) []string {
	kp := keyPairs{
		kpMap:      make(map[string]int, len(keys)),
		keys:       keys,
		structKeys: parser.Keys(v, option),
	}
	kp.add()
	return kp.match()
}
