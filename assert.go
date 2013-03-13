/*
Are you fed up writing tests like this?

  func TestSomething(t *testing.T) {
          result := something()
          expected := []int{2, 9, 6}
          if !reflect.DeepEqual(expected, result) {
  		t.Errorf("something(): %#v != %#v\n", expected, result)
          }
  }

Would you prefer to write tests like this?

  func TestSomething(t *testing.T) {
  	assert.Equal(t, "something()", []int{2, 9, 6}, something())
  }

This package makes it easy.

All functions return true if the test passes, and false if the test fails.
*/
package assert

import (
	"reflect"
	"strings"
)

/*
An interface to enable writing tests for this package; you'll pass *testing.T.
*/
type T interface {
	Errorf(format string, args ...interface{})
}

// If reflect.DeepEqual(a, b) fails, call:
//     t.Errorf("%s: %#v != %#v\n", message, a, b)
func Equal(t T, message string, a, b interface{}) bool {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%s: %#v != %#v\n", message, a, b)
		return false
	}
	return true
}

// If err != nil, call:
//   t.Errorf("Unexpected error: %s: %s\n", message, err)
func ErrIsNil(t T, message string, err error) bool {
	if err != nil {
		t.Errorf("Unexpected error: %s: %s\n", message, err)
		return false
	}
	return true
}

// If err == nil, call:
//   t.Errorf("Error is nil: %s: %s\n", message)
// If !strings.Contains(err.Error(), substr), call:
//   t.Errorf("Expected substring missing: %s\nsubstring: %s\nerror: %v\n",
//	message, substr, err)
func ErrContains(t T, message string, err error, substr string) bool {
	if err == nil {
		t.Errorf("Error is nil: %s: %s\n", message)
		return false
	}
	if !strings.Contains(err.Error(), substr) {
		t.Errorf("Expected substring missing: %s\nsubstring: %s\nerror: %v\n",
			message, substr, err)
		return false
	}
	return true
}
