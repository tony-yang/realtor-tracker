// Package main implements a gRPC indexer client that serves as an example
// for other apps to use to fetch the collected data by the indexer server.
package main

import (
	"context"
	"flag"
	"time"

	"github.com/sirupsen/logrus"
	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:80", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func getListing(c mlspb.MlsServiceClient) {
	logrus.Info("Get the mls property")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := &mlspb.Request{}
	property, err := c.GetListing(ctx, request)
	if err != nil {
		logrus.Fatalf("%v.GetListing() failed, %v", c, err)
	}
	logrus.Info(property)
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			logrus.Fatalf("failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		logrus.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := mlspb.NewMlsServiceClient(conn)

	getListing(client)
}
