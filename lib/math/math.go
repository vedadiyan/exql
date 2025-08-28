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

func Abs() (string, lang.Function) {
	name := "abs"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Abs(num)), nil
	}
	return name, fn
}

func Sign() (string, lang.Function) {
	name := "sign"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if num > 0 {
			return lang.NumberValue(1), nil
		} else if num < 0 {
			return lang.NumberValue(-1), nil
		}
		return lang.NumberValue(0), nil
	}
	return name, fn
}

func Max() (string, lang.Function) {
	name := "max"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) == 0 {
			return nil, lib.ArgumentErrorMin(name, 1)
		}
		max, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: argument 0 %w", name, err)
		}
		for i := 1; i < len(args); i++ {
			val, err := lib.ToNumber(args[i])
			if err != nil {
				return nil, fmt.Errorf("%s: argument %d %w", name, i, err)
			}
			if val > max {
				max = val
			}
		}
		return lang.NumberValue(max), nil
	}
	return name, fn
}

func Min() (string, lang.Function) {
	name := "min"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) == 0 {
			return nil, lib.ArgumentErrorMin(name, 1)
		}
		min, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: argument 0 %w", name, err)
		}
		for i := 1; i < len(args); i++ {
			val, err := lib.ToNumber(args[i])
			if err != nil {
				return nil, fmt.Errorf("%s: argument %d %w", name, i, err)
			}
			if val < min {
				min = val
			}
		}
		return lang.NumberValue(min), nil
	}
	return name, fn
}

func Clamp() (string, lang.Function) {
	name := "clamp"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, lib.ArgumentError(name, 3)
		}
		value, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: value %w", name, err)
		}
		min, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: min %w", name, err)
		}
		max, err := lib.ToNumber(args[2])
		if err != nil {
			return nil, fmt.Errorf("%s: max %w", name, err)
		}
		if value < min {
			return lang.NumberValue(min), nil
		}
		if value > max {
			return lang.NumberValue(max), nil
		}
		return lang.NumberValue(value), nil
	}
	return name, fn
}

func Ceil() (string, lang.Function) {
	name := "ceil"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Ceil(num)), nil
	}
	return name, fn
}

func Floor() (string, lang.Function) {
	name := "floor"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Floor(num)), nil
	}
	return name, fn
}

func Round() (string, lang.Function) {
	name := "round"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, lib.ArgumentErrorRange(name, 1, 2)
		}
		value, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: value %w", name, err)
		}
		precision := 0
		if len(args) == 2 {
			precisionFloat, err := lib.ToNumber(args[1])
			if err != nil {
				return nil, fmt.Errorf("%s: precision %w", name, err)
			}
			precision = int(precisionFloat)
		}
		if precision == 0 {
			return lang.NumberValue(math.Round(value)), nil
		}
		multiplier := math.Pow(10, float64(precision))
		return lang.NumberValue(math.Round(value*multiplier) / multiplier), nil
	}
	return name, fn
}

func Trunc() (string, lang.Function) {
	name := "trunc"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Trunc(num)), nil
	}
	return name, fn
}

func Pow() (string, lang.Function) {
	name := "pow"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		base, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: base %w", name, err)
		}
		exponent, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: exponent %w", name, err)
		}
		return lang.NumberValue(math.Pow(base, exponent)), nil
	}
	return name, fn
}

func Sqrt() (string, lang.Function) {
	name := "sqrt"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if num < 0 {
			return lang.NumberValue(math.NaN()), nil
		}
		return lang.NumberValue(math.Sqrt(num)), nil
	}
	return name, fn
}

func Cbrt() (string, lang.Function) {
	name := "cbrt"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Cbrt(num)), nil
	}
	return name, fn
}

func Exp() (string, lang.Function) {
	name := "exp"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Exp(num)), nil
	}
	return name, fn
}

func Exp2() (string, lang.Function) {
	name := "exp2"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Exp2(num)), nil
	}
	return name, fn
}

func Log() (string, lang.Function) {
	name := "log"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if num <= 0 {
			return lang.NumberValue(math.NaN()), nil
		}
		return lang.NumberValue(math.Log(num)), nil
	}
	return name, fn
}

