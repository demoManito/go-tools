package convert

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

// BoolToInt bool 转 int
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// IntToBool int 转 bool
func IntToBool(i int) bool {
	return i != 0
}

// NewBool bool 转指针类型
func NewBool(b bool) *bool {
	return &b
}

// NewInt creates new int pointer
func NewInt(i int) *int {
	return &i
}

// NewInt64 creates new int64 pointer
func NewInt64(i int64) *int64 {
	return &i
}

// NewFloat64 creates new float64 pointer
func NewFloat64(i float64) *float64 {
	return &i
}

const (
	JSON = "json"
	LOAD = "load"
	YAML = "yaml"
	FORM = "form"
	BSON = "bson"
)

// StructToURLValues Struct => url.Values
// support: string、int、float、bool、slice
// because the url.Values type does not need too many parameters to modify, it does not need to support all types
func StructToURLValues(value interface{}, tagName string) (url.Values, error) {
	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("param not is struct")
	}

	form := url.Values{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		name := v.Type().Field(i).Tag.Get(tagName)
		for field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface {
			field = field.Elem()
		}
		var fieldVal string
		switch field.Kind() {
		case reflect.String:
			fieldVal = field.String()
		case reflect.Int:
			fieldVal = strconv.Itoa(int(field.Int()))
		case reflect.Float64:
			fieldVal = strconv.FormatFloat(field.Float(), 'f', 2, 64)
		case reflect.Bool:
			fieldVal = strconv.FormatBool(field.Bool())
		case reflect.Slice, reflect.Map:
			fieldVal = fmt.Sprint(field.Interface())
		}
		form.Add(name, fieldVal)
	}
	return form, nil
}

// SliceToMap slice to map
func SliceToMap() {

}
