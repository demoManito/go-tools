package kc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testKeys struct {
	ID   int64  `gorm:"column:id;primary_key" json:"-" bson:",omitempty" form:"-"`
	Name string `gorm:"column:name" json:"name" bson:"name,omitempty" form:"-"`
	Age  string `gorm:"column:age" json:"age" bson:"age" form:"-"`
}

func TestTagParser_Keys(t *testing.T) {
	assert := assert.New(t)

	tp := TagParser{TagName: GORM.TagName}
	keys := tp.Keys(new(testKeys), "")
	assert.Len(keys, 3)

	tp = TagParser{TagName: JSON.TagName}
	keys = tp.Keys(new(testKeys), "")
	assert.Len(keys, 2)

	tp = TagParser{TagName: BSON.TagName}
	keys = tp.Keys(new(testKeys), "")
	assert.Len(keys, 2)
	keys = tp.Keys(new(testKeys), "omitempty")
	assert.Len(keys, 1)
}

func TestTagParser_Name(t *testing.T) {
	assert := assert.New(t)

	tp := TagParser{TagName: "test"}
	assert.Equal(tp.TagName, tp.Name())
}
