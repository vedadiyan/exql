package math

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// Math Functions
// These functions provide comprehensive mathematical operations and calculations

// Basic Math Functions
func mathAbs(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("abs: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("abs: %w", err)
	}
	return lang.NumberValue(math.Abs(num)), nil
}

func mathSign(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("sign: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}
	if num > 0 {
		return lang.NumberValue(1), nil
	} else if num < 0 {
		return lang.NumberValue(-1), nil
	}
	return lang.NumberValue(0), nil
}

func mathMax(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		return nil, errors.New("max: expected at least 1 argument")
	}
	max, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("max: argument 0 %w", err)
	}
	for i := 1; i < len(args); i++ {
		val, err := lib.ToNumber(args[i])
		if err != nil {
			return nil, fmt.Errorf("max: argument %d %w", i, err)
		}
		if val > max {
			max = val
		}
	}
	return lang.NumberValue(max), nil
}

func mathMin(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		return nil, errors.New("min: expected at least 1 argument")
	}
	min, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("min: argument 0 %w", err)
	}
	for i := 1; i < len(args); i++ {
		val, err := lib.ToNumber(args[i])
		if err != nil {
			return nil, fmt.Errorf("min: argument %d %w", i, err)
		}
		if val < min {
			min = val
		}
	}
	return lang.NumberValue(min), nil
}

func mathClamp(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("clamp: expected 3 arguments")
	}
	value, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("clamp: value %w", err)
	}
	min, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("clamp: min %w", err)
	}
	max, err := lib.ToNumber(args[2])
	if err != nil {
		return nil, fmt.Errorf("clamp: max %w", err)
	}

	if value < min {
		return lang.NumberValue(min), nil
	}
	if value > max {
		return lang.NumberValue(max), nil
	}
	return lang.NumberValue(value), nil
}

// Rounding Functions
func mathCeil(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("ceil: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("ceil: %w", err)
	}
	return lang.NumberValue(math.Ceil(num)), nil
}

func mathFloor(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("floor: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("floor: %w", err)
	}
	return lang.NumberValue(math.Floor(num)), nil
}

func mathRound(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("round: expected 1 or 2 arguments")
	}

	value, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("round: value %w", err)
	}
	precision := 0

	if len(args) == 2 {
		precisionFloat, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("round: precision %w", err)
		}
		precision = int(precisionFloat)
	}

	if precision == 0 {
		return lang.NumberValue(math.Round(value)), nil
	}

	multiplier := math.Pow(10, float64(precision))
	return lang.NumberValue(math.Round(value*multiplier) / multiplier), nil
}

func mathTrunc(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("trunc: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("trunc: %w", err)
	}
	return lang.NumberValue(math.Trunc(num)), nil
}

// Power and Root Functions
func mathPow(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("pow: expected 2 arguments")
	}
	base, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("pow: base %w", err)
	}
	exponent, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("pow: exponent %w", err)
	}
	return lang.NumberValue(math.Pow(base, exponent)), nil
}

func mathSqrt(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("sqrt: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("sqrt: %w", err)
	}
	if num < 0 {
		return lang.NumberValue(math.NaN()), nil
	}
	return lang.NumberValue(math.Sqrt(num)), nil
}

func mathCbrt(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("cbrt: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("cbrt: %w", err)
	}
	return lang.NumberValue(math.Cbrt(num)), nil
}

func mathExp(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("exp: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("exp: %w", err)
	}
	return lang.NumberValue(math.Exp(num)), nil
}

func mathExp2(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("exp2: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("exp2: %w", err)
	}
	return lang.NumberValue(math.Exp2(num)), nil
}

// Logarithmic Functions
func mathLog(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("log: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("log: %w", err)
	}
	if num <= 0 {
		return lang.NumberValue(math.NaN()), nil
	}
	return lang.NumberValue(math.Log(num)), nil
}

func mathLog10(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("log10: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("log10: %w", err)
	}
	if num <= 0 {
		return lang.NumberValue(math.NaN()), nil
	}
	return lang.NumberValue(math.Log10(num)), nil
}

func mathLog2(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("log2: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("log2: %w", err)
	}
	if num <= 0 {
		return lang.NumberValue(math.NaN()), nil
	}
	return lang.NumberValue(math.Log2(num)), nil
}

