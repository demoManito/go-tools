package kc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type keyParser interface {
	Keys(v interface{}, option string) []string
	Name() string
}

// KeysChecker keys checker
type KeysChecker struct {
	t       *testing.T
	kp      keyParser
	assert  *assert.Assertions
	require *require.Assertions
}

// NewKeysChecker new keys checker
func NewKeysChecker(t *testing.T, kp keyParser) *KeysChecker {
	return &KeysChecker{
		t:       t,
		kp:      kp,
		assert:  assert.New(t),
		require: require.New(t),
	}
}

// Check assert check keys
func (kc *KeysChecker) Check(v interface{}, keys ...string) {
	kc.assert.Emptyf(ExistsKeys(kc.kp, "", v, keys...), "checker: %s, type: %s",
		kc.kp.Name(), reflect.TypeOf(v).Name())
}

// RequireCheck require check keys
func (kc *KeysChecker) RequireCheck(v interface{}, keys ...string) {
	kc.require.Emptyf(ExistsKeys(kc.kp, "", v, keys...), "checker: %s, type: %s",
		kc.kp.Name(), reflect.TypeOf(v).Name())
}

// CheckOmitEmpty check omitEmpty
func (kc *KeysChecker) CheckOmitEmpty(v interface{}, keys ...string) {
	kc.assert.Emptyf(ExistsKeys(kc.kp, "omitempty", v, keys...), "checker: %s, type: %s",
		kc.kp.Name(), reflect.TypeOf(v).Name())
}

// RequireCheckOmitEmpty require check omitEmpty
func (kc *KeysChecker) RequireCheckOmitEmpty(v interface{}, keys ...string) {
	kc.require.Emptyf(ExistsKeys(kc.kp, "omitempty", v, keys...), "checker: %s, type: %s",
		kc.kp.Name(), reflect.TypeOf(v).Name())
}
