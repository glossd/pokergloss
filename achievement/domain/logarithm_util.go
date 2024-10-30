package domain

import "math"

func Logarithm(of, base int64) float64 {
	return math.Log10(float64(of)) / math.Log10(float64(base))
}