func Log10() (string, lang.Function) {
	name := "log10"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if num <= 0 {
			return lang.NumberValue(math.NaN()), nil
		}
		return lang.NumberValue(math.Log10(num)), nil
	}
	return name, fn
}

func Log2() (string, lang.Function) {
	name := "log2"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if num <= 0 {
			return lang.NumberValue(math.NaN()), nil
		}
		return lang.NumberValue(math.Log2(num)), nil
	}
	return name, fn
}

func Sin() (string, lang.Function) {
	name := "sin"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Sin(num)), nil
	}
	return name, fn
}

func Cos() (string, lang.Function) {
	name := "cos"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Cos(num)), nil
	}
	return name, fn
}

func Tan() (string, lang.Function) {
	name := "tan"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Tan(num)), nil
	}
	return name, fn
}

func Asin() (string, lang.Function) {
	name := "asin"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if num < -1 || num > 1 {
			return lang.NumberValue(math.NaN()), nil
		}
		return lang.NumberValue(math.Asin(num)), nil
	}
	return name, fn
}

func Acos() (string, lang.Function) {
	name := "acos"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if num < -1 || num > 1 {
			return lang.NumberValue(math.NaN()), nil
		}
		return lang.NumberValue(math.Acos(num)), nil
	}
	return name, fn
}

func Atan() (string, lang.Function) {
	name := "atan"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Atan(num)), nil
	}
	return name, fn
}

func Atan2() (string, lang.Function) {
	name := "atan2"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		y, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: y %w", name, err)
		}
		x, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: x %w", name, err)
		}
		return lang.NumberValue(math.Atan2(y, x)), nil
	}
	return name, fn
}

func Sinh() (string, lang.Function) {
	name := "sinh"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Sinh(num)), nil
	}
	return name, fn
}

func Cosh() (string, lang.Function) {
	name := "cosh"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Cosh(num)), nil
	}
	return name, fn
}

func Tanh() (string, lang.Function) {
	name := "tanh"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(math.Tanh(num)), nil
	}
	return name, fn
}

func Radians() (string, lang.Function) {
	name := "radians"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		degrees, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(degrees * math.Pi / 180), nil
	}
	return name, fn
}

func Degrees() (string, lang.Function) {
	name := "degrees"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		radians, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(radians * 180 / math.Pi), nil
	}
	return name, fn
}

func Sum() (string, lang.Function) {
	name := "sum"
	fn := func(args []lang.Value) (lang.Value, error) {
		sum := 0.0
		for i, arg := range args {
			if list, ok := arg.(lang.ListValue); ok {
				for j, item := range list {
					val, err := lib.ToNumber(item)
					if err != nil {
						return nil, fmt.Errorf("%s: list argument %d item %d %w", name, i, j, err)
					}
					sum += val
				}
			} else {
				val, err := lib.ToNumber(arg)
				if err != nil {
					return nil, fmt.Errorf("%s: argument %d %w", name, i, err)
				}
				sum += val
			}
		}
		return lang.NumberValue(sum), nil
	}
	return name, fn
}

func Mean() (string, lang.Function) {
	name := "mean"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) == 0 {
			return nil, lib.ArgumentErrorMin(name, 1)
		}
		count := 0
		sum := 0.0
		for i, arg := range args {
			if list, ok := arg.(lang.ListValue); ok {
				for j, item := range list {
					val, err := lib.ToNumber(item)
					if err != nil {
						return nil, fmt.Errorf("%s: list argument %d item %d %w", name, i, j, err)
					}
					sum += val
					count++
				}
			} else {
				val, err := lib.ToNumber(arg)
				if err != nil {
					return nil, fmt.Errorf("%s: argument %d %w", name, i, err)
				}
				sum += val
				count++
			}
		}
		if count == 0 {
			return nil, fmt.Errorf("%s: no numeric values found", name)
		}
		return lang.NumberValue(sum / float64(count)), nil
	}
	return name, fn
}

