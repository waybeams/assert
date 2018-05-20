package assert_test

import (
	"fmt"
	"github.com/waybeams/assert"
	"testing"
)

type foo struct{}

type fakeConst int

const (
	fakeConstValueA = iota
	fakeConstValueB
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
			assert.Match("foo", "sdffoosdf")
		})

		t.Run("Failure", func(t *testing.T) {
			defer ensurePanicWith("Expected: \"foo\", but received: \"sdf\"")()
			assert.Match("foo", "sdf")
		})
	})

	t.Run("NotNil", func(t *testing.T) {
		// The following is NOT TRUE: returnNil() != nil
		// This makes me such a sad panda.
		var returnNil = func() *foo {
			return nil
		}

		t.Run("Object results should be nil", func(t *testing.T) {
			t.Skip("Please help make the following call fail (and the test pass)")
			defer ensurePanicWith("Expected nil return value to fail NoNil check")()
			assert.NotNil(returnNil())
		})

		t.Run("Failure", func(t *testing.T) {
			defer ensurePanicWith("Expected <nil> to not be nil")()
			assert.NotNil(nil)

		})

		t.Run("Success", func(t *testing.T) {
			assert.NotNil(true)
		})
	})

	t.Run("Failure message", func(t *testing.T) {
		defer ensurePanicWith("Expected <nil> to not be nil")()
		assert.NotNil(nil)
	})

	t.Run("Nil", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			assert.Nil(nil)
		})

		t.Run("Failure", func(t *testing.T) {
			defer ensurePanicWith("Expected true of type: bool to be nil")()
			assert.Nil(true)
		})
	})

	t.Run("True", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			assert.True(true)
		})

		t.Run("Failure", func(t *testing.T) {
			defer ensurePanicWith("Expected false to be true")()
			assert.True(false)
		})
	})

	t.Run("False", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			assert.False(false)
		})

		t.Run("Failure message", func(t *testing.T) {
			defer ensurePanicWith("Expected true to be false")()
			assert.False(true)
		})
	})

	t.Run("StrictEqual", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			assert.StrictEqual(0.0, 0.0)
		})

		t.Run("Failure", func(t *testing.T) {
			defer ensurePanicWith("Expected 0 to STRICTLY equal 0")()
			assert.StrictEqual(0.0, 0)
		})
	})

	t.Run("Equal", func(t *testing.T) {
		t.Run("0.0 == 0.0", func(t *testing.T) {
			assert.Equal(0.0, 0.0)
		})

		t.Run("0.0 == 0", func(t *testing.T) {
			assert.Equal(0.0, 0)
		})

		t.Run("0 == 0.0", func(t *testing.T) {
			assert.Equal(0, 0.0)
		})

		t.Run("0 == 0", func(t *testing.T) {
			assert.Equal(0, 0)
		})

		t.Run("Enum values match", func(t *testing.T) {
			assert.Equal(fakeConstValueB, fakeConstValueB)
		})

		t.Run("Enum values mismatch", func(t *testing.T) {
			defer ensurePanicWith("Expected 0 to equal 1")()
			assert.Equal(fakeConstValueA, fakeConstValueB)
		})

		t.Run("failure with custom message", func(t *testing.T) {
			defer ensurePanicWith("Expected 1 to equal 2. Fake custom message")()
			assert.Equal(1, 2, "Fake custom message")
		})
	})
}
