package assert

import (
	"errors"
	"testing"
)

type testlogger struct {
	message string
	counter int
}

func NewTestLogger() *testlogger {
	return &testlogger{}
}

func (t *testlogger) Errorf(message string, _ ...interface{}) {
	t.counter++
	t.message = message
}

func TestEqual(t *testing.T) {
	logger := NewTestLogger()
	Equal(logger, "message", 7, 23)
	if 1 != logger.counter {
		t.Errorf("Equal: Errorf not called for inputs 7, 23\n")
	}
	logger = NewTestLogger()
	Equal(logger, "message", 42, 42)
	if 0 != logger.counter {
		t.Errorf("Equal: Errorf unexpectedly called for inputs 42, 42\n")
	}
}

func TestErrIsNil(t *testing.T) {
	logger := NewTestLogger()
	ErrIsNil(logger, "message", errors.New("A dummy error"))
	if 1 != logger.counter {
		t.Error("ErrIsNil: Errorf not called")
	}
	logger = NewTestLogger()
	ErrIsNil(logger, "message", nil)
	if 0 != logger.counter {
		t.Error("ErrIsNil: Errorf unexpectedly called")
	}
}

func TestErrContains(t *testing.T) {
	logger := NewTestLogger()
	ErrContains(logger, "message", nil, "")
	if 1 != logger.counter {
		t.Error("ErrContains: Errorf not called for nil error")
	}
	logger = NewTestLogger()
	ErrContains(logger, "message", errors.New("asdf"), "qwerty")
	if 1 != logger.counter {
		t.Error("ErrContains: Errorf not called for bad error message")
	}
	logger = NewTestLogger()
	ErrContains(logger, "message", errors.New("asdf"), "asd")
	if 0 != logger.counter {
		t.Error("ErrContains: Errorf unexpectedly called")
	}
}
