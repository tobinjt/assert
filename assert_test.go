package assert

import (
	"errors"
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

func (t *testlogger) Errorf(message string, _ ...interface{}) {
	t.counter++
	t.message = message
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
		t.Errorf("Equal: Errorf unexpectedly called for inputs 42, 42\n")
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
		t.Error("ErrIsNil: Errorf unexpectedly called")
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
		t.Error("ErrContains: Errorf unexpectedly called")
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
		t.Error("Panics: Errorf was called even though panic was called")
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
	// ErrContains calls t.Errorf, but the Errorf we use in testing doesn't do formatting, so it won't contain the message passed to Panics; instead we check for the fixed string from ErrContains.
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
	if logger.message != "%s: recover() returned nil: %s\n" {
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
