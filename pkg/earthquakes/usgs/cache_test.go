// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

package usgs

import (
	"testing"

	"github.com/navibyte/quake/internal/geolib"

	pb "github.com/navibyte/quake/api/v1"
)

func TestCacheGet(t *testing.T) {
	testCacheGet(t, pb.Magnitude_MAGNITUDE_M45_PLUS,
		pb.Past_PAST_DAY)
}

func testCacheGet(t *testing.T, magnitude pb.Magnitude, past pb.Past) {
	// get data from a cache (that fetches data from USGS web service if needed)
	col1, err := cacheGetList(magnitude, past)
	if err != nil {
		t.Fatal(err)
	}
	stat1 := cacheGetStat(magnitude, past)
	if stat1.fetchCount != 1 || stat1.hitCount != 0 {
		t.Error("invalid cache fetch or hit count")
	}

	// log some information from parsed data
	bounds := col1.Bounds
	if bounds == nil {
		t.Error("bounds is nil")
	}
	t.Logf("bounds min-lat: %.6f min-lon: %.6f max-lat: %.6f max-lon: %.6f",
		geolib.LatFromE7(bounds.MinLatitude), geolib.LonFromE7(bounds.MinLongitude),
		geolib.LatFromE7(bounds.MaxLatitude), geolib.LonFromE7(bounds.MaxLongitude))
	for _, eq := range col1.Features {
		t.Logf("magnitude %.1f near %s", eq.Magnitude, eq.Place)
	}

	// get data again, this should come from a cache
	col2, err := cacheGetList(magnitude, past)
	if err != nil {
		t.Fatal(err)
	}
	stat2 := cacheGetStat(magnitude, past)
	if stat2.fetchCount != 1 || stat2.hitCount != 1 {
		t.Error("invalid cache fetch or hit count")
	}
	if col1.Metadata.GeneratedTime != col2.Metadata.GeneratedTime {
		t.Error("did not cache fetches properly")
	}

}
