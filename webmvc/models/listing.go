package models

import (
	"context"
	"fmt"
	"time"

	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
	"github.com/tony-yang/realtor-tracker/webmvc/base"

	"google.golang.org/grpc"
)

type Listing struct {
	Address       string
	Bathrooms     string
	Bedrooms      string
	LandSize      string
	MlsId         string
	MlsNumber     string
	MlsUrl        string
	Parking       string
	PhotoUrl      string
	Price         string
	PublicRemarks string
	Stories       string
	PropertyType  string
	ListTimestamp string
}

func (l *Listing) ReadListing(mlsNumber string) (map[string]Listing, error) {
	base.Debug("reading listing: mlsNumber =", mlsNumber)
	listings := make(map[string]Listing)
	listings["19016318"] = Listing{
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

	if property, ok := listings[mlsNumber]; ok {
		return map[string]Listing{mlsNumber: property}, nil
	} else {
		return listings, nil
	}
}

func (l *Listing) ReadAllListings() (*mlspb.Listings, error) {
	addr := "127.0.0.1:9000"
	opts := []grpc.DialOption{grpc.WithInsecure()}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := mlspb.NewMlsServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := &mlspb.Request{}
	listings, err := c.GetListing(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("Failed to get read all listings: %v", err)
	}

	return listings, nil
}
