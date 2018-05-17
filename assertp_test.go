package assert_test

import (
	"fmt"
	"github.com/waybeams/assert"
	"testing"
)

func ensurePanicWith(expectedMessage string) func() {
	return func() {
		if r := recover(); r != nil {
			if r != expectedMessage {
				panic(fmt.Sprintf("Received unexpected panic: %v", r))
			}
		} else {
			panic(fmt.Sprintf("Expected panic (%v) but did not receive one", expectedMessage))
		}
	}
}

func TestPanicAsserts(t *testing.T) {
	t.Run("Match", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			assert.MatchP("foo", "sdffoosdf")
		})

		t.Run("Failure", func(t *testing.T) {
			defer ensurePanicWith("Expected: \"foo\", but received: \"sdf\"")()
			assert.MatchP("foo", "sdf")
		})
	})

	t.Run("NotNil", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			assert.NotNilP(true)
		})
	})

	t.Run("Failure message", func(t *testing.T) {
		defer ensurePanicWith("Expected <nil> to not be nil")()
		assert.NotNilP(nil)
	})

	t.Run("Nil", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			assert.NilP(nil)
		})

		t.Run("Failure", func(t *testing.T) {
			defer ensurePanicWith("Expected true of type: bool to be nil")()
			assert.NilP(true)
		})
	})

	t.Run("True", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			assert.TrueP(true)
		})

		t.Run("Failure", func(t *testing.T) {
			defer ensurePanicWith("Expected false to be true")()
			assert.TrueP(false)
		})
	})

	t.Run("False", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			assert.FalseP(false)
		})

		t.Run("Failure message", func(t *testing.T) {
			defer ensurePanicWith("Expected true to be false")()
			assert.FalseP(true)
		})
	})

	t.Run("StrictEqual", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			assert.StrictEqualP(0.0, 0.0)
		})

		t.Run("Failure", func(t *testing.T) {
			defer ensurePanicWith("Expected 0 to STRICTLY equal 0")()
			assert.StrictEqualP(0.0, 0)
		})
	})

	t.Run("Equal", func(t *testing.T) {
		t.Run("0.0 == 0.0", func(t *testing.T) {
			assert.EqualP(0.0, 0.0)
		})

		t.Run("0.0 == 0", func(t *testing.T) {
			assert.EqualP(0.0, 0)
		})

		t.Run("0 == 0.0", func(t *testing.T) {
			assert.EqualP(0, 0.0)
		})

		t.Run("0 == 0", func(t *testing.T) {
			assert.EqualP(0, 0)
		})

		t.Run("Enum values match", func(t *testing.T) {
			assert.EqualP(fakeConstValueB, fakeConstValueB)
		})

		t.Run("Enum values mismatch", func(t *testing.T) {
			defer ensurePanicWith("Expected 0 to equal 1")()
			assert.EqualP(fakeConstValueA, fakeConstValueB)
		})

		t.Run("failure with custom message", func(t *testing.T) {
			defer ensurePanicWith("Expected 1 to equal 2. Fake custom message")()
			assert.EqualP(1, 2, "Fake custom message")
		})
	})
}
