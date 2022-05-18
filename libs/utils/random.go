package utils

import (
	"math/rand"
	"time"
)

// RandomIntn 返回一个[0,n)的随机数，不包含n
func RandomIntn(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

// RandomFloat32 返回一个取值范围在[0.0, 1.0)的伪随机float32值。
func RandomFloat32() float32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float32()
}

// RandomFloat64 返回一个取值范围在[0.0, 1.0)的伪随机float64值。
func RandomFloat64() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}
