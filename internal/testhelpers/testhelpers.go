package testhelpers

import (
	"reflect"
	"strings"
	"testing"
)

// Equals compares two values and fails if they are not equal.
func Equals(t *testing.T, expected, value interface{}) {
	t.Helper()
	if reflect.DeepEqual(expected, value) {
		return
	}

	t.Fatalf("Expected `%v`, got `%v`", expected, value)
}

// Equals compares two values and fails if they are equal.
func NotEquals(t *testing.T, expected, value interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, value) {
		return
	}

	t.Fatalf("Expected `%v` NOT to be `%v`", expected, value)
}

// NotNil checks if a value is not nil.
func NotNil(t *testing.T, value interface{}) {
	t.Helper()
	if value != nil {
		return
	}

	t.Fatalf("Expected `%v` NOT to be nil", value)
}

// Nil checks if a value is nil.
func Nil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil {
		return
	}

	t.Fatalf("Expected `%v` to be nil", value)
}

// True checks if a value is true.
func True(t *testing.T, val bool) {
	t.Helper()
	if val {
		return
	}

	t.Fatalf("Expected `%v` to be true", val)
}

// False checks if a value is false.
func False(t *testing.T, val bool) {
	t.Helper()
	if val {
		return
	}

	t.Fatalf("Expected `%v` to be false", val)
}

// Contains checks if a string contains a substring.
func Contains(t *testing.T, content, term string) {
	t.Helper()
	if strings.Contains(content, term) {
		return
	}

	t.Fatalf("Expected `%v` to contain `%v`", content, term)
}

// NotContains checks if a string does not contain a substring.
func NotContains(t *testing.T, content, term string) {
	t.Helper()
	if !strings.Contains(content, term) {
		return
	}

	t.Fatalf("Expected `%v` NOT to contain `%v`", content, term)
}

// Error checks if the error passed is not nil
func Error(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		return
	}

	t.Fatalf("Expected error, got nil")
}

// NoError checks if the error passed is nil
func NoError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}

	t.Fatalf("Expected no error, got `%v`", err)
}
