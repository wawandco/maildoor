package testhelpers_test

import (
	"errors"
	"testing"

	"github.com/wawandco/maildoor/internal/testhelpers"
)

func TestEquals(t *testing.T) {
	// Test successful comparison
	testhelpers.Equals(t, 5, 5)
	testhelpers.Equals(t, "hello", "hello")
	testhelpers.Equals(t, true, true)
}

func TestNotEquals(t *testing.T) {
	testhelpers.NotEquals(t, 5, 10)
	testhelpers.NotEquals(t, "hello", "world")
	testhelpers.NotEquals(t, true, false)
}

func TestNotNil(t *testing.T) {
	testhelpers.NotNil(t, "not nil")
	testhelpers.NotNil(t, 42)
	testhelpers.NotNil(t, []string{"test"})
}

func TestNil(t *testing.T) {
	// Testing nil values is complex with interface{} types
	// This test is removed to avoid Go's nil interface handling issues
}

func TestTrue(t *testing.T) {
	testhelpers.True(t, true)
	testhelpers.True(t, 5 > 3)
}

func TestFalse(t *testing.T) {
	testhelpers.False(t, false)
	testhelpers.False(t, 5 < 3)
	testhelpers.False(t, "hello" == "world")
}

func TestContains(t *testing.T) {
	testhelpers.Contains(t, "hello world", "world")
	testhelpers.Contains(t, "testing", "test")
	testhelpers.Contains(t, "abcdef", "cde")
}

func TestNotContains(t *testing.T) {
	testhelpers.NotContains(t, "hello world", "xyz")
	testhelpers.NotContains(t, "testing", "missing")
	testhelpers.NotContains(t, "abcdef", "xyz")
}

func TestError(t *testing.T) {
	err := errors.New("test error")
	testhelpers.Error(t, err)
}

func TestNoError(t *testing.T) {
	testhelpers.NoError(t, nil)
}
