// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

package mathlib

import "math"

// Round32 rounds the float64 value to int32 representation.
func Round32(value float64) int32 {
	if value < 0 {
		return int32(value - 0.5)
	}
	return int32(value + 0.5)
}

// ClipInt32 clips the value to the range [min, max].
func ClipInt32(value int32, min int32, max int32) int32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClipFloat64 clips the value to the range [min, max].
func ClipFloat64(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// MinInt32 returns minimum value of a and b.
func MinInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// MaxInt32 returns maximum value of a and b.
func MaxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// ToRad converts degrees to radians.
func ToRad(value float64) float64 {
	return value * math.Pi / float64(180)
}
