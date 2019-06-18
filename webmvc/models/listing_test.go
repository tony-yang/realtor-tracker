package models_test

import (
	"encoding/json"
	"testing"
	"webmvc/base"
	"webmvc/models"
	"webmvc/tester"
)

func TestFetchMlsListing(t *testing.T) {
	t.Run("returns all listings when requested without subpath", func(t *testing.T) {
		listing := models.Listing{}
		result, _, _ := listing.FetchMlsListing("")

		wanted := make(map[string]models.Listing)

		wanted["19016318"] = models.Listing{
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
		wantedResult, err := json.Marshal(wanted)
		if err != nil {
			base.Error("failed to create json listing:", err)
		}

		tester.AssertStringEqual(t, result, string(wantedResult))
	})
}
