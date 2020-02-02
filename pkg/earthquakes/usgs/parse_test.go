// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

package usgs

import (
	"io/ioutil"
	"testing"

	pb "github.com/navibyte/quake/api/v1"
)

func TestParsingGeoJSON(t *testing.T) {
	// for parsing we use a locally (on dev environment) stored file
	b, err := ioutil.ReadFile("testdata/4.5_day.json")
	if err != nil {
		t.Fatal(err)
	}

	// ensure we can parse test GeoJSON data to earthquakes
	col, err := ToEarthquakeCollection(b, true)
	if err != nil {
		t.Fatal(err)
	}
	if col.Metadata.GeneratedTime != 1577968050000/1000 {
		t.Error("invalid generated time")
	}
	//fmt.Println("Generated time: ", ts.Format(time.UnixDate))
	bounds := col.Bounds
	if bounds == nil {
		t.Error("bounds is nil")
	}
	if bounds.MinLatitude != -53_0705000 || bounds.MinLongitude != -178_2268000 || bounds.MinHeight != -268_710_00 ||
		bounds.MaxLatitude != 55_5519000 || bounds.MaxLongitude != 170_3780000 || bounds.MaxHeight != -7_920_00 {
		t.Error("invalid bounds")
	}
	if len(col.Features) != 21 {
		t.Error("invalid feature count")
	} else {
		eq := col.Features[3]
		if eq.Id != "us70006tf3" {
			t.Error("invalid id")
		}
		if eq.Place != "218km NW of Saumlaki, Indonesia" {
			t.Error("invalid properties on a feature")
		}
		if eq.Alert !=  pb.Alert_ALERT_UNSPECIFIED {
			t.Error("invalid alert")
		}
		if eq.Position == nil {
			t.Error("no position")
		} else if eq.Position.Longitude != 130_1236000 {
			t.Error("invalid longitude")
		}
		details := eq.Details
		if details == nil {
			t.Error("no details found")
		} else {
			if details.Status != pb.Status_STATUS_REVIEWED {
				t.Error("invalid status")
			}
			if details.Ids != ",us70006tf3," {
				t.Error("invalid id")
			}
		}
	}
}
