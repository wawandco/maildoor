package testhelpers

import (
	"reflect"
	"testing"
)

func Equals(t *testing.T, expected, value interface{}) {
	t.Helper()
	if reflect.DeepEqual(expected, value) {
		return
	}

	t.Fatalf("Expected `%v`, got `%v`", expected, value)
}

func NotEquals(t *testing.T, expected, value interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, value) {
		return
	}

	t.Fatalf("Expected `%v` NOT to be `%v`", expected, value)
}

func NotNil(t *testing.T, value interface{}) {
	t.Helper()
	if value != nil {
		return
	}

	t.Fatalf("Expected `%v` NOT to be nil", value)
}

func Nil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil {
		return
	}

	t.Fatalf("Expected `%v` to be nil", value)
}

func True(t *testing.T, val bool) {
	t.Helper()
	if val {
		return
	}

	t.Fatalf("Expected `%v` to be true", val)
}

func False(t *testing.T, val bool) {
	t.Helper()
	if val {
		return
	}

	t.Fatalf("Expected `%v` to be false", val)
}
