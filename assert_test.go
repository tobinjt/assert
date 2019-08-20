package assert

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type testlogger struct {
	message string
	counter int
}

func newTestLogger() *testlogger {
	return &testlogger{}
}

func (t *testlogger) Errorf(message string, args ...interface{}) {
	t.counter++
	t.message = fmt.Sprintf(message, args...)
}

func TestErrorf(t *testing.T) {
	logger := newTestLogger()
	f := func() {
		logger.Errorf("message")
	}
	f()
	f()
	f()
	if 3 != logger.counter {
		t.Errorf("Logger: counter: got %v, want %v", logger.counter, 3)
	}
	if "message" != logger.message {
		t.Errorf("Logger: message: got %v, want message", logger.message)
	}
}

func TestFailedAssertionCounter(t *testing.T) {
	if 0 != failedAssertionCounter {
		t.Errorf("failedAssertionCounter: want 0, got %v", failedAssertionCounter)
	}
	logger := newTestLogger()
	Equal(logger, "message", 1, 2)
	Equal(logger, "message", 1, 2)
	Equal(logger, "message", 1, 2)
	if 3 != failedAssertionCounter {
		t.Errorf("failedAssertionCounter: want 3, got %v", failedAssertionCounter)
	}
	if 3 != FailedAssertionCounter() {
		t.Errorf("FailedAssertionCounter: want 3, got %v", FailedAssertionCounter())
	}
	ResetFailedAssertionCounter()
	if 0 != failedAssertionCounter {
		t.Errorf("failedAssertionCounter: want 0, got %v", failedAssertionCounter)
	}
	if 0 != FailedAssertionCounter() {
		t.Errorf("FailedAssertionCounter: want 0, got %v", FailedAssertionCounter())
	}
}

func TestEqual(t *testing.T) {
	logger := newTestLogger()
	Equal(logger, "message", 7, 23)
	if 1 != logger.counter {
		t.Errorf("Equal: Errorf not called for inputs 7, 23\n")
	}
	logger = newTestLogger()
	Equal(logger, "message", 42, 42)
	if 0 != logger.counter {
		t.Errorf("Equal: Errorf unexpectedly called for inputs 42, 42: %s\n", logger.message)
	}
}

func TestFloatsAreClose(t *testing.T) {
	logger := newTestLogger()
	FloatsAreClose(logger, "message", 2.123, 2.124, 2)
	if 0 != logger.counter {
		t.Errorf("FloatsAreClose: Errorf unexpectedly called for inputs 2.123, 2.124, 2: %s\n", logger.message)
	}
	logger = newTestLogger()
	FloatsAreClose(logger, "message", 2.123, 2.125, 3)
	if 1 != logger.counter {
		t.Errorf("FloatsAreClose: Errorf not called for inputs 2.123, 2.125, 3; message: %s\n", logger.message)
	}
}

func TestErrIsNil(t *testing.T) {
	logger := newTestLogger()
	ErrIsNil(logger, "message", errors.New("A dummy error"))
	if 1 != logger.counter {
		t.Error("ErrIsNil: Errorf not called")
	}
	logger = newTestLogger()
	ErrIsNil(logger, "message", nil)
	if 0 != logger.counter {
		t.Errorf("ErrIsNil: Errorf unexpectedly called: %s", logger.message)
	}
}

func TestErrContains(t *testing.T) {
	logger := newTestLogger()
	ErrContains(logger, "message", nil, "")
	if 1 != logger.counter {
		t.Error("ErrContains: Errorf not called for nil error")
	}
	logger = newTestLogger()
	ErrContains(logger, "message", errors.New("asdf"), "qwerty")
	if 1 != logger.counter {
		t.Error("ErrContains: Errorf not called for bad error message")
	}
	logger = newTestLogger()
	ErrContains(logger, "message", errors.New("asdf"), "asd")
	if 0 != logger.counter {
		t.Errorf("ErrContains: Errorf unexpectedly called: %s", logger.message)
	}
}

func TestPanics(t *testing.T) {
	// Panics as expected.
	logger := newTestLogger()
	func() {
		defer Panics(logger, "Function calls panic", "panicked")
		panic("I panicked :(")
	}()
	if 0 != logger.counter {
		t.Errorf("Panics: Errorf was called even though panic was called: %s", logger.message)
	}

	// Panics with bad error message.
	logger = newTestLogger()
	func() {
		defer Panics(logger, "Function calls panic with unexpected message", "panicked")
		panic("Unexpected message")
	}()
	if 1 != logger.counter {
		t.Error("Panics: Errorf was not called even though panic was called")
	}
	if !strings.Contains(logger.message, "Expected substring missing") {
		t.Errorf("Panics: bad message: got %q, want %q", logger.message, "Expected substring missing")
	}

	// Panic is not called.
	logger = newTestLogger()
	func() {
		defer Panics(logger, "Panic is not called", "panicked")
	}()
	if 1 != logger.counter {
		t.Error("Panics: Errorf was not called even though panic was not called")
	}
	if !strings.Contains(logger.message, "Nothing called panic():") {
		t.Errorf("Panics: bad message: got %q, want %q", logger.message, "%s: recover() returned nil: %s\n")
	}

	logger = newTestLogger()
	func() {
		defer Panics(logger, "", "panicked")
		panic([]int64{7})
	}()
	if 1 != logger.counter {
		t.Error("Panics: Errorf was not called even though argument to panic() was not a string")
	}
	if !strings.Contains(logger.message, "Return value from recover() wasn't a string") {
		t.Errorf("Panics: bad message: got %q, want %q", logger.message, "Return value from recover() wasn't a string")
	}
}
