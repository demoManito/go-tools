package slice

import (
	"errors"
	"reflect"
)

type uniqFunc func(reflect.Value) []interface{}

func uniq(data interface{}, f uniqFunc) ([]interface{}, error) {
	v := reflect.ValueOf(data)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, errors.New("cannot resolve type")
	}
	if v.Len() == 0 {
		return []interface{}{}, nil
	}
	return f(v), nil
}

// Uniq deduplicate and return the deduplicated slice/array/struct
// params: <data> array to be deduplicated
// params: <args> tag and tag_name for structs
/*
	int:
	[1,1,2,2,3,3] => [1,2,3]

	string:
	["1","1","2","2","3","3"] => ["1","2","3"]

	struct:
	type test struct {
	  A int `json:"a"`
	  B int `json:"b"`
	}
	data=[{A:1, B:1}, {A:2, B:2}, {A:1, B:3}]
	Uniq(data, "json", "a") => [{A:1, B:1}, {A:1, B:3}]
*/
func Uniq(data interface{}, args ...string) ([]interface{}, error) {
	return uniq(data, func(v reflect.Value) []interface{} {
		uniqMap := make(map[interface{}]int)
		uniqArr := make([]interface{}, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			var field interface{}
			var fieldVal interface{}
			// WORKAROUND: when the element is interface{} && primitive data type, null pointer exception will be generated
			index := v.Index(i)
			for index.Kind() == reflect.Ptr || index.Kind() == reflect.Interface {
				index = index.Elem()
			}
			switch index.Kind() {
			case reflect.Struct:
				for j := 0; j < index.NumField(); j++ {
					if index.Type().Field(j).Tag.Get(args[0]) == args[1] {
						field = index.Field(j).Addr().Elem().Interface()
						fieldVal = v.Index(i).Addr().Elem().Interface()
						break
					}
					// WORKAROUND: handle when <tag> does not exist
					// this tag contains tag_name and tag, Example: [A int `json:"a"`] includes <json> and <a>
					field = i
					fieldVal = v.Index(i).Addr().Elem().Interface()
				}
			default:
				field = v.Index(i).Addr().Elem().Interface()
				fieldVal = field
			}
			if _, ok := uniqMap[field]; !ok {
				uniqMap[field] = i
				uniqArr = append(uniqArr, fieldVal)
			}
		}
		return uniqArr
	})
}

// UniqItem return duplicate elements
func UniqItem(data interface{}) ([]interface{}, error) {
	return uniq(data, func(v reflect.Value) []interface{} {
		uniqMap := make(map[interface{}]int)
		uniqArr := make([]interface{}, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Addr().Elem().Interface()
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = i
				continue
			}
			uniqArr = append(uniqArr, item)
		}
		uniq, _ := Uniq(uniqArr)
		return uniq
	})
}

// UniqIndex returns the index of the repeated data
func UniqIndex(data interface{}) ([]interface{}, error) {
	return uniq(data, func(v reflect.Value) []interface{} {
		uniqMap := make(map[interface{}]int)
		uniqArr := make([]interface{}, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Addr().Elem().Interface()
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = i
				continue
			}
			uniqArr = append(uniqArr, i)
		}
		return uniqArr
	})
}
