// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

package main

import (
	"fmt"
	"time"

	pb "github.com/navibyte/quake/api/v1"
)

// mockEarthquakeCollection returns a dummy collection only for dev testing
func mockEarthquakeCollection(count int, details bool) *pb.EarthquakeCollection {
	var list []*pb.Earthquake
	for i := 0; i < count; i++ {
		list = append(list, mockEarthquake(fmt.Sprintf("Test%d", i), details))
	}
	col := &pb.EarthquakeCollection{
		Metadata: &pb.EarthquakeMetadata{
			GeneratedTime: time.Now().Unix(),
			Url:           "http://example.org/",
			Title:         "USGS Magnitude 1+ Earthquakes, Past Day",
			Api:           "1.0",
			Count:         15,
			HttpStatus:    "200",
		},
		Bounds: &pb.GeoBoundsE7{
			MinLatitude:  -1000000,
			MinLongitude: -1000000,
			MinHeight:    -5000000,
			MaxLatitude:  1000000,
			MaxLongitude: 1000000,
			MaxHeight:    -3000000,
		},
		Features: list,
	}
	return col
}

// mockEarthquake returns a dummy earthquake only for dev testing
func mockEarthquake(id string, details bool) *pb.Earthquake {
	eq := &pb.Earthquake{
		Id: id,
		Position: &pb.GeoPointE7{
			Latitude:  0,
			Longitude: 0,
			Height:    0,
		},
		Magnitude:      5.0,
		Place:          "No Where",
		Time:           time.Now().Unix(),
		UpdatedTime:    time.Now().Unix(),
		TimezoneOffset: 0,
		Alert:          pb.Alert_ALERT_ORANGE,
		Significance:   5,
	}
	if details {
		eq.Details = mockEarthquakeDetails(id)
	}
	return eq
}

func mockEarthquakeDetails(id string) *pb.EarthquakeDetails {
	return &pb.EarthquakeDetails{
		Id:                 id,
		Url:                "http://example.org/earthquake/" + id,
		DetailFeedUrl:      "http://example.org/earthquake/" + id + "/details",
		Felt:               50,
		ReportedIntensity:  4.9,
		EstimatedIntensity: 5.2,
		Status:             pb.Status_STATUS_REVIEWED,
		Tsunami:            false,
		Network:            "us",
		Code:               "2013lgaz",
		Ids:                ",ci15296281,us2013mqbd,at00mji9pf,",
		Sources:            ",us,nc,ci,",
		ProductTypes:       ",cap,dyfi,general-link,origin,p-wave-travel-times,phase-data,",
		Nst:                10,
		Dmin:               4.0,
		Rms:                0.5,
		Gap:                135.0,
		MagType:            "Md",
		Type:               pb.Type_TYPE_EARTHQUAKE,
	}
}
