package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
	"github.com/tony-yang/realtor-tracker/webmvc/base"
	"github.com/tony-yang/realtor-tracker/webmvc/models"
)

type Listing struct {
	base.Controller
	models.Listing
}

func (l *Listing) Get(subpath string, queries map[string]string) *base.HttpResponse {
	fmt.Println("subpath =", subpath, "queries =", queries)
	statusCode := http.StatusOK
	var bodyContent *mlspb.Listings
	var err error

	if subpath != "" {
		// bodyContent, err = l.ReadListing(subpath)
		fmt.Println("given subpath query individual")
	} else {
		bodyContent, err = l.ReadAllListings()
	}
	if err != nil {
		base.Error("error fetch mls listing:", err)
		statusCode = http.StatusInternalServerError
	}

	body, err := json.Marshal(bodyContent)
	if err != nil {
		body = []byte{}
		statusCode = http.StatusInternalServerError
	}

	return &base.HttpResponse{
		Body:       string(body),
		StatusCode: statusCode,
	}
}
