package testhelpers

import (
	"reflect"
	"strings"
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

func Contains(t *testing.T, content, term string) {
	t.Helper()
	if strings.Contains(content, term) {
		return
	}

	t.Fatalf("Expected `%v` to contain `%v`", content, term)
}

func NotContains(t *testing.T, content, term string) {
	t.Helper()
	if !strings.Contains(content, term) {
		return
	}

	t.Fatalf("Expected `%v` NOT to contain `%v`", content, term)
}

func Error(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		return
	}

	t.Fatalf("Expected error, got nil")
}

func NoError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}

	t.Fatalf("Expected no error, got `%v`", err)
}
