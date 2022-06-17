package kc

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/jinzhu/gorm"
)

// currently supported Tags
var (
	BSON = TagParser{TagName: "bson"}
	GORM = TagParser{TagName: "gorm"}
	JSON = TagParser{TagName: "json"}
	YAML = TagParser{TagName: "yaml"}
	FORM = TagParser{TagName: "form"}
)

type TagParser struct {
	TagName string
}

func (tp TagParser) Keys(v interface{}, option string) []string {
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	if rt.Kind() == reflect.Ptr {
		rt = reflect.TypeOf(reflect.Indirect(rv).Interface())
	}
	pre := make([]string, 0, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		// 非导出字段不做校验
		if unicode.IsLower(rune(field.Name[0])) {
			continue
		}

		tag := field.Tag.Get(tp.TagName)
		if tag != "" && tag != "-" {
			pre = append(pre, tag)
		}
		if tag == "" && !field.Anonymous {
			// 没有 tag 时, gorm 按照默认规则导出, json、bson、yaml 等其他类型直接导出结构体字段名
			if tp.TagName == GORM.TagName {
				pre = append(pre, fmt.Sprintf("colume:%s", gorm.ToColumnName(field.Name)))
			} else {
				pre = append(pre, field.Name)
			}
		}
	}

	res := make([]string, 0, len(pre))
	appendNotEmpty := func(key string) {
		if key != "" {
			res = append(res, key)
		}
	}

	switch tp.TagName {
	case GORM.TagName:
		tags := make([]map[string]string, len(pre))
		for i, p := range pre {
			parts := strings.Split(p, ";") // example: `gorm:"column:id;primary_key"`
			m := make(map[string]string, len(parts))
			for _, part := range parts {
				gormSplitTag := strings.Split(part+":", ":")
				m[gormSplitTag[0]] = gormSplitTag[1] // example: {"column":"id"} or {"primary_key": ""}
			}
			tags[i] = m // example: [{"column":"id", "primary_key": ""}, {"column": "name"}]
		}
		for _, tag := range tags {
			if t, ok := tag["column"]; ok {
				res = append(res, t)
			}
		}
		return res
	default: // bson, form, json, yaml...
		if option == "omitempty" {
			for i := range pre {
				if strings.Contains(pre[i], ",omitempty") { // 防止 tag 就是 omitempty
					appendNotEmpty(strings.Split(pre[i], ",")[0])
				}
			}
		} else {
			for i := range pre {
				appendNotEmpty(strings.Split(pre[i], ",")[0])
			}
		}
		return res
	}
}

func (tp TagParser) Name() string {
	return tp.TagName
}
