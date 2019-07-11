// Package main implements a gRPC indexer server that will serve the collected
// data to the client in a standardized format defined by the proto.
package main

import (
	"context"
	"fmt"
	"net"

  "github.com/sirupsen/logrus"
	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"

	"google.golang.org/grpc"
)

var (
	port = 80
)

type indexerServer struct{}

func (s *indexerServer) GetListing(ctx context.Context, r *mlspb.Request) (*mlspb.Listings, error) {
	logrus.Println("Request =", r)
	listings := &mlspb.Listings{}
	p := &mlspb.Property{
		Address:       "1234 street|city, province A0B1C2",
		Bathrooms:     "1",
		Bedrooms:      "3 + 0",
		LandSize:      "0X",
		MlsId:         "1234",
		MlsNumber:     "19016318",
		MlsUrl:        "/abc/20552312/house",
		Parking:       "None",
		PhotoUrl:      "https://picture/listings/high/456.jpg",
		Price:         "$10,000",
		PublicRemarks: "HOUSE DESCRIPTION",
		Stories:       "1.5",
		PropertyType:  "House",
		ListTimestamp: "123456789",
	}
	listings.Property = append(listings.Property, p)
	return listings, nil
}

func newServer() *indexerServer {
	s := &indexerServer{}
	return s
}

func main() {
	logrus.Println("Starting the Server")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Fatalf("failed to Listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	mlspb.RegisterMlsServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