func Median() (string, lang.Function) {
	name := "median"
	fn := func(args []lang.Value) (lang.Value, error) {
		var values []float64
		for i, arg := range args {
			if list, ok := arg.(lang.ListValue); ok {
				for j, item := range list {
					val, err := lib.ToNumber(item)
					if err != nil {
						return nil, fmt.Errorf("%s: list argument %d item %d %w", name, i, j, err)
					}
					values = append(values, val)
				}
			} else {
				val, err := lib.ToNumber(arg)
				if err != nil {
					return nil, fmt.Errorf("%s: argument %d %w", name, i, err)
				}
				values = append(values, val)
			}
		}
		if len(values) == 0 {
			return nil, fmt.Errorf("%s: no numeric values found", name)
		}
		sort.Float64s(values)
		n := len(values)
		if n%2 == 0 {
			return lang.NumberValue((values[n/2-1] + values[n/2]) / 2), nil
		}
		return lang.NumberValue(values[n/2]), nil
	}
	return name, fn
}

func Mode() (string, lang.Function) {
	name := "mode"
	fn := func(args []lang.Value) (lang.Value, error) {
		frequency := make(map[float64]int)
		for i, arg := range args {
			if list, ok := arg.(lang.ListValue); ok {
				for j, item := range list {
					val, err := lib.ToNumber(item)
					if err != nil {
						return nil, fmt.Errorf("%s: list argument %d item %d %w", name, i, j, err)
					}
					frequency[val]++
				}
			} else {
				val, err := lib.ToNumber(arg)
				if err != nil {
					return nil, fmt.Errorf("%s: argument %d %w", name, i, err)
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
	return name, fn
}

func Variance() (string, lang.Function) {
	name := "variance"
	fn := func(args []lang.Value) (lang.Value, error) {
		_, meanFunc := Mean()
		meanVal, err := meanFunc(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		mean, _ := lib.ToNumber(meanVal)
		var values []float64
		for i, arg := range args {
			if list, ok := arg.(lang.ListValue); ok {
				for j, item := range list {
					val, err := lib.ToNumber(item)
					if err != nil {
						return nil, fmt.Errorf("%s: list argument %d item %d %w", name, i, j, err)
					}
					values = append(values, val)
				}
			} else {
				val, err := lib.ToNumber(arg)
				if err != nil {
					return nil, fmt.Errorf("%s: argument %d %w", name, i, err)
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
	return name, fn
}

func StdDev() (string, lang.Function) {
	name := "stddev"
	fn := func(args []lang.Value) (lang.Value, error) {
		_, varianceFunc := Variance()
		varianceVal, err := varianceFunc(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		variance, _ := lib.ToNumber(varianceVal)
		return lang.NumberValue(math.Sqrt(variance)), nil
	}
	return name, fn
}

func Random() (string, lang.Function) {
	name := "random"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) == 0 {
			return lang.NumberValue(rand.Float64()), nil
		}
		if len(args) == 1 {
			maxVal, err := lib.ToNumber(args[0])
			if err != nil {
				return nil, fmt.Errorf("%s: max %w", name, err)
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
				return nil, fmt.Errorf("%s: min %w", name, err)
			}
			maxVal, err := lib.ToNumber(args[1])
			if err != nil {
				return nil, fmt.Errorf("%s: max %w", name, err)
			}
			min := int(minVal)
			max := int(maxVal)
			if max <= min {
				return lang.NumberValue(float64(min)), nil
			}
			return lang.NumberValue(float64(rand.Intn(max-min) + min)), nil
		}
		return nil, lib.ArgumentErrorMultiRange(name, []int{0, 1, 2})
	}
	return name, fn
}

func RandomSeed() (string, lang.Function) {
	name := "random_seed"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) == 0 {
			rand.Seed(time.Now().UnixNano())
		} else if len(args) == 1 {
			seedVal, err := lib.ToNumber(args[0])
			if err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
			seed := int64(seedVal)
			rand.Seed(seed)
		} else {
			return nil, lib.ArgumentErrorRange(name, 0, 1)
		}
		return lang.BoolValue(true), nil
	}
	return name, fn
}

func RandomFloat() (string, lang.Function) {
	name := "random_float"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) == 0 {
			return lang.NumberValue(rand.Float64()), nil
		}
		if len(args) == 1 {
			maxVal, err := lib.ToNumber(args[0])
			if err != nil {
				return nil, fmt.Errorf("%s: max %w", name, err)
			}
			return lang.NumberValue(rand.Float64() * maxVal), nil
		}
		if len(args) == 2 {
			minVal, err := lib.ToNumber(args[0])
			if err != nil {
				return nil, fmt.Errorf("%s: min %w", name, err)
			}
			maxVal, err := lib.ToNumber(args[1])
			if err != nil {
				return nil, fmt.Errorf("%s: max %w", name, err)
			}
			return lang.NumberValue(minVal + rand.Float64()*(maxVal-minVal)), nil
		}
		return nil, lib.ArgumentErrorMultiRange(name, []int{0, 1, 2})
	}
	return name, fn
}

func IsNaN() (string, lang.Function) {
	name := "is_nan"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.BoolValue(math.IsNaN(num)), nil
	}
	return name, fn
}

func IsInf() (string, lang.Function) {
	name := "is_inf"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.BoolValue(math.IsInf(num, 0)), nil
	}
	return name, fn
}

func IsFinite() (string, lang.Function) {
	name := "is_finite"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.BoolValue(!math.IsNaN(num) && !math.IsInf(num, 0)), nil
	}
	return name, fn
}

func GCD() (string, lang.Function) {
	name := "gcd"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 {
			return nil, lib.ArgumentErrorMin(name, 2)
		}
		firstVal, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: argument 0 %w", name, err)
		}
		result := int(math.Abs(firstVal))
		for i := 1; i < len(args); i++ {
			val, err := lib.ToNumber(args[i])
			if err != nil {
				return nil, fmt.Errorf("%s: argument %d %w", name, i, err)
			}
			b := int(math.Abs(val))
			result = gcd(result, b)
		}
		return lang.NumberValue(float64(result)), nil
	}
	return name, fn
}

