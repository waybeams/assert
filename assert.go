// Package assert includes helper functions that work with the native Go
// testing package.
package assert

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

func messagesToString(mainMessage string, optMessages ...string) string {
	switch len(optMessages) {
	case 0:
		return mainMessage
	case 1:
		return fmt.Sprintf("%s. %s", mainMessage, optMessages[0])
	case 2:
		return fmt.Sprintf("%s %s\n%s", mainMessage, optMessages[0], optMessages[1])
	default:
		panic("Custom assertion provided with unexpected messages")
	}
	return ""
}

func isTrue(value bool, mainMessage string, messages ...string) {
	if !value {
		msg := messagesToString(mainMessage, messages...)
		panic(errors.New(msg).Error())
	}
}

// Panic fails if the provided handler does not trigger a panic that includes an error
// or message that matches the provided expression string.
func Panic(expr string, handler func()) {
	defer func() {
		r := recover()
		if r != nil {
			var err error
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			Match(expr, err.Error())
		} else {
			panic("Did not receive expected panic")
		}
	}()
	handler()
}

// StrictEqual fails if the provided values are not == to one another.
func StrictEqual(found interface{}, expected interface{}, messages ...string) {
	if found != expected {
		mainMessage := fmt.Sprintf("Expected %v to STRICTLY equal %v", found, expected)
		panic(messagesToString(mainMessage, messages...))
	}
}

// Equal fails if the provided values are not equal in a "best effort" comparison.
// This method will (perhaps incorrectly to reasonably folks) claim 1.0 is
// equal to 1.
// This coercion helps with test brevity and flexibility. If you'd like
// something more precise, use StrictEqual instead.
func Equal(found interface{}, expected interface{}, messages ...string) {
	if found != expected {
		kindA := reflect.ValueOf(found).Kind()
		switch kindA {
		case reflect.Bool:
			fallthrough
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			foundStr := fmt.Sprintf("%v", found)
			expectedStr := fmt.Sprintf("%v", expected)
			if foundStr != expectedStr {
				mainMessage := fmt.Sprintf("Expected %v to equal %v", found, expected)
				panic(messagesToString(mainMessage, messages...))
			}
			return
		}

		if found != expected {
			mainMessage := fmt.Sprintf("Custom Equal expected %v to equal %v", found, expected)
			panic(messagesToString(mainMessage, messages...))
		}
	}
}

// Match fails if the the provided exprStr is not found in the provided str value as
// a regular expression.
func Match(exprStr string, str string) {
	matched, _ := regexp.MatchString(exprStr, str)
	if !matched {
		panic(fmt.Sprintf("Expected: \"%v\", but received: \"%v\"", exprStr, str))
	}
}

// True fails if the provided value is not true
func True(value bool, messages ...string) {
	isTrue(value, fmt.Sprintf("Expected %v to be true", value), messages...)
}

// False fails if the provided value is not false
func False(value bool, messages ...string) {
	isTrue(!value, fmt.Sprintf("Expected %v to be false", value), messages...)
}

// NotNil fails if the provided value is nil
func NotNil(value interface{}, messages ...string) {
	if isNil(value) {
		msg := fmt.Sprintf("Expected %v to not be nil", value)
		panic(messagesToString(msg, messages...))
	}
}

func isNil(actual interface{}) bool {
	if actual == nil {
		return true
	}

	value := reflect.ValueOf(actual)
	kind := value.Kind()
	nilable := kind == reflect.Slice ||
		kind == reflect.Chan ||
		kind == reflect.Func ||
		kind == reflect.Ptr ||
		kind == reflect.Map

	// Careful: reflect.Value.IsNil() will panic unless it's an interface, chan, map, func, slice, or ptr
	// Reference: http://golang.org/pkg/reflect/#Value.IsNil
	return nilable && value.IsNil()
}

// Nil fails if the provided value is not nil
func Nil(value interface{}, messages ...string) {
	if !isNil(value) {
		typeOf := reflect.TypeOf(value).String()
		msg := fmt.Sprintf("Expected %v of type: %v to be nil", value, typeOf)
		panic(messagesToString(msg, messages...))
	}
}
