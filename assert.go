/*Package assert makes it easier to write tests.

Are you fed up writing tests like this?

	func TestSomething(t *testing.T) {
		result := something()
		expected := []int{2, 9, 6}
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("something(): got %#v, want %#v\n", result, expected)
		}
	}

Would you prefer to write tests like this?

	func TestSomething(t *testing.T) {
		assert.Equal(t, "something()", []int{2, 9, 6}, something())
	}

This package supports the latter style of writing tests.

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
	"errors"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strings"
)

/*
T is an interface to enable writing tests for this package; you'll pass *testing.T in normal usage of this package.
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

var failedAssertionCounter int

// FailedAssertionCounter returns the counter of failed assertions.
func FailedAssertionCounter() int {
	return failedAssertionCounter
}

// ResetFailedAssertionCounter resets the counter of failed assertions.
func ResetFailedAssertionCounter() {
	failedAssertionCounter = 0
}

// Equal checks that a == b.
// If reflect.DeepEqual(a, b) fails, call:
//     t.Errorf("%s:%d: %s: got %#v, want %#v\n", file, line, message, b, a)
func Equal(t T, message string, a, b interface{}) bool {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%s: %s: got %#v, want %#v\n", getCallerSourceLocation(),
			message, b, a)
		failedAssertionCounter++
		return false
	}
	return true
}

// FloatsAreClose checks that the first precision digits after the decimal point
// of a and b are equal.  E.g. FloatsAreClose(t, "", 0.1234567, 0.1234569, 6)
// will return true
// BEWARE: floating point precision is inaccurate enough that
// FloatsAreClose(2.123, 2.124, 3) will return true because 0.1**3 ==
// 0.0010000000000000002 and 2.124-2.123 == 0.0009999999999998899.
func FloatsAreClose(t T, message string, a, b float64, precision int) bool {
	threshold := math.Pow(0.1, float64(precision))
	difference := math.Abs(a - b)
	if difference > threshold {
		t.Errorf("%s: %s: floats %v and %v differ by %v > threshold %v",
			getCallerSourceLocation(), message, a, b, difference, threshold)
		failedAssertionCounter++
		return false
	}
	return true
}

// ErrIsNil checks that err == nil.
// If err != nil, call:
//   t.Errorf("%s:%d: Unexpected error: %s: %s\n", file, line, message, err)
func ErrIsNil(t T, message string, err error) bool {
	if err != nil {
		t.Errorf("%s: Unexpected error: %s: %s\n",
			getCallerSourceLocation(), message, err)
		failedAssertionCounter++
		return false
	}
	return true
}

// ErrContains checks that an error contains the expected substring.
// If err == nil, call:
//   t.Errorf("%s:%d: Error is nil: %s: %s\n", file, line, message)
// If !strings.Contains(err.Error(), substr), call:
//   t.Errorf("%s:%d: Expected substring missing: %s\nsubstring: %s\nerror: %v\n",
//	file, line, message, substr, err)
func ErrContains(t T, message string, err error, substr string) bool {
	if err == nil {
		t.Errorf("%s: Error is nil: %s\n",
			getCallerSourceLocation(), message)
		failedAssertionCounter++
		return false
	}
	if !strings.Contains(err.Error(), substr) {
		t.Errorf("%s: Expected substring missing: %s\nsubstring: %s\n    error: %v\n",
			getCallerSourceLocation(), message, substr, err)
		failedAssertionCounter++
		return false
	}
	return true
}

// Panics checks that a function called panic() with an expected error message.
// The argument to panic() must be a string, and that string must contain the expected substring.
// Usage is more complex than the other functions:
//	func TestPanics(t *testing.T) {
//		func() {
//			// Do setup.
//			// ....
//			// assert.Panics() will be called when the enclosing anonymous function exits.
//			defer assert.Panics(t, "something() did not panic properly", "expected substring")
//			something(arg1, arg2)
//		}()
//
//		// Next test of something().
//	}
func Panics(t T, message string, substr string) {
	r := recover()
	if r == nil {
		t.Errorf("Nothing called panic(): %s\n", message)
		return
	}
	if str, ok := r.(string); ok {
		err := errors.New(str)
		ErrContains(t, message, err, substr)
	} else {
		t.Errorf("Return value from recover() wasn't a string: %v", message)
	}
}
