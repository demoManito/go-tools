package iterator

import (
	"errors"
	"reflect"
)

type IIterator interface {
	HasNext() bool
	Next() interface{}
}

// Iterator 迭代器
type Iterator struct {
	d     reflect.Value
	index int
}

// New
func New(data interface{}) (IIterator, error) {
	d := reflect.ValueOf(data)
	for d.Kind() == reflect.Interface || d.Kind() == reflect.Ptr {
		d = d.Elem()
	}
	if d.Kind() != reflect.Slice && d.Kind() != reflect.Array {
		return nil, errors.New("data is not slice or array")
	}
	iterator := &Iterator{
		d:     d,
		index: 0,
	}
	return iterator, nil
}

// HasNext next is there exist item
func (i *Iterator) HasNext() bool {
	return i.index < i.d.Len()
}

// Next get now item && index incr
func (i *Iterator) Next() interface{} {
	if i.HasNext() {
		data := i.d.Index(i.index)
		i.index++
		return data
	}
	return nil
}