// Trigonometric Functions (angles in radians)
func mathSin(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("sin: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("sin: %w", err)
	}
	return lang.NumberValue(math.Sin(num)), nil
}

func mathCos(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("cos: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("cos: %w", err)
	}
	return lang.NumberValue(math.Cos(num)), nil
}

func mathTan(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("tan: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("tan: %w", err)
	}
	return lang.NumberValue(math.Tan(num)), nil
}

func mathAsin(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("asin: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("asin: %w", err)
	}
	if num < -1 || num > 1 {
		return lang.NumberValue(math.NaN()), nil
	}
	return lang.NumberValue(math.Asin(num)), nil
}

func mathAcos(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("acos: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("acos: %w", err)
	}
	if num < -1 || num > 1 {
		return lang.NumberValue(math.NaN()), nil
	}
	return lang.NumberValue(math.Acos(num)), nil
}

func mathAtan(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("atan: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("atan: %w", err)
	}
	return lang.NumberValue(math.Atan(num)), nil
}

func mathAtan2(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("atan2: expected 2 arguments")
	}
	y, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("atan2: y %w", err)
	}
	x, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("atan2: x %w", err)
	}
	return lang.NumberValue(math.Atan2(y, x)), nil
}

// Hyperbolic Functions
func mathSinh(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("sinh: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("sinh: %w", err)
	}
	return lang.NumberValue(math.Sinh(num)), nil
}

func mathCosh(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("cosh: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("cosh: %w", err)
	}
	return lang.NumberValue(math.Cosh(num)), nil
}

func mathTanh(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("tanh: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("tanh: %w", err)
	}
	return lang.NumberValue(math.Tanh(num)), nil
}

// Angle Conversion
func mathRadians(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("radians: expected 1 argument")
	}
	degrees, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("radians: %w", err)
	}
	return lang.NumberValue(degrees * math.Pi / 180), nil
}

func mathDegrees(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("degrees: expected 1 argument")
	}
	radians, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("degrees: %w", err)
	}
	return lang.NumberValue(radians * 180 / math.Pi), nil
}

// Statistical Functions
func mathSum(args []lang.Value) (lang.Value, error) {
	sum := 0.0
	for i, arg := range args {
		if list, ok := arg.(lang.ListValue); ok {
			for j, item := range list {
				val, err := lib.ToNumber(item)
				if err != nil {
					return nil, fmt.Errorf("sum: list argument %d item %d %w", i, j, err)
				}
				sum += val
			}
		} else {
			val, err := lib.ToNumber(arg)
			if err != nil {
				return nil, fmt.Errorf("sum: argument %d %w", i, err)
			}
			sum += val
		}
	}
	return lang.NumberValue(sum), nil
}

func mathMean(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		return nil, errors.New("mean: expected at least 1 argument")
	}

	count := 0
	sum := 0.0

	for i, arg := range args {
		if list, ok := arg.(lang.ListValue); ok {
			for j, item := range list {
				val, err := lib.ToNumber(item)
				if err != nil {
					return nil, fmt.Errorf("mean: list argument %d item %d %w", i, j, err)
				}
				sum += val
				count++
			}
		} else {
			val, err := lib.ToNumber(arg)
			if err != nil {
				return nil, fmt.Errorf("mean: argument %d %w", i, err)
			}
			sum += val
			count++
		}
	}

	if count == 0 {
		return nil, errors.New("mean: no numeric values found")
	}
	return lang.NumberValue(sum / float64(count)), nil
}

func mathMedian(args []lang.Value) (lang.Value, error) {
	var values []float64

	for i, arg := range args {
		if list, ok := arg.(lang.ListValue); ok {
			for j, item := range list {
				val, err := lib.ToNumber(item)
				if err != nil {
					return nil, fmt.Errorf("median: list argument %d item %d %w", i, j, err)
				}
				values = append(values, val)
			}
		} else {
			val, err := lib.ToNumber(arg)
			if err != nil {
				return nil, fmt.Errorf("median: argument %d %w", i, err)
			}
			values = append(values, val)
		}
	}

	if len(values) == 0 {
		return nil, errors.New("median: no numeric values found")
	}

	sort.Float64s(values)
	n := len(values)

	if n%2 == 0 {
		return lang.NumberValue((values[n/2-1] + values[n/2]) / 2), nil
	}
	return lang.NumberValue(values[n/2]), nil
}

