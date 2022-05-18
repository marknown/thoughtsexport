package utils

import "strconv"

// Float64FromString 字符串转成 float64
func Float64FromString(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return f
}
