package kc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExistsKeys(t *testing.T) {
	assert := assert.New(t)

	testk := new(testKeys)

	diff := ExistsKeys(GORM, "", testk, "id", "name", "age")
	assert.Empty(diff)

	diff = ExistsKeys(GORM, "", testk, "id", "name")
	assert.Len(diff, 1)
	assert.Equal([]string{"age"}, diff)

	diff = ExistsKeys(GORM, "", testk, "id")
	assert.Len(diff, 2)
	assert.Equal([]string{"age", "name"}, diff)

	diff = ExistsKeys(FORM, "", testk, "id", "name", "age")
	assert.Len(diff, 3)
	assert.Equal([]string{"age", "id", "name"}, diff)

	diff = ExistsKeys(YAML, "", testk, "id", "name", "age")
	assert.Len(diff, 6)
	assert.Equal([]string{"Age", "ID", "Name", "age", "id", "name"}, diff)
}
