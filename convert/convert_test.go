package convert

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructToURLValues(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	type User struct {
		Name      string       `json:"name"`
		Age       int          `json:"age"`
		Adult     bool         `json:"adult"`
		Chinese   *int         `json:"chinese"`
		Guardian  map[int]User `json:"guardian"`
		Friends   []User       `json:"friends"`
		Classmate []*User      `json:"classmate"`
	}

	form, err := StructToURLValues(&User{
		Name:      "张三",
		Age:       10,
		Adult:     false,
		Chinese:   NewInt(1),
		Guardian:  map[int]User{11: {Name: "11"}, 22: {Name: "22"}},
		Friends:   []User{{Name: "1"}, {Name: "2"}},
		Classmate: []*User{{Name: "111"}, {Name: "222"}},
	}, JSON)
	require.NoError(err)
	fmt.Println(form)
	assert.NotEmpty(form.Get("name"))
	assert.NotEmpty(form.Get("age"))
	assert.NotEmpty(form.Get("adult"))
	assert.NotEmpty(form.Get("chinese"))
	assert.NotEmpty(form.Get("friends"))
	assert.NotEmpty(form.Get("guardian"))

	form, err = StructToURLValues(User{Name: "张三", Age: 10, Adult: false}, JSON)
	require.NoError(err)
	assert.NotEmpty(form.Get("name"))
	assert.NotEmpty(form.Get("age"))
	assert.NotEmpty(form.Get("adult"))
	assert.NotEmpty(form.Get("friends"))
	assert.NotEmpty(form.Get("guardian"))

	form, err = StructToURLValues(1, JSON)
	require.EqualError(err, "param not is struct")
	assert.Nil(form)
}
