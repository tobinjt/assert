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

All functions return true if the test passes, and false if the test fails.  This
allows you to write tests like:
  func TestSomething(t *testing.T) {
	  result, err := something()
	  if assert.ErrIsNil(t, "something()", err) {
		  assert.Equal(t, "something()", 7, result)
	  }
  }
*/
package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

/*
An interface to enable writing tests for this package; you'll pass *testing.T.
*/
type T interface {
	Errorf(format string, args ...interface{})
}

// Returns "filename:line_number" for the caller's caller.
func getCallerSourceLocation() string {
	_, file, line, ok := runtime.Caller(2)
	result := "unknown:unknown"
	if ok {
		result = fmt.Sprintf("%s:%d", file, line)
	}
	return result
}

// If reflect.DeepEqual(a, b) fails, call:
//     t.Errorf("%s:%d: %s: %#v != %#v\n", file, line, message, a, b)
func Equal(t T, message string, a, b interface{}) bool {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%s: %s: %#v != %#v\n", getCallerSourceLocation(),
			message, a, b)
		return false
	}
	return true
}

// If err != nil, call:
//   t.Errorf("%s:%d: Unexpected error: %s: %s\n", file, line, message, err)
func ErrIsNil(t T, message string, err error) bool {
	if err != nil {
		t.Errorf("%s: Unexpected error: %s: %s\n",
			getCallerSourceLocation(), message, err)
		return false
	}
	return true
}

// If err == nil, call:
//   t.Errorf("%s:%d: Error is nil: %s: %s\n", file, line, message)
// If !strings.Contains(err.Error(), substr), call:
//   t.Errorf("%s:%d: Expected substring missing: %s\nsubstring: %s\nerror: %v\n",
//	file, line, message, substr, err)
func ErrContains(t T, message string, err error, substr string) bool {
	if err == nil {
		t.Errorf("%s: Error is nil: %s: %s\n",
			getCallerSourceLocation(), message)
		return false
	}
	if !strings.Contains(err.Error(), substr) {
		t.Errorf("%s: Expected substring missing: %s\nsubstring: %s\nerror: %v\n",
			getCallerSourceLocation(), message, substr, err)
		return false
	}
	return true
}
