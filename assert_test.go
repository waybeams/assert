package assert_test

import (
	"fmt"
	"github.com/waybeams/assert"
	"reflect"
	"testing"
)

type fakeConst int

const (
	fakeConstValueA = iota
	fakeConstValueB
)

type CustomT struct {
	testing.T
	failureMsg string
}

func (c *CustomT) Errorf(format string, args ...interface{}) {
	c.failureMsg = fmt.Sprintf(format, args...)
}

func (c *CustomT) Error(msgOrErr ...interface{}) {
	msg := msgOrErr[0]
	msgType := reflect.TypeOf(msg).String()

	switch msgType {
	case "string":
		c.failureMsg = msg.(string)
	case "error":
		c.failureMsg = msg.(error).Error()
	default:
		panicMsg := fmt.Sprintf("Unexpected call to CustomT.Error with type: %s", msgType)
		panic(panicMsg)
	}
}

func NewCustomT() *CustomT {
	return &CustomT{}
}

func TestAssertions(t *testing.T) {
	t.Run("Match", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			ct := NewCustomT()
			assert.Match(ct, "foo", "sdffoosdf")
			if ct.failureMsg != "" {
				t.Errorf("Unexpected failure %s", ct.failureMsg)
			}
		})

		t.Run("Failure message", func(t *testing.T) {
			ct := NewCustomT()
			assert.Match(ct, "foo", "sdf")
			if ct.failureMsg != "Expected: \"foo\", but received: \"sdf\"" {
				t.Error(ct.failureMsg)

			}
		})
	})

	t.Run("NotNil", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			ct := NewCustomT()
			assert.NotNil(ct, true)
			if ct.failureMsg != "" {
				t.Error(ct.failureMsg)
			}

		})

		t.Run("Failure message", func(t *testing.T) {
			ct := NewCustomT()
			assert.NotNil(ct, nil)
			if ct.failureMsg == "" {
				t.Error("Expected failure")
			}
			assert.Match(t, "not be nil", ct.failureMsg)
		})
	})

	t.Run("Nil", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			ct := NewCustomT()
			assert.Nil(ct, nil)
			if ct.failureMsg != "" {
				t.Error(ct.failureMsg)
			}
		})

	})

	t.Run("True", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			ct := NewCustomT()
			assert.True(ct, true)
			if ct.failureMsg != "" {
				t.Error(ct.failureMsg)
			}
		})
	})

	t.Run("False", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			ct := NewCustomT()
			assert.False(ct, false)
			if ct.failureMsg != "" {
				t.Error(ct.failureMsg)
			}
		})

		t.Run("Failure message", func(t *testing.T) {
			ct := NewCustomT()
			assert.False(ct, true)
			if ct.failureMsg == "" {
				t.Error("Expected a failure message")
			}
			assert.Match(t, "Expected true to be false", ct.failureMsg)
		})
	})

	t.Run("StrictEqual", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			ct := NewCustomT()
			assert.StrictEqual(ct, 0.0, 0.0)
			if ct.failureMsg != "" {
				t.Error(ct.failureMsg)
			}
		})

		t.Run("Failure", func(t *testing.T) {
			ct := NewCustomT()
			assert.StrictEqual(ct, 0.0, 0)
			if ct.failureMsg == "" {
				t.Error("Expected StrictEqual failure")
			}

		})
	})

	t.Run("Equal", func(t *testing.T) {
		t.Run("0.0 == 0.0", func(t *testing.T) {
			ct := NewCustomT()
			assert.Equal(ct, 0.0, 0.0)
			if ct.failureMsg != "" {
				t.Error(ct.failureMsg)
			}
		})

		t.Run("0.0 == 0", func(t *testing.T) {
			ct := NewCustomT()
			assert.Equal(ct, 0.0, 0)
			if ct.failureMsg != "" {
				t.Error(ct.failureMsg)
			}
		})

		t.Run("0 == 0", func(t *testing.T) {
			ct := NewCustomT()
			assert.Equal(ct, 0, 0)
			if ct.failureMsg != "" {
				t.Error(ct.failureMsg)
			}
		})

		t.Run("0 == 0.0", func(t *testing.T) {
			ct := NewCustomT()
			assert.Equal(ct, 0, 0.0)
			if ct.failureMsg != "" {
				t.Error(ct.failureMsg)
			}
		})

		t.Run("Enum values match", func(t *testing.T) {
			ct := NewCustomT()
			assert.Equal(ct, fakeConstValueB, fakeConstValueB)
			if ct.failureMsg != "" {
				t.Error(ct)
			}
		})

		t.Run("Enum values mismatch", func(t *testing.T) {
			ct := NewCustomT()
			assert.Equal(ct, fakeConstValueA, fakeConstValueB)
			assert.Match(t, "expected 0 to equal 1", ct.failureMsg)
		})

		t.Run("failure with custom message", func(t *testing.T) {
			ct := NewCustomT()

			assert.Equal(ct, 1, 2, "Fake custom message")
			assert.Match(t, "Fake custom message", ct.failureMsg)
		})
	})
}
