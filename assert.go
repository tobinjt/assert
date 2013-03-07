/*
Are you fed up writing tests like this?

  import "reflect"
  import "testing"
  func TestSomething(t *testing.T) {
          result := something()
          expected := []int{2, 9, 6}
          if !reflect.DeepEqual(expected, result) {
  		t.Errorf("something(): %#v != %#v\n", expected, result)
          }
  }

Would you prefer to write tests like this?

  import "github.com/tobinjt/assert"
  import "testing"
  func TestSomething(t *testing.T) {
  	assert.Equal(t, "something()", []int{2, 9, 6}, something())
  }

This package makes it easy.

All functions return true if the test passes, and false if the test fails.
*/
package assert

import (
	"reflect"
	"testing"
)

// If reflect.DeepEqual(a, b) fails, call:
//     t.Errorf("%s: %#v != %#v\n", message, a, b)
func Equal(t *testing.T, message string, a, b interface{}) bool {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%s: %#v != %#v\n", message, a, b)
		return false
	}
	return true
}

// If err != nil, call:
//   t.Errorf("Unexpected error: %s: %s\n", message, err)
func ErrIsNil(t *testing.T, message string, err error) bool {
	if err != nil {
		t.Errorf("Unexpected error: %s: %s\n", message, err)
		return false
	}
	return true
}
