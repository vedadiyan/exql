package math

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/vedadiyan/exql/lang"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestAbs(t *testing.T) {
	_, fn := Abs()

	tests := []struct {
		name     string
		input    float64
		expected float64
		hasError bool
	}{
		{"positive number", 5.5, 5.5, false},
		{"negative number", -5.5, 5.5, false},
		{"zero", 0, 0, false},
		{"large negative", -1000.5, 1000.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestStdDev(t *testing.T) {
	_, fn := StdDev()

	tests := []struct {
		name     string
		inputs   []lang.Value
		expected float64
		hasError bool
	}{
		{
			"simple stddev",
			[]lang.Value{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)},
			1,
			false,
		},
		{
			"single value",
			[]lang.Value{lang.NumberValue(5)},
			0,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.inputs)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestRandom(t *testing.T) {
	_, fn := Random()

	tests := []struct {
		name     string
		args     []lang.Value
		hasError bool
	}{
		{"no args", []lang.Value{}, false},
		{"single max", []lang.Value{lang.NumberValue(10)}, false},
		{"min and max", []lang.Value{lang.NumberValue(5), lang.NumberValue(10)}, false},
		{"too many args", []lang.Value{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Validate range
			val := float64(result.(lang.NumberValue))
			if len(tt.args) == 0 {
				if val < 0 || val >= 1 {
					t.Errorf("Random value %f out of range [0,1)", val)
				}
			} else if len(tt.args) == 1 {
				max := float64(tt.args[0].(lang.NumberValue))
				if val < 0 || val >= max {
					t.Errorf("Random value %f out of range [0,%f)", val, max)
				}
			} else if len(tt.args) == 2 {
				min := float64(tt.args[0].(lang.NumberValue))
				max := float64(tt.args[1].(lang.NumberValue))
				if val < min || val >= max {
					t.Errorf("Random value %f out of range [%f,%f)", val, min, max)
				}
			}
		})
	}
}

func TestRandomSeed(t *testing.T) {
	_, fn := RandomSeed()

	tests := []struct {
		name     string
		args     []lang.Value
		hasError bool
	}{
		{"no seed", []lang.Value{}, false},
		{"with seed", []lang.Value{lang.NumberValue(12345)}, false},
		{"too many args", []lang.Value{lang.NumberValue(1), lang.NumberValue(2)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if !bool(result.(lang.BoolValue)) {
				t.Errorf("Expected true result")
			}
		})
	}
}

func TestRandomFloat(t *testing.T) {
	_, fn := RandomFloat()

	tests := []struct {
		name     string
		args     []lang.Value
		hasError bool
	}{
		{"no args", []lang.Value{}, false},
		{"single max", []lang.Value{lang.NumberValue(5.5)}, false},
		{"min and max", []lang.Value{lang.NumberValue(1.5), lang.NumberValue(5.5)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			val := float64(result.(lang.NumberValue))
			if len(tt.args) == 0 {
				if val < 0 || val >= 1 {
					t.Errorf("Random float %f out of range [0,1)", val)
				}
			}
		})
	}
}

func TestIsNaN(t *testing.T) {
	_, fn := IsNaN()

	tests := []struct {
		name     string
		input    float64
		expected bool
	}{
		{"normal number", 5.5, false},
		{"zero", 0, false},
		{"infinity", math.Inf(1), false},
		{"NaN", math.NaN(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestIsInf(t *testing.T) {
	_, fn := IsInf()

	tests := []struct {
		name     string
		input    float64
		expected bool
	}{
		{"normal number", 5.5, false},
		{"zero", 0, false},
		{"positive infinity", math.Inf(1), true},
		{"negative infinity", math.Inf(-1), true},
		{"NaN", math.NaN(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestIsFinite(t *testing.T) {
	_, fn := IsFinite()

	tests := []struct {
		name     string
		input    float64
		expected bool
	}{
		{"normal number", 5.5, true},
		{"zero", 0, true},
		{"large number", 1e308, true},
		{"infinity", math.Inf(1), false},
		{"NaN", math.NaN(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestGCD(t *testing.T) {
	_, fn := GCD()

	tests := []struct {
		name     string
		inputs   []float64
		expected float64
		hasError bool
	}{
		{"two numbers", []float64{12, 8}, 4, false},
		{"three numbers", []float64{12, 8, 16}, 4, false},
		{"coprime numbers", []float64{7, 13}, 1, false},
		{"with zero", []float64{12, 0}, 12, false},
		{"negative numbers", []float64{-12, 8}, 4, false},
		{"single argument", []float64{12}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]lang.Value, len(tt.inputs))
			for i, input := range tt.inputs {
				args[i] = lang.NumberValue(input)
			}

			result, err := fn(args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestLCM(t *testing.T) {
	_, fn := LCM()

	tests := []struct {
		name     string
		inputs   []float64
		expected float64
		hasError bool
	}{
		{"two numbers", []float64{12, 8}, 24, false},
		{"three numbers", []float64{4, 6, 8}, 24, false},
		{"coprime numbers", []float64{7, 13}, 91, false},
		{"with zero", []float64{12, 0}, 0, false},
		{"single argument", []float64{12}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]lang.Value, len(tt.inputs))
			for i, input := range tt.inputs {
				args[i] = lang.NumberValue(input)
			}

			result, err := fn(args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestFactorial(t *testing.T) {
	_, fn := Factorial()

	tests := []struct {
		name     string
		input    float64
		expected float64
		isNaN    bool
	}{
		{"zero factorial", 0, 1, false},
		{"one factorial", 1, 1, false},
		{"small factorial", 5, 120, false},
		{"negative factorial", -1, 0, true}, // Should return NaN
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if tt.isNaN {
				if !math.IsNaN(actual) {
					t.Errorf("Expected NaN, got %f", actual)
				}
			} else {
				if actual != tt.expected {
					t.Errorf("Expected %f, got %f", tt.expected, actual)
				}
			}
		})
	}
}

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() (string, lang.Function)
		expected float64
	}{
		{"Pi", Pi, math.Pi},
		{"E", E, math.E},
		{"Phi", Phi, math.Phi},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, fn := tt.fn()
			result, err := fn([]lang.Value{})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if math.Abs(actual-tt.expected) > 1e-15 {
				t.Errorf("Expected %f, got %f", tt.expected, actual)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("gcd helper", func(t *testing.T) {
		tests := []struct {
			a, b     int
			expected int
		}{
			{12, 8, 4},
			{7, 13, 1},
			{0, 5, 5},
			{5, 0, 5},
		}

		for _, tt := range tests {
			result := gcd(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("gcd(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		}
	})

	t.Run("lcm helper", func(t *testing.T) {
		tests := []struct {
			a, b     int
			expected int
		}{
			{12, 8, 24},
			{7, 13, 91},
			{0, 5, 0},
			{5, 0, 0},
		}

		for _, tt := range tests {
			result := lcm(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("lcm(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		}
	})
}

func TestArgumentErrors(t *testing.T) {
	functions := []struct {
		fn           func() (string, lang.Function)
		expectedArgs int
	}{
		{Abs, 1},
		{Sign, 1},
		{Clamp, 3},
		{Ceil, 1},
		{Floor, 1},
		{Trunc, 1},
		{Pow, 2},
		{Sqrt, 1},
		{Cbrt, 1},
		{Exp, 1},
		{Log, 1},
		{Sin, 1},
		{Cos, 1},
		{Tan, 1},
		{Atan2, 2},
		{Radians, 1},
		{Degrees, 1},
		{IsNaN, 1},
		{IsInf, 1},
		{IsFinite, 1},
		{Factorial, 1},
		{Pi, 0},
		{E, 0},
		{Phi, 0},
	}

	for _, tf := range functions {
		name, fn := tf.fn()
		t.Run(name+"_wrong_args", func(t *testing.T) {
			var args []lang.Value
			if tf.expectedArgs == 0 {
				args = []lang.Value{lang.NumberValue(1)} // Add extra arg for 0-arg functions
			} else {
				args = make([]lang.Value, tf.expectedArgs-1) // One less than expected
				for i := range args {
					args[i] = lang.NumberValue(1)
				}
			}
			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for wrong argument count in %s", name)
			}
		})
	}
}

func TestNonNumericInputs(t *testing.T) {
	functions := []func() (string, lang.Function){
		Abs, Sign, Max, Min, Ceil, Floor, Round, Trunc,
		Sqrt, Cbrt, Exp, Log, Sin, Cos, Tan, Asin, Acos, Atan,
		Sinh, Cosh, Tanh, Radians, Degrees, IsNaN, IsInf, IsFinite, Factorial,
	}

	nonNumericInput := lang.StringValue("not a number")

	for _, getFn := range functions {
		name, fn := getFn()
		t.Run(name+"_non_numeric", func(t *testing.T) {
			_, err := fn([]lang.Value{nonNumericInput})
			if err == nil {
				t.Errorf("Expected error for non-numeric input in %s", name)
			}
		})
	}
}

func TestExport(t *testing.T) {
	functions := Export()

	expectedFunctions := []string{
		"abs", "sign", "max", "min", "clamp", "ceil", "floor", "round", "trunc",
		"pow", "sqrt", "cbrt", "exp", "exp2", "log", "log10", "log2",
		"sin", "cos", "tan", "asin", "acos", "atan", "atan2",
		"sinh", "cosh", "tanh", "radians", "degrees",
		"sum", "mean", "median", "mode", "variance", "stddev",
		"random", "randomSeed", "randomFloat",
		"isNan", "isInf", "isFinite", "gcd", "lcm", "factorial",
		"pi", "e", "phi",
	}

	if len(functions) != len(expectedFunctions) {
		t.Errorf("Expected %d functions, got %d", len(expectedFunctions), len(functions))
	}

	for _, name := range expectedFunctions {
		if _, exists := functions[name]; !exists {
			t.Errorf("Expected function %s not found", name)
		}
	}
}

func BenchmarkAbs(b *testing.B) {
	_, fn := Abs()
	args := []lang.Value{lang.NumberValue(-42.5)}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkSqrt(b *testing.B) {
	_, fn := Sqrt()
	args := []lang.Value{lang.NumberValue(144)}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkSin(b *testing.B) {
	_, fn := Sin()
	args := []lang.Value{lang.NumberValue(math.Pi / 4)}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkSum(b *testing.B) {
	_, fn := Sum()
	args := []lang.Value{
		lang.ListValue{
			lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3),
			lang.NumberValue(4), lang.NumberValue(5), lang.NumberValue(6),
		},
	}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func intPtr(i int) *int {
	return &i
}

func TestSign(t *testing.T) {
	_, fn := Sign()

	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"positive", 5.5, 1},
		{"negative", -5.5, -1},
		{"zero", 0, 0},
		{"positive small", 0.001, 1},
		{"negative small", -0.001, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestMax(t *testing.T) {
	_, fn := Max()

	tests := []struct {
		name     string
		inputs   []float64
		expected float64
		hasError bool
	}{
		{"two numbers", []float64{5, 10}, 10, false},
		{"three numbers", []float64{5, 10, 3}, 10, false},
		{"single number", []float64{42}, 42, false},
		{"negative numbers", []float64{-5, -10, -3}, -3, false},
		{"mixed positive negative", []float64{-5, 10, -3}, 10, false},
		{"no arguments", []float64{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]lang.Value, len(tt.inputs))
			for i, input := range tt.inputs {
				args[i] = lang.NumberValue(input)
			}

			result, err := fn(args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestMin(t *testing.T) {
	_, fn := Min()

	tests := []struct {
		name     string
		inputs   []float64
		expected float64
		hasError bool
	}{
		{"two numbers", []float64{5, 10}, 5, false},
		{"three numbers", []float64{5, 10, 3}, 3, false},
		{"single number", []float64{42}, 42, false},
		{"negative numbers", []float64{-5, -10, -3}, -10, false},
		{"no arguments", []float64{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]lang.Value, len(tt.inputs))
			for i, input := range tt.inputs {
				args[i] = lang.NumberValue(input)
			}

			result, err := fn(args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestClamp(t *testing.T) {
	_, fn := Clamp()

	tests := []struct {
		name     string
		value    float64
		min      float64
		max      float64
		expected float64
	}{
		{"within range", 5, 0, 10, 5},
		{"below min", -5, 0, 10, 0},
		{"above max", 15, 0, 10, 10},
		{"at min", 0, 0, 10, 0},
		{"at max", 10, 0, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{
				lang.NumberValue(tt.value),
				lang.NumberValue(tt.min),
				lang.NumberValue(tt.max),
			}
			result, err := fn(args)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestCeil(t *testing.T) {
	_, fn := Ceil()

	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"positive decimal", 4.2, 5},
		{"negative decimal", -4.2, -4},
		{"whole number", 5, 5},
		{"zero", 0, 0},
		{"small positive", 0.1, 1},
		{"small negative", -0.1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestFloor(t *testing.T) {
	_, fn := Floor()

	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"positive decimal", 4.8, 4},
		{"negative decimal", -4.2, -5},
		{"whole number", 5, 5},
		{"zero", 0, 0},
		{"small positive", 0.9, 0},
		{"small negative", -0.1, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestRound(t *testing.T) {
	_, fn := Round()

	tests := []struct {
		name      string
		value     float64
		precision *int
		expected  float64
		hasError  bool
	}{
		{"basic round", 4.5, nil, 5, false},
		{"round down", 4.4, nil, 4, false},
		{"with precision", 4.567, intPtr(2), 4.57, false},
		{"zero precision", 4.567, intPtr(0), 5, false},
		{"negative precision", 123.456, intPtr(-1), 120, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.NumberValue(tt.value)}
			if tt.precision != nil {
				args = append(args, lang.NumberValue(float64(*tt.precision)))
			}

			result, err := fn(args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if math.Abs(actual-tt.expected) > 1e-10 {
				t.Errorf("Expected %f, got %f", tt.expected, actual)
			}
		})
	}
}

func TestTrunc(t *testing.T) {
	_, fn := Trunc()

	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"positive decimal", 4.8, 4},
		{"negative decimal", -4.8, -4},
		{"whole number", 5, 5},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestPow(t *testing.T) {
	_, fn := Pow()

	tests := []struct {
		name     string
		base     float64
		exponent float64
		expected float64
	}{
		{"square", 2, 2, 4},
		{"cube", 3, 3, 27},
		{"power of zero", 5, 0, 1},
		{"zero base", 0, 3, 0},
		{"negative base even exponent", -2, 2, 4},
		{"negative base odd exponent", -2, 3, -8},
		{"fractional exponent", 9, 0.5, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.NumberValue(tt.base), lang.NumberValue(tt.exponent)}
			result, err := fn(args)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if math.Abs(actual-tt.expected) > 1e-10 {
				t.Errorf("Expected %f, got %f", tt.expected, actual)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	_, fn := Sqrt()

	tests := []struct {
		name     string
		input    float64
		expected float64
		isNaN    bool
	}{
		{"perfect square", 9, 3, false},
		{"non-perfect square", 8, 2.828427125, false},
		{"zero", 0, 0, false},
		{"negative number", -4, 0, true}, // Should return NaN
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if tt.isNaN {
				if !math.IsNaN(actual) {
					t.Errorf("Expected NaN, got %f", actual)
				}
			} else {
				if math.Abs(actual-tt.expected) > 1e-6 {
					t.Errorf("Expected %f, got %f", tt.expected, actual)
				}
			}
		})
	}
}

func TestCbrt(t *testing.T) {
	_, fn := Cbrt()

	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"perfect cube", 27, 3},
		{"negative cube", -8, -2},
		{"zero", 0, 0},
		{"one", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if math.Abs(actual-tt.expected) > 1e-10 {
				t.Errorf("Expected %f, got %f", tt.expected, actual)
			}
		})
	}
}

func TestTrigonometricFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() (string, lang.Function)
		input    float64
		expected float64
	}{
		{"sin(0)", Sin, 0, 0},
		{"sin(π/2)", Sin, math.Pi / 2, 1},
		{"cos(0)", Cos, 0, 1},
		{"cos(π)", Cos, math.Pi, -1},
		{"tan(0)", Tan, 0, 0},
		{"tan(π/4)", Tan, math.Pi / 4, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, fn := tt.fn()
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if math.Abs(actual-tt.expected) > 1e-10 {
				t.Errorf("Expected %f, got %f", tt.expected, actual)
			}
		})
	}
}

func TestInverseTrigonometricFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() (string, lang.Function)
		input    float64
		expected float64
		isNaN    bool
	}{
		{"asin(0)", Asin, 0, 0, false},
		{"asin(1)", Asin, 1, math.Pi / 2, false},
		{"asin(2)", Asin, 2, 0, true}, // Out of range
		{"acos(0)", Acos, 0, math.Pi / 2, false},
		{"acos(1)", Acos, 1, 0, false},
		{"acos(2)", Acos, 2, 0, true}, // Out of range
		{"atan(0)", Atan, 0, 0, false},
		{"atan(1)", Atan, 1, math.Pi / 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, fn := tt.fn()
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if tt.isNaN {
				if !math.IsNaN(actual) {
					t.Errorf("Expected NaN, got %f", actual)
				}
			} else {
				if math.Abs(actual-tt.expected) > 1e-10 {
					t.Errorf("Expected %f, got %f", tt.expected, actual)
				}
			}
		})
	}
}

func TestAtan2(t *testing.T) {
	_, fn := Atan2()

	tests := []struct {
		name     string
		y        float64
		x        float64
		expected float64
	}{
		{"first quadrant", 1, 1, math.Pi / 4},
		{"second quadrant", 1, -1, 3 * math.Pi / 4},
		{"third quadrant", -1, -1, -3 * math.Pi / 4},
		{"fourth quadrant", -1, 1, -math.Pi / 4},
		{"positive x-axis", 0, 1, 0},
		{"negative x-axis", 0, -1, math.Pi},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.NumberValue(tt.y), lang.NumberValue(tt.x)}
			result, err := fn(args)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if math.Abs(actual-tt.expected) > 1e-10 {
				t.Errorf("Expected %f, got %f", tt.expected, actual)
			}
		})
	}
}

func TestHyperbolicFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() (string, lang.Function)
		input    float64
		expected float64
	}{
		{"sinh(0)", Sinh, 0, 0},
		{"cosh(0)", Cosh, 0, 1},
		{"tanh(0)", Tanh, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, fn := tt.fn()
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if math.Abs(actual-tt.expected) > 1e-10 {
				t.Errorf("Expected %f, got %f", tt.expected, actual)
			}
		})
	}
}

func TestRadiansAndDegrees(t *testing.T) {
	t.Run("radians", func(t *testing.T) {
		_, fn := Radians()
		result, err := fn([]lang.Value{lang.NumberValue(180)})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		actual := float64(result.(lang.NumberValue))
		if math.Abs(actual-math.Pi) > 1e-10 {
			t.Errorf("Expected %f, got %f", math.Pi, actual)
		}
	})

	t.Run("degrees", func(t *testing.T) {
		_, fn := Degrees()
		result, err := fn([]lang.Value{lang.NumberValue(math.Pi)})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		actual := float64(result.(lang.NumberValue))
		if math.Abs(actual-180) > 1e-10 {
			t.Errorf("Expected %f, got %f", 180.0, actual)
		}
	})
}

func TestLogarithmicFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() (string, lang.Function)
		input    float64
		expected float64
		isNaN    bool
	}{
		{"log(e)", Log, math.E, 1, false},
		{"log(1)", Log, 1, 0, false},
		{"log(0)", Log, 0, 0, true},   // Should be NaN
		{"log(-1)", Log, -1, 0, true}, // Should be NaN
		{"log10(10)", Log10, 10, 1, false},
		{"log10(100)", Log10, 100, 2, false},
		{"log2(2)", Log2, 2, 1, false},
		{"log2(8)", Log2, 8, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, fn := tt.fn()
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if tt.isNaN {
				if !math.IsNaN(actual) {
					t.Errorf("Expected NaN, got %f", actual)
				}
			} else {
				if math.Abs(actual-tt.expected) > 1e-10 {
					t.Errorf("Expected %f, got %f", tt.expected, actual)
				}
			}
		})
	}
}

func TestExponentialFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() (string, lang.Function)
		input    float64
		expected float64
	}{
		{"exp(0)", Exp, 0, 1},
		{"exp(1)", Exp, 1, math.E},
		{"exp2(0)", Exp2, 0, 1},
		{"exp2(3)", Exp2, 3, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, fn := tt.fn()
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if math.Abs(actual-tt.expected) > 1e-10 {
				t.Errorf("Expected %f, got %f", tt.expected, actual)
			}
		})
	}
}

func TestSum(t *testing.T) {
	_, fn := Sum()

	tests := []struct {
		name     string
		inputs   []lang.Value
		expected float64
		hasError bool
	}{
		{
			"numbers only",
			[]lang.Value{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)},
			6,
			false,
		},
		{
			"with list",
			[]lang.Value{lang.NumberValue(1), lang.ListValue{lang.NumberValue(2), lang.NumberValue(3)}},
			6,
			false,
		},
		{
			"list only",
			[]lang.Value{lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)}},
			6,
			false,
		},
		{
			"empty",
			[]lang.Value{},
			0,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.inputs)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestMean(t *testing.T) {
	_, fn := Mean()

	tests := []struct {
		name     string
		inputs   []lang.Value
		expected float64
		hasError bool
	}{
		{
			"simple average",
			[]lang.Value{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)},
			2,
			false,
		},
		{
			"with list",
			[]lang.Value{lang.ListValue{lang.NumberValue(2), lang.NumberValue(4), lang.NumberValue(6)}},
			4,
			false,
		},
		{
			"no arguments",
			[]lang.Value{},
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.inputs)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestMedian(t *testing.T) {
	_, fn := Median()

	tests := []struct {
		name     string
		inputs   []lang.Value
		expected float64
		hasError bool
	}{
		{
			"odd count",
			[]lang.Value{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)},
			2,
			false,
		},
		{
			"even count",
			[]lang.Value{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3), lang.NumberValue(4)},
			2.5,
			false,
		},
		{
			"unsorted",
			[]lang.Value{lang.NumberValue(3), lang.NumberValue(1), lang.NumberValue(2)},
			2,
			false,
		},
		{
			"no arguments",
			[]lang.Value{},
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.inputs)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestMode(t *testing.T) {
	_, fn := Mode()

	tests := []struct {
		name     string
		inputs   []lang.Value
		expected float64
		hasError bool
	}{
		{
			"clear mode",
			[]lang.Value{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(2), lang.NumberValue(3)},
			2,
			false,
		},
		{
			"all same",
			[]lang.Value{lang.NumberValue(5), lang.NumberValue(5), lang.NumberValue(5)},
			5,
			false,
		},
		{
			"no arguments",
			[]lang.Value{},
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.inputs)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestVariance(t *testing.T) {
	_, fn := Variance()

	tests := []struct {
		name     string
		inputs   []lang.Value
		expected float64
		hasError bool
	}{
		{
			"simple variance",
			[]lang.Value{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)},
			1,
			false,
		},
		{
			"single value",
			[]lang.Value{lang.NumberValue(5)},
			0,
			false,
		},
		{
			"no arguments",
			[]lang.Value{},
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.inputs)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}
