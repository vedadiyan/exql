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
package exql

import (
	"testing"

	"github.com/vedadiyan/exql/lang"
)

func TestContextOptions(t *testing.T) {
	// Test different context configurations
	t.Run("empty context", func(t *testing.T) {
		ctx := NewDefaultContext()

		result, err := Eval("42", ctx)
		if err != nil {
			t.Fatalf("error with empty context: %v", err)
		}

		if !valueEqual(result, lang.NumberValue(42)) {
			t.Errorf("expected 42, got %v", result)
		}
	})

	t.Run("context with built-in library", func(t *testing.T) {
		ctx := NewDefaultContext(WithBuiltInLibrary())

		result, err := Eval("util", ctx)
		if err != nil {
			t.Fatalf("error accessing built-in library: %v", err)
		}

		if result == nil {
			t.Error("built-in library should not be nil")
		}
	})

	t.Run("context with custom functions", func(t *testing.T) {
		customFuncs := map[string]lang.Function{
			"double": func(args []lang.Value) (lang.Value, error) {
				return lang.NumberValue(toNumber(args[0]) * 2), nil
			},
		}

		ctx := NewDefaultContext(WithFunctions(customFuncs))

		result, err := Eval("double(21)", ctx)
		if err != nil {
			t.Fatalf("error with custom function: %v", err)
		}

		if !valueEqual(result, lang.NumberValue(42)) {
			t.Errorf("expected 42, got %v", result)
		}
	})

	t.Run("context with both options", func(t *testing.T) {
		customFuncs := map[string]lang.Function{
			"test": func(args []lang.Value) (lang.Value, error) {
				return lang.StringValue("success"), nil
			},
		}

		ctx := NewDefaultContext(WithBuiltInLibrary(), WithFunctions(customFuncs))

		// Test custom function
		result1, err := Eval("test()", ctx)
		if err != nil {
			t.Fatalf("error with custom function: %v", err)
		}

		if !valueEqual(result1, lang.StringValue("success")) {
			t.Errorf("expected 'success', got %v", result1)
		}

		// Test built-in library
		result2, err := Eval("string", ctx)
		if err != nil {
			t.Fatalf("error with built-in library: %v", err)
		}

		if result2 == nil {
			t.Error("built-in library should not be nil")
		}
	})
}
