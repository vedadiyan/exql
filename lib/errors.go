/*
 * Copyright 2025 Pouya Vedadiyan
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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
