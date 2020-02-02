// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

package usgs

import (
	"sort"

	pb "github.com/navibyte/quake/api/v1"
	"github.com/navibyte/quake/internal/geolib"
	"github.com/navibyte/quake/internal/mathlib"
)

func GetEarthquake(id string) (*pb.Earthquake, error) {
	return cacheGetById(id)
}

func ListEarthquakes(magnitude pb.Magnitude, past pb.Past,
	limit int, details bool) (*pb.EarthquakeCollection, error) {

	// get collection from the cache
	col, err := cacheGetList(magnitude, past)
	if err != nil {
		return nil, err
	}

	// return collection "as-is" if details was asked and no too many features
	noLimit := limit <= 0
	if details && (noLimit || len(col.Features) <= limit) {
		return col, nil
	}

	// still here, we have to filter before returning feature collection
	result := copyCollection(col, limit, details, nil, nil)
	return result, nil
}

func ListEarthquakesFocusPosition(magnitude pb.Magnitude, past pb.Past,
	limit int, details bool, pos *pb.GeoPointE7) (*pb.EarthquakeCollection, error) {

	// get collection from the cache
	col, err := cacheGetList(magnitude, past)
	if err != nil {
		return nil, err
	}

	// filter resulting collection (and sort it by focusing on a position)
	result := copyCollection(col, limit, details, pos, nil)
	return result, nil
}

func ListEarthquakesFocusBounds(magnitude pb.Magnitude, past pb.Past,
	limit int, details bool, bounds *pb.GeoBoundsE7) (*pb.EarthquakeCollection, error) {

	// get collection from the cache
	col, err := cacheGetList(magnitude, past)
	if err != nil {
		return nil, err
	}

	// focus point (for sorting) at mid of bounding box
	pos := &pb.GeoPointE7{
		Latitude:  bounds.MinLatitude + (bounds.MaxLatitude-bounds.MinLatitude)/2,
		Longitude: bounds.MinLongitude + (bounds.MaxLongitude-bounds.MinLongitude)/2,
		Height:    bounds.MinHeight + (bounds.MaxHeight-bounds.MinHeight)/2,
	}

	// filter resulting collection (and sort it by focusing on a position for
	// those earthquakes that locates inside bounds)
	result := copyCollection(col, limit, details, pos, bounds)
	return result, nil
}

func copyCollection(from *pb.EarthquakeCollection, limit int, details bool,
	focus *pb.GeoPointE7, bounds *pb.GeoBoundsE7) *pb.EarthquakeCollection {

	to := &pb.EarthquakeCollection{}
	if m := from.Metadata; m != nil {
		to.Metadata = &pb.EarthquakeMetadata{
			GeneratedTime: m.GeneratedTime,
			Url:           m.Url,
			Title:         m.Title,
			Api:           m.Api,
			Count:         m.Count,
			HttpStatus:    m.HttpStatus,
		}
	}

	list := from.Features
	if focus != nil {
		// need to focus and sort earthquake features
		type sorter struct {
			eq   *pb.Earthquake
			dist float64
		}
		var sorting []sorter
		for _, eq := range list {
			pos := eq.Position
			lat := pos.Latitude
			lon := pos.Longitude
			if bounds != nil {
				if lat < bounds.MinLatitude || lat > bounds.MaxLatitude ||
					lon < bounds.MinLongitude || lon > bounds.MaxLongitude {
					continue // out of bounds, so skip
				}
			}
			dist := geolib.DistanceE7(lat, lon, focus.Latitude, focus.Longitude)
			sorting = append(sorting, sorter{eq: eq, dist: dist})
		}
		sort.Slice(sorting, func(i, j int) bool {
			return sorting[i].dist < sorting[j].dist
		})
		list = make([]*pb.Earthquake, len(sorting))
		for i, s := range sorting {
			list[i] = s.eq
		}
	}

	// append features (up to number of limit) to resulting collection
	for _, eq := range list {
		if limit > 0 && len(to.Features) >= limit {
			break
		}
		if details {
			to.Features = append(to.Features, eq)
		} else {
			to.Features = append(to.Features,
				cloneEarthquakeWithoutDetails(eq))
		}
		if to.Bounds == nil {
			to.Bounds = createBounds(eq.Position)
		} else {
			addToBounds(to.Bounds, eq.Position)
		}
	}
	if to.Metadata != nil {
		to.Metadata.Count = int32(len(to.Features))
	}
	return to
}

func cloneEarthquakeWithoutDetails(eq *pb.Earthquake) *pb.Earthquake {
	return &pb.Earthquake{
		Id:             eq.Id,
		Position:       eq.Position,
		Magnitude:      eq.Magnitude,
		Place:          eq.Place,
		Time:           eq.Time,
		UpdatedTime:    eq.UpdatedTime,
		TimezoneOffset: eq.TimezoneOffset,
		Alert:          eq.Alert,
		Significance:   eq.Significance,
	}
}

func createBounds(pos *pb.GeoPointE7) *pb.GeoBoundsE7 {
	return &pb.GeoBoundsE7{
		MinLatitude:  pos.Latitude,
		MinLongitude: pos.Longitude,
		MinHeight:    pos.Height,
		MaxLatitude:  pos.Latitude,
		MaxLongitude: pos.Longitude,
		MaxHeight:    pos.Height,
	}
}

func addToBounds(bounds *pb.GeoBoundsE7, pos *pb.GeoPointE7) {
	bounds.MinLatitude = mathlib.MinInt32(bounds.MinLatitude, pos.Latitude)
	bounds.MinLongitude = mathlib.MinInt32(bounds.MinLongitude, pos.Longitude)
	bounds.MinHeight = mathlib.MinInt32(bounds.MinHeight, pos.Height)
	bounds.MaxLatitude = mathlib.MaxInt32(bounds.MaxLatitude, pos.Latitude)
	bounds.MaxLongitude = mathlib.MaxInt32(bounds.MaxLongitude, pos.Longitude)
	bounds.MaxHeight = mathlib.MaxInt32(bounds.MaxHeight, pos.Height)
}
