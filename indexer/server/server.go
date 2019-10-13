// Package main implements a gRPC indexer server that will serve the collected
// data to the client in a standardized format defined by the proto.
package server

import (
	"context"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/tony-yang/realtor-tracker/indexer/collector"
	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"

	"google.golang.org/grpc"
)

var (
	port = 80
)

type indexerServer struct{}

func (s *indexerServer) GetListing(ctx context.Context, r *mlspb.Request) (*mlspb.Listings, error) {
	listings := &mlspb.Listings{}

	for name, c := range collector.Collectors {
		logrus.Infof("Read from the '%s' collector", name)
		result, err := c.GetDB().ReadListings()
		if err != nil {
			logrus.Errorf("reading property listing failed: %v", err)
		}
		logrus.Debug(result.String())
		listings.Property = append(listings.Property, result.Property...)
	}
	return listings, nil
}

func newServer() *indexerServer {
	s := &indexerServer{}
	return s
}

func StartServer() {
	logrus.Info("Starting the Server")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Fatalf("failed to Listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	mlspb.RegisterMlsServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
