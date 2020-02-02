// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

package usgs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	pb "github.com/navibyte/quake/api/v1"
)

const (
	apiBaseURL        = "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/"
	apiBaseURLPostfix = ".geojson"
)

var (
	httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
)

// ErrUnknownDataRequest is returned by a parser when could not formulate a request
var ErrUnknownDataRequest = errors.New("unknown earthquake data request")

func fetch(magnitude pb.Magnitude, past pb.Past) ([]byte, error) {
	url, err := resolveURL(magnitude, past)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	data, err := fetchFromURL(url)
	if err != nil {
		log.Printf("error %v fetching %s", err, url)
	} else {
		kilos := float64(len(data)) / 1024.0
		ms := time.Now().Sub(start).Milliseconds()
		log.Printf("fetched %.1f KB in %d ms from %s", kilos, ms, url)
	}
	return data, err
}

// resolveUrl creates an URL to fetch earthquakes from the GeoJSON Summary
// data sources from USGS
// (see https://earthquake.usgs.gov/earthquakes/feed/v1.0/geojson.php).
func resolveURL(magnitude pb.Magnitude, past pb.Past) (string, error) {
	var magnID string
	var pastID string
	switch magnitude {
	case pb.Magnitude_MAGNITUDE_SIGNIFICANT:
		magnID = "significant"
	case pb.Magnitude_MAGNITUDE_M45_PLUS:
		magnID = "4.5"
	case pb.Magnitude_MAGNITUDE_M25_PLUS:
		magnID = "2.5"
	case pb.Magnitude_MAGNITUDE_M10_PLUS:
		magnID = "1.0"
	case pb.Magnitude_MAGNITUDE_ALL:
		magnID = "all"
	}
	switch past {
	case pb.Past_PAST_HOUR:
		pastID = "hour"
	case pb.Past_PAST_DAY:
		pastID = "day"
	case pb.Past_PAST_7DAYS:
		pastID = "week"
	case pb.Past_PAST_30DAYS:
		pastID = "month"
	}

	if magnID == "" || pastID == "" {
		return "", ErrUnknownDataRequest
	}
	url := apiBaseURL + magnID + "_" + pastID + apiBaseURLPostfix

	return url, nil
}

// fecthFromURL fetches data as []byte from an external HTTP resource
func fetchFromURL(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("resouce %s returned %d", url, resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
