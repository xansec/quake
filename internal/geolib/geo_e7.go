// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

package geolib

import (
 "math"
 "github.com/navibyte/quake/internal/mathlib"
)

const factorE7 = 1e7

// LatFromE7 converts E7 integer representation to latitude.
// The argument lat must be (and clipped) on the range [-90_0000000, 90_0000000. 
func LatFromE7(e7 int32) float64 {
	return float64(mathlib.ClipInt32(e7, -90_0000000, 90_0000000)) / factorE7
}

// LonFromE7 converts E7 integer representation to latitude.
// The argument lon must be (and clipped) on the range [-180_0000000, 180_0000000]. 
func LonFromE7(e7 int32) float64 {
	return float64(mathlib.ClipInt32(e7, -180_0000000, 180_0000000)) / factorE7
}

// LatToE7 converts latitude to E7 integer representation.
// The argument lat must be (and clipped) on the range [-90.0, 90.0]. 
func LatToE7(lat float64) int32 {
	return mathlib.Round32(mathlib.ClipFloat64(lat, -90.0, 90.0) * factorE7)
}

// LonToE7 converts longitude to E7 integer representation.
// The argument lon must be (and clipped) on the range [-180.0, 180.0]. 
func LonToE7(lon float64) int32 {
	return mathlib.Round32(mathlib.ClipFloat64(lon, -180.0, 180.0) * factorE7)
}

// DistanceE7 returns a distance between two points. Result is meters.
func DistanceE7(lat1, lon1, lat2, lon2 int32) float64 {
	return Distance(LatFromE7(lat1), LonFromE7(lon1),
			LatFromE7(lat2), LonFromE7(lon2))
}

// Distance returns a distance between two points. Result is meters.
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// using "haversine" formula
	// see: http://mathforum.org/library/drmath/view/51879.html
	
	const earthRadius = float64(6371000)

	lat1Rad := mathlib.ToRad(lat1)
	lat2Rad := mathlib.ToRad(lat2)
	dlat := mathlib.ToRad(lat2 - lat1)
	dlon := mathlib.ToRad(lon2 - lon1)
 	a := math.Sin(dlat/2) * math.Sin(dlat/2) + 
		 math.Cos(lat1Rad) * math.Cos(lat2Rad) * 
		 math.Sin(dlon/2) * math.Sin(dlon/2)
  	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a)) 
	
	return earthRadius * c
}
