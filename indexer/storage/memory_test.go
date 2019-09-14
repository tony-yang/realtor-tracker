package storage

import (
	"strconv"
	"testing"

	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
)

func TestSaveNewListing(t *testing.T) {
	mDB := newMemoryDB()

	mlsNumber := "19016318"
	listings := map[string]*mlspb.Property{
		mlsNumber: {
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
		},
	}
	if err := mDB.SaveNewListing(listings); err != nil {
		t.Errorf("Failed to save the new listing: %v", err)
	}

	if strconv.Itoa(mDB.mls[mlsNumber].mlsID) != listings[mlsNumber].MlsId {
		t.Errorf("mlsID incorrectly saved, expected %d, got %s", listings[mlsNumber].MlsId, mDB.mls[mlsNumber].mlsID)
	}

	if mDB.mls[mlsNumber].mlsURL != listings[mlsNumber].MlsUrl {
		t.Errorf("mlsURL incorrectly saved, expected %s, got %s", listings[mlsNumber].MlsUrl, mDB.mls[mlsNumber].mlsURL)
	}

	if mDB.property[mlsNumber].address != listings[mlsNumber].Address {
		t.Errorf("property address incorrectly saved, expected %s, got %s", listings[mlsNumber].Address, mDB.property[mlsNumber].address)
	}

	if mDB.priceHistory[mlsNumber][0].price != 10000 {
		t.Errorf("price incorrectly saved, expected 10000, got %d", mDB.property[mlsNumber].address)
	}
}
