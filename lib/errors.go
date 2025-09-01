package lib

import (
	"fmt"

	"github.com/vedadiyan/exql/lang"
)

func ArgumentError(name string, expected int) error {
	return fmt.Errorf("%s: expected %d arguments", name, expected)
}

func ContextError(name string, value lang.Value) error {
	return fmt.Errorf("%s: first argument must be map, got %T", name, value)
}

func StringError(name string, value lang.Value) error {
	return fmt.Errorf("%s: expected string, got %T", name, value)
}

func ListError(name string, value lang.Value) error {
	return fmt.Errorf("%s: expected list, got %T", name, value)
}

func MapError(name string, value lang.Value) error {
	return fmt.Errorf("%s: expected map, got %T", name, value)
}

func RangeError(name string, min, max int) error {
	return fmt.Errorf("%s: expected %d to %d arguments", name, min, max)
}

func ArgumentErrorMin(name string, expected int) error {
	return fmt.Errorf("%s: expected at least %d argument(s)", name, expected)
}

func ArgumentErrorRange(name string, min, max int) error {
	return fmt.Errorf("%s: expected %d or %d arguments", name, min, max)
}

func RgumentErrorMultiRange(name string, expected []int) error {
	return fmt.Errorf("%s: expected %v arguments", name, expected)
}

func ArgumentErrorMultiRange(name string, expected []int) error {
	return fmt.Errorf("%s: expected %v arguments", name, expected)
}
