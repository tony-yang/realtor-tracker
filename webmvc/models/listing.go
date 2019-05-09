package models

import (
	"encoding/json"
	"net/http"
	"webmvc/base"
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
}

func (l *Listing) FetchMlsListing(subpath string) (string, int, error) {
	var (
		listings   []byte
		err        error
		result     = make(map[string]Listing)
		statusCode = http.StatusOK
	)

	// TODO: Convert this to fetch from a backend storage
	result["19016318"] = Listing{
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
	}

	if property, ok := result[subpath]; ok {
		listings, err = json.Marshal(property)
	} else {
		listings, err = json.Marshal(result)
	}

	if err != nil {
		base.Error("failed to create json listing:", err)
		statusCode = http.StatusInternalServerError
	}

	return string(listings), statusCode, err
}
