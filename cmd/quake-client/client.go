// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

// Package main implements a (test) client for the QuakeService
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/navibyte/quake/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type cliParams struct {
	address string
	op      string
}

const (
	defaultAddress       = "localhost:50051"
	timeout              = 30 * time.Second
	unsecureForLocalhost = true
	useRootCA            = true
)

func main() {
	// Require at least command with two args (op + param)
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	// handle help command
	if os.Args[1] == "help" || os.Args[1] == "Help" {
		printUsage()
		os.Exit(0)
	}

	// Get address for the server
	address := os.Getenv("QUAKE_SERVICE")
	if address == "" {
		address = defaultAddress
	}
	fmt.Println("Using address: ", address)

	// context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// open connection to the gRPC server
	var client pb.QuakeServiceClient
	if unsecureForLocalhost && strings.HasPrefix(address, "localhost") {
		// assuming local service with insecure connection
		conn, err := grpc.DialContext(ctx, address,
			grpc.WithInsecure(),
			grpc.WithBlock())
		if err != nil {
			log.Fatalf("failed to connect: %v", err)
		}
		defer conn.Close()
		client = pb.NewQuakeServiceClient(conn)
	} else {
		// assuming remote service (SSL / TLS)
		var conf *tls.Config
		if useRootCA {
			rootCAs, err := x509.SystemCertPool()
			if err != nil {
				log.Fatalf("no system root CAs available: %v", err)
			}
			conf = &tls.Config{
				InsecureSkipVerify: false,
				RootCAs:            rootCAs,
			}
		} else {
			conf = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		conn, err := grpc.DialContext(ctx, address,
			grpc.WithTransportCredentials(credentials.NewTLS(conf)),
			grpc.WithBlock())
		if err != nil {
			log.Fatalf("failed to connect: %v", err)
		}
		defer conn.Close()
		client = pb.NewQuakeServiceClient(conn)
	}
	fmt.Println("Connected to service")

	// Check which method to call and then call it
	switch os.Args[1] {
	case "GetEarthquake":
		r, err := client.GetEarthquake(ctx,
			&pb.GetEarthquakeRequest{
				Id: os.Args[2],
			})
		if err != nil {
			log.Fatalf("failed to get earthquake: %v", err)
		}
		printEarthquake(r.Feature)
		break
	case "ListEarthquakes":
		req, err := parseListEarthquakesRequest()
		if err != nil {
			log.Fatalf("bad request: %v", err)
		}
		r, err := client.ListEarthquakes(ctx, req)
		if err != nil {
			log.Fatalf("failed to list earthquakes: %v", err)
		}
		printEarthquakes(r.Collection)
		break
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: quake-client <method> <params>..")
	fmt.Println("Methods:")
	fmt.Println("  help")
	fmt.Println("  GetEarthquake <id>")
	fmt.Println("     id: {string}")
	fmt.Println("  ListEarthquakes <magnitude> <past> <limit> <details>")
	fmt.Println("     magnitude: significant | 4.5 | 2.5 | 1.0 | all")
	fmt.Println("     past: hour | day | 7days | 30days")
	fmt.Println("     limit: {integer}")
	fmt.Println("     details: true | false")
	fmt.Println("Optionally use env QUAKE_SERVICE to set server address.")
	fmt.Println("Otherwise a default address is used: ", defaultAddress)
}

func printEarthquakes(col *pb.EarthquakeCollection) {
	for _, eq := range col.Features {
		printEarthquake(eq)
	}
}

func printEarthquake(eq *pb.Earthquake) {
	timeFormatted := time.Unix(eq.Time, 0).Format(time.UnixDate)
	fmt.Printf("%s at %s M%.1f near %s",
		eq.Id, timeFormatted, eq.Magnitude, eq.Place)
	fmt.Println("")
}

func parseListEarthquakesRequest() (*pb.ListEarthquakesRequest, error) {
	req := &pb.ListEarthquakesRequest{}
	if len(os.Args) >= 3 {
		switch os.Args[2] {
		case "significant":
			req.Magnitude = pb.Magnitude_MAGNITUDE_SIGNIFICANT
			break
		case "4.5":
			req.Magnitude = pb.Magnitude_MAGNITUDE_M45_PLUS
			break
		case "2.5":
			req.Magnitude = pb.Magnitude_MAGNITUDE_M25_PLUS
			break
		case "1.0":
			req.Magnitude = pb.Magnitude_MAGNITUDE_M10_PLUS
			break
		case "all":
			req.Magnitude = pb.Magnitude_MAGNITUDE_ALL
			break
		default:
			return nil, errors.New("unknown magnitude: " + os.Args[2])
		}
	}
	if len(os.Args) >= 4 {
		switch os.Args[3] {
		case "hour":
			req.Past = pb.Past_PAST_HOUR
			break
		case "day":
			req.Past = pb.Past_PAST_DAY
			break
		case "7days":
			req.Past = pb.Past_PAST_7DAYS
			break
		case "30days":
			req.Past = pb.Past_PAST_30DAYS
			break
		default:
			return nil, errors.New("unknown past: " + os.Args[3])
		}
	}
	if len(os.Args) >= 5 {
		limit, err := strconv.Atoi(os.Args[4])
		if err != nil || limit < 0 {
			return nil, errors.New("invalid limit: " + os.Args[4])
		}
		req.Limit = uint64(limit)
	}
	if len(os.Args) >= 6 {
		switch os.Args[5] {
		case "true":
			req.Details = true
			break
		case "false":
			req.Details = false
			break
		default:
			return nil, errors.New("unknown details: " + os.Args[5])
		}
	}
	return req, nil
}