func LCM() (string, lang.Function) {
	name := "lcm"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 {
			return nil, lib.ArgumentErrorMin(name, 2)
		}
		firstVal, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: argument 0 %w", name, err)
		}
		result := int(math.Abs(firstVal))
		for i := 1; i < len(args); i++ {
			val, err := lib.ToNumber(args[i])
			if err != nil {
				return nil, fmt.Errorf("%s: argument %d %w", name, i, err)
			}
			b := int(math.Abs(val))
			result = lcm(result, b)
		}
		return lang.NumberValue(float64(result)), nil
	}
	return name, fn
}

func Factorial() (string, lang.Function) {
	name := "factorial"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
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
	return name, fn
}

func Pi() (string, lang.Function) {
	name := "pi"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 0 {
			return nil, lib.ArgumentError(name, 0)
		}
		return lang.NumberValue(math.Pi), nil
	}
	return name, fn
}

func E() (string, lang.Function) {
	name := "e"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 0 {
			return nil, lib.ArgumentError(name, 0)
		}
		return lang.NumberValue(math.E), nil
	}
	return name, fn
}

func Phi() (string, lang.Function) {
	name := "phi"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 0 {
			return nil, lib.ArgumentError(name, 0)
		}
		return lang.NumberValue(math.Phi), nil
	}
	return name, fn
}

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

var MathFunctions = []func() (string, lang.Function){
	Abs,
	Sign,
	Max,
	Min,
	Clamp,
	Ceil,
	Floor,
	Round,
	Trunc,
	Pow,
	Sqrt,
	Cbrt,
	Exp,
	Exp2,
	Log,
	Log10,
	Log2,
	Sin,
	Cos,
	Tan,
	Asin,
	Acos,
	Atan,
	Atan2,
	Sinh,
	Cosh,
	Tanh,
	Radians,
	Degrees,
	Sum,
	Mean,
	Median,
	Mode,
	Variance,
	StdDev,
	Random,
	RandomSeed,
	RandomFloat,
	IsNaN,
	IsInf,
	IsFinite,
	GCD,
	LCM,
	Factorial,
	Pi,
	E,
	Phi,
}