func mathMode(args []lang.Value) (lang.Value, error) {
	frequency := make(map[float64]int)

	for i, arg := range args {
		if list, ok := arg.(lang.ListValue); ok {
			for j, item := range list {
				val, err := lib.ToNumber(item)
				if err != nil {
					return nil, fmt.Errorf("mode: list argument %d item %d %w", i, j, err)
				}
				frequency[val]++
			}
		} else {
			val, err := lib.ToNumber(arg)
			if err != nil {
				return nil, fmt.Errorf("mode: argument %d %w", i, err)
			}
			frequency[val]++
		}
	}

	if len(frequency) == 0 {
		return nil, errors.New("mode: no numeric values found")
	}

	var mode float64
	maxFreq := 0

	for val, freq := range frequency {
		if freq > maxFreq {
			maxFreq = freq
			mode = val
		}
	}

	return lang.NumberValue(mode), nil
}

func mathVariance(args []lang.Value) (lang.Value, error) {
	meanVal, err := mathMean(args)
	if err != nil {
		return nil, fmt.Errorf("variance: %w", err)
	}
	mean, _ := lib.ToNumber(meanVal)

	var values []float64
	for i, arg := range args {
		if list, ok := arg.(lang.ListValue); ok {
			for j, item := range list {
				val, err := lib.ToNumber(item)
				if err != nil {
					return nil, fmt.Errorf("variance: list argument %d item %d %w", i, j, err)
				}
				values = append(values, val)
			}
		} else {
			val, err := lib.ToNumber(arg)
			if err != nil {
				return nil, fmt.Errorf("variance: argument %d %w", i, err)
			}
			values = append(values, val)
		}
	}

	if len(values) <= 1 {
		return lang.NumberValue(0), nil
	}

	sumSquaredDiff := 0.0
	for _, val := range values {
		diff := val - mean
		sumSquaredDiff += diff * diff
	}

	return lang.NumberValue(sumSquaredDiff / float64(len(values)-1)), nil
}

func mathStdDev(args []lang.Value) (lang.Value, error) {
	varianceVal, err := mathVariance(args)
	if err != nil {
		return nil, fmt.Errorf("stddev: %w", err)
	}
	variance, _ := lib.ToNumber(varianceVal)
	return lang.NumberValue(math.Sqrt(variance)), nil
}

// Random Functions
func mathRandom(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		return lang.NumberValue(rand.Float64()), nil
	}
	if len(args) == 1 {
		maxVal, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("random: max %w", err)
		}
		max := int(maxVal)
		if max <= 0 {
			return lang.NumberValue(0), nil
		}
		return lang.NumberValue(float64(rand.Intn(max))), nil
	}
	if len(args) == 2 {
		minVal, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("random: min %w", err)
		}
		maxVal, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("random: max %w", err)
		}
		min := int(minVal)
		max := int(maxVal)
		if max <= min {
			return lang.NumberValue(float64(min)), nil
		}
		return lang.NumberValue(float64(rand.Intn(max-min) + min)), nil
	}
	return nil, errors.New("random: expected 0, 1, or 2 arguments")
}

func mathRandomSeed(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		rand.Seed(time.Now().UnixNano())
	} else if len(args) == 1 {
		seedVal, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("random_seed: %w", err)
		}
		seed := int64(seedVal)
		rand.Seed(seed)
	} else {
		return nil, errors.New("random_seed: expected 0 or 1 argument")
	}
	return lang.BoolValue(true), nil
}

func mathRandomFloat(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		return lang.NumberValue(rand.Float64()), nil
	}
	if len(args) == 1 {
		maxVal, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("random_float: max %w", err)
		}
		return lang.NumberValue(rand.Float64() * maxVal), nil
	}
	if len(args) == 2 {
		minVal, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("random_float: min %w", err)
		}
		maxVal, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("random_float: max %w", err)
		}
		return lang.NumberValue(minVal + rand.Float64()*(maxVal-minVal)), nil
	}
	return nil, errors.New("random_float: expected 0, 1, or 2 arguments")
}

// Utility Functions
func mathIsNaN(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_nan: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_nan: %w", err)
	}
	return lang.BoolValue(math.IsNaN(num)), nil
}

func mathIsInf(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_inf: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_inf: %w", err)
	}
	return lang.BoolValue(math.IsInf(num, 0)), nil
}

