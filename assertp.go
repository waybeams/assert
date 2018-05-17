// Package assert includes helper functions that work with the native Go
// testing package.
package assert

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

func messagesToStringP(mainMessage string, optMessages ...string) string {
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

func isTrueP(value bool, mainMessage string, messages ...string) {
	if !value {
		msg := messagesToStringP(mainMessage, messages...)
		panic(errors.New(msg).Error())
	}
}

// Panic fails if the provided handler does not trigger a panic that includes an error
// or message that matches the provided expression string.
func PanicP(expr string, handler func()) {
	defer func() {
		r := recover()
		if r != nil {
			var err error
			fmt.Println("Recovered in f", r)
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			MatchP(expr, err.Error())
		} else {
			panic("Did not receive expected panic")
		}
	}()
	handler()
}

// StrictEqual fails if the provided values are not == to one another.
func StrictEqualP(found interface{}, expected interface{}, messages ...string) {
	if found != expected {
		mainMessage := fmt.Sprintf("Expected %v to STRICTLY equal %v", found, expected)
		panic(messagesToStringP(mainMessage, messages...))
	}
}

// Equal fails if the provided values are not equal in a "best effort" comparison.
// This method will (perhaps incorrectly to reasonably folks) claim 1.0 is
// equal to 1.
// This coercion helps with test brevity and flexibility. If you'd like
// something more precise, use StrictEqual instead.
func EqualP(found interface{}, expected interface{}, messages ...string) {
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
				panic(messagesToStringP(mainMessage, messages...))
			}
			return
		}

		if found != expected {
			mainMessage := fmt.Sprintf("Custom Equal expected %v to equal %v", found, expected)
			panic(messagesToStringP(mainMessage, messages...))
		}
	}
}

// Match fails if the the provided exprStr is not found in the provided str value as
// a regular expression.
func MatchP(exprStr string, str string) {
	matched, _ := regexp.MatchString(exprStr, str)
	fmt.Println("MATCH P WITH:", matched)
	if !matched {
		panic(fmt.Sprintf("Expected: \"%v\", but received: \"%v\"", exprStr, str))
	}
}

// True fails if the provided value is not true
func TrueP(value bool, messages ...string) {
	isTrueP(value, fmt.Sprintf("Expected %v to be true", value), messages...)
}

// False fails if the provided value is not false
func FalseP(value bool, messages ...string) {
	isTrueP(!value, fmt.Sprintf("Expected %v to be false", value), messages...)
}

// NotNil fails if the provided value is nil
func NotNilP(value interface{}, messages ...string) {
	if value == nil {
		msg := fmt.Sprintf("Expected %v to not be nil", value)
		panic(messagesToStringP(msg, messages...))
	}
}

// Nil fails if the provided value is not nil
func NilP(value interface{}, messages ...string) {
	if value != nil {
		typeOf := reflect.TypeOf(value).String()
		msg := fmt.Sprintf("Expected %v of type: %v to be nil", value, typeOf)
		panic(messagesToStringP(msg, messages...))
	}
}
