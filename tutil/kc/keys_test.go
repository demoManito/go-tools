package kc

import "testing"

func TestKeysChecker_Check(t *testing.T) {
	kc := NewKeysChecker(t, GORM)
	kc.Check(testKeys{}, "id", "name", "age")
	kc.RequireCheck(testKeys{}, "id", "name", "age")
}
