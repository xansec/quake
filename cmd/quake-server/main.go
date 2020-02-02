// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

const (
	defaultPort = "50051"
)

func main() {
	// PORT set by Cloud Run or default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// create the server with the actual service injected by registerServer()
	lst, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to open tcp listener: %v", err)
	}
	s := grpc.NewServer()
	registerServer(s)
	if err := s.Serve(lst); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
