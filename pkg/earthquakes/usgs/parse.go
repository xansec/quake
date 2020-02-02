// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

package usgs

import (
	"errors"

	pb "github.com/navibyte/quake/api/v1"
	"github.com/navibyte/quake/internal/geolib"
	"github.com/navibyte/quake/internal/jsonlib"
	"github.com/navibyte/quake/internal/mathlib"
	"github.com/tidwall/gjson"
)

// ErrInvalidJSON can be returned by the parser
var ErrInvalidJSON = errors.New("invalid JSON data")

// ErrInvalidGeoJSON can be returned by the parser
var ErrInvalidGeoJSON = errors.New("JSON data is not valid GeoJSON")

// ToEarthquakeCollection parses GeoJSON data from USGS to earthquake objects.
func ToEarthquakeCollection(data []byte, details bool) (*pb.EarthquakeCollection, error) {
	// validate JSON data
	if !gjson.ValidBytes(data) {
		return nil, ErrInvalidJSON
	}

	// init response
	var col pb.EarthquakeCollection

	// parse JSON data and check that is GeoJSON containing FeatureCollection
	root := gjson.ParseBytes(data)
	if root.Get("type").String() != "FeatureCollection" {
		return nil, ErrInvalidGeoJSON
	}
	c := jsonlib.NewCursor(root)

	// parse metadata
	if m := c.Get("metadata"); m.IsObject() {
		col.Metadata = &pb.EarthquakeMetadata{
			GeneratedTime: m.Int64("generated") / 1000,
			Url:           m.String("url"),
			Title:         m.String("title"),
			Api:           m.String("api"),
			Count:         m.Int32("meta"),
			HttpStatus:    m.String("status"),
		}
	}

	// Parse (collection) bbox.
	// Note GeoJSON from USGS has bbox:
	//    [min-lon, min-lat, min-depth, max-lon, max-lat, max-depth]
	// GeoBoundsE7 has height (above sea level) instead of depth (of GeoJSON).
	if b := c.Get("bbox"); b.IsArray() {
		col.Bounds = &pb.GeoBoundsE7{
			MinLatitude:  geolib.LatToE7(b.Float64("1")),
			MinLongitude: geolib.LonToE7(b.Float64("0")),
			MinHeight:    depthToHeightCentimeters(b.Float64("5")),
			MaxLatitude:  geolib.LatToE7(b.Float64("4")),
			MaxLongitude: geolib.LonToE7(b.Float64("3")),
			MaxHeight:    depthToHeightCentimeters(b.Float64("2")),
		}
	}

	// parse features (that is earthquakes)
	c.ForEachArray("features", func(feature jsonlib.Cursor) bool {
		if feature.String("type") != "Feature" {
			return true // ignore this element but keep iterating
		}

		// parse an Earthquake from a GeoJSON Feature structure
		var eq pb.Earthquake
		// first parse "id" of a feature
		eqID := feature.String("id")
		eq.Id = eqID
		// parse point geometry
		if feature.String("geometry.type") == "Point" {
			if coord := feature.Get("geometry.coordinates"); coord.IsArray() {
				eq.Position = &pb.GeoPointE7{
					Latitude:  geolib.LatToE7(coord.Float64("1")),
					Longitude: geolib.LonToE7(coord.Float64("0")),
					Height:    depthToHeightCentimeters(coord.Float64("2")),
				}
			}
		}
		// parse also properties
		if prop := feature.Get("properties"); prop.IsObject() {
			// parse core properties
			eq.Magnitude = prop.Float32("mag")
			eq.Place = prop.String("place")
			eq.Time = prop.Int64("time") / 1000
			eq.UpdatedTime = prop.Int64("updated") / 1000
			eq.TimezoneOffset = prop.Int32("tz")
			eq.Alert = parseAlert(prop.String("alert"))
			eq.Significance = prop.Int32("sig")

			// detailed properties are parsed only if needed
			if details {
				eq.Details = &pb.EarthquakeDetails{
					Id:                 eqID,
					Url:                prop.String("url"),
					DetailFeedUrl:      prop.String("detail"),
					Felt:               prop.Int32("felt"),
					ReportedIntensity:  prop.Float32("cdi"),
					EstimatedIntensity: prop.Float32("mmi"),
					Status:             parseStatus(prop.String("status")),
					Tsunami:            prop.Int32("tsunami") == 1,
					Network:            prop.String("net"),
					Code:               prop.String("code"),
					Ids:                prop.String("ids"),
					Sources:            prop.String("sources"),
					ProductTypes:       prop.String("types"),
					Nst:                prop.Int32("nst"),
					Dmin:               prop.Float32("dmin"),
					Rms:                prop.Float32("rms"),
					Gap:                prop.Float32("gap"),
					MagType:            prop.String("magType"),
					Type:               parseType(prop.String("type")),
				}
			}
		}
		// append a new Earthquake to collection
		col.Features = append(col.Features, &eq)

		return true // keep iterating
	})

	return &col, nil
}

func depthToHeightCentimeters(depthKM float64) int32 {
	// Earthquake depth of the event is kilometers, range about [0, 1000].
	// For this application depth is converted to "height above sea, cm"
	return mathlib.Round32(-depthKM * 100000)
}

func parseAlert(value string) pb.Alert {
	switch value {
	case "red":
		return pb.Alert_ALERT_RED
	case "orange":
		return pb.Alert_ALERT_ORANGE
	case "yellow":
		return pb.Alert_ALERT_YELLOW
	case "green":
		return pb.Alert_ALERT_GREEN
	default:
		return pb.Alert_ALERT_UNSPECIFIED
	}
}

func parseStatus(value string) pb.Status {
	switch value {
	case "automatic":
		return pb.Status_STATUS_AUTOMATIC
	case "reviewed":
		return pb.Status_STATUS_REVIEWED
	case "deleted":
		return pb.Status_STATUS_DELETED
	default:
		return pb.Status_STATUS_UNSPECIFIED
	}
}

func parseType(value string) pb.Type {
	switch value {
	case "earthquake":
		return pb.Type_TYPE_EARTHQUAKE
	case "quarry":
		return pb.Type_TYPE_QUARRY
	default:
		return pb.Type_TYPE_UNSPECIFIED
	}
}
