// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

package usgs

import (
	"testing"

	"github.com/navibyte/quake/internal/geolib"
	pb "github.com/navibyte/quake/api/v1"
)

func TestRepository(t *testing.T) {
	testRepository(t, pb.Magnitude_MAGNITUDE_M25_PLUS,
		pb.Past_PAST_DAY, 5, nil)
	testRepository(t, pb.Magnitude_MAGNITUDE_M25_PLUS,
		pb.Past_PAST_DAY, 5,
		&pb.GeoBoundsE7{
			MinLatitude:  10_0000000,
			MinLongitude: -70_0000000,
			MinHeight:    -100000,
			MaxLatitude:  20_0000000,
			MaxLongitude: -60_0000000,
			MaxHeight:    0,
		})
}

func testRepository(t *testing.T, magnitude pb.Magnitude, past pb.Past,
	limit int, bounds *pb.GeoBoundsE7) {

	var col *pb.EarthquakeCollection
	var err error

	// list data from repository (true => ask for details)
	if bounds == nil {
		col, err = ListEarthquakes(magnitude, past, limit, true)
	} else {
		col, err = ListEarthquakesFocusBounds(magnitude, past, limit, true, bounds)
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(col.Features) > limit {
		t.Error("too many features")
	}
	for _, eq := range col.Features {
		if eq.Details == nil {
			t.Error("has no details even if asked")
		}
		t.Logf("magnitude %.1f %.2f %.2f", eq.Magnitude,
			geolib.LatFromE7(eq.Position.Latitude),
			geolib.LonFromE7(eq.Position.Longitude))
	}

	// list data again from repository (false => without details)
	if bounds == nil {
		col, err = ListEarthquakes(magnitude, past, limit, false)
	} else {
		col, err = ListEarthquakesFocusBounds(magnitude, past, limit, false, bounds)
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(col.Features) > limit {
		t.Error("too many features")
	}
	for _, eq := range col.Features {
		if eq.Details != nil {
			t.Error("has details even if asked not")
		}
	}

}