func mathIsFinite(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_finite: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_finite: %w", err)
	}
	return lang.BoolValue(!math.IsNaN(num) && !math.IsInf(num, 0)), nil
}

func mathGCD(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 {
		return nil, errors.New("gcd: expected at least 2 arguments")
	}

	firstVal, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("gcd: argument 0 %w", err)
	}
	result := int(math.Abs(firstVal))

	for i := 1; i < len(args); i++ {
		val, err := lib.ToNumber(args[i])
		if err != nil {
			return nil, fmt.Errorf("gcd: argument %d %w", i, err)
		}
		b := int(math.Abs(val))
		result = gcd(result, b)
	}

	return lang.NumberValue(float64(result)), nil
}

func mathLCM(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 {
		return nil, errors.New("lcm: expected at least 2 arguments")
	}

	firstVal, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("lcm: argument 0 %w", err)
	}
	result := int(math.Abs(firstVal))

	for i := 1; i < len(args); i++ {
		val, err := lib.ToNumber(args[i])
		if err != nil {
			return nil, fmt.Errorf("lcm: argument %d %w", i, err)
		}
		b := int(math.Abs(val))
		result = lcm(result, b)
	}

	return lang.NumberValue(float64(result)), nil
}

func mathFactorial(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("factorial: expected 1 argument")
	}

	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("factorial: %w", err)
	}
	n := int(num)

	if n < 0 {
		return lang.NumberValue(math.NaN()), nil
	}
	if n == 0 || n == 1 {
		return lang.NumberValue(1), nil
	}

	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}

	return lang.NumberValue(float64(result)), nil
}

// Helper functions
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return a * b / gcd(a, b)
}

// Math Constants
func mathPi(args []lang.Value) (lang.Value, error) {
	if len(args) != 0 {
		return nil, errors.New("pi: expected 0 arguments")
	}
	return lang.NumberValue(math.Pi), nil
}

func mathE(args []lang.Value) (lang.Value, error) {
	if len(args) != 0 {
		return nil, errors.New("e: expected 0 arguments")
	}
	return lang.NumberValue(math.E), nil
}

func mathPhi(args []lang.Value) (lang.Value, error) {
	if len(args) != 0 {
		return nil, errors.New("phi: expected 0 arguments")
	}
	return lang.NumberValue(math.Phi), nil
}

// Functions that would be in the BuiltinFunctions map:
var MathFunctions = map[string]func([]lang.Value) (lang.Value, error){
	// Basic operations
	"abs":   mathAbs,
	"sign":  mathSign,
	"max":   mathMax,
	"min":   mathMin,
	"clamp": mathClamp,

	// Rounding
	"ceil":  mathCeil,
	"floor": mathFloor,
	"round": mathRound,
	"trunc": mathTrunc,

	// Power and roots
	"pow":  mathPow,
	"sqrt": mathSqrt,
	"cbrt": mathCbrt,
	"exp":  mathExp,
	"exp2": mathExp2,

	// Logarithms
	"log":   mathLog,
	"log10": mathLog10,
	"log2":  mathLog2,

	// Trigonometry
	"sin":   mathSin,
	"cos":   mathCos,
	"tan":   mathTan,
	"asin":  mathAsin,
	"acos":  mathAcos,
	"atan":  mathAtan,
	"atan2": mathAtan2,

	// Hyperbolic
	"sinh": mathSinh,
	"cosh": mathCosh,
	"tanh": mathTanh,

	// Angle conversion
	"radians": mathRadians,
	"degrees": mathDegrees,

	// Statistics
	"sum":      mathSum,
	"mean":     mathMean,
	"avg":      mathMean, // Alias
	"median":   mathMedian,
	"mode":     mathMode,
	"variance": mathVariance,
	"stddev":   mathStdDev,

	// Random
	"random":       mathRandom,
	"random_seed":  mathRandomSeed,
	"random_float": mathRandomFloat,

	// Utilities
	"is_nan":    mathIsNaN,
	"is_inf":    mathIsInf,
	"is_finite": mathIsFinite,
	"gcd":       mathGCD,
	"lcm":       mathLCM,
	"factorial": mathFactorial,

	// Constants
	"pi":  mathPi,
	"e":   mathE,
	"phi": mathPhi,
}
