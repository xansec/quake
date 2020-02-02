// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

package main

import (
	"context"

	pb "github.com/navibyte/quake/api/v1"
	"github.com/navibyte/quake/pkg/earthquakes/usgs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// registerServer is used by main()
func registerServer(s *grpc.Server) {
	pb.RegisterQuakeServiceServer(s, &server{})
}

// -----------------------------------------------------------------------------

// server implementation for the QuakeService
type server struct {
	pb.UnimplementedQuakeServiceServer
}

func (*server) ListEarthquakes(ctx context.Context,
	req *pb.ListEarthquakesRequest) (*pb.ListEarthquakesResponse, error) {

	// use (USGS) earthquake repository to get collection usings a right method
	var col *pb.EarthquakeCollection
	var err error
	if pos := req.GetPosition(); pos != nil {
		// list earthquakes nearest to the position
		col, err = usgs.ListEarthquakesFocusPosition(
			req.Magnitude, req.Past, int(req.Limit), req.Details, pos)
	} else if bounds := req.GetBounds(); bounds != nil {
		// list earthquakes inside bounds (and earthquakes nearest to the
		// center of bounds coming first on the list)
		col, err = usgs.ListEarthquakesFocusBounds(
			req.Magnitude, req.Past, int(req.Limit), req.Details, bounds)
	} else {
		// list earthquakes on a order they are fetched from USGS
		col, err = usgs.ListEarthquakes(
			req.Magnitude, req.Past, int(req.Limit), req.Details)
	}

	// check if repository returned some error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %s", err.Error())
	}
	if col == nil {
		return nil, status.Errorf(codes.Internal, "internal error: collection nil")
	}

	// no error, so return valid response to RCP caller
	res := &pb.ListEarthquakesResponse{
		Collection: col,
	}
	return res, nil
}

func (*server) GetEarthquake(ctx context.Context,
	req *pb.GetEarthquakeRequest) (*pb.GetEarthquakeResponse, error) {

	// use (USGS) earthquake repository to get a specific earthquake (by id)
	eq, err := usgs.GetEarthquake(req.Id)

	// check if repository returned some error
	if err != nil {
		if err == usgs.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "no earthquake for %s", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "internal error: %s", err.Error())
	}
	if eq == nil {
		return nil, status.Errorf(codes.Internal, "internal error: earthquake nil")
	}

	// no error, so return valid response to RCP caller
	res := &pb.GetEarthquakeResponse{
		Feature: eq,
	}
	return res, nil
}

// -----------------------------------------------------------------------------

// mockServer test implementation for the QuakeService
type mockServer struct {
	pb.UnimplementedQuakeServiceServer
}

func (*mockServer) ListEarthquakes(ctx context.Context,
	req *pb.ListEarthquakesRequest) (*pb.ListEarthquakesResponse, error) {

	// return mock collection
	res := &pb.ListEarthquakesResponse{
		Collection: mockEarthquakeCollection(int(req.Limit), req.Details),
	}
	return res, nil

	//return nil, status.Errorf(codes.Unimplemented, "method ListEarthquakes not implemented")
}
func (*mockServer) GetEarthquake(ctx context.Context,
	req *pb.GetEarthquakeRequest) (*pb.GetEarthquakeResponse, error) {

	// return mock earthquake
	res := &pb.GetEarthquakeResponse{
		Feature: mockEarthquake("Test123", true),
	}
	return res, nil

	//return nil, status.Errorf(codes.Unimplemented, "method GetEarthquake not implemented")
}
