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
*/
package assert

import (
	"reflect"
	"testing"
)

// Uses reflect.DeepEqual() to compare a and b, and calls
//     t.Errorf("%s: %#v != %#v\n", message, a, b)
// if the comparison fails.
func Equal(t *testing.T, message string, a, b interface{}) bool {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%s: %#v != %#v\n", message, a, b)
		return false
	}
	return true
}
