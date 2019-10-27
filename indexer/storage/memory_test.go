package storage

import (
	"testing"
	"time"

	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
)

func TestSaveNewListing(t *testing.T) {
	t.Run("save a new listing", func(t *testing.T) {
		cityIndex := map[string]*City{
			"city,province": {
				Name:      "City",
				State:     "Province",
				MlsNumber: make(map[string]bool),
			},
		}
		mDB, _ := NewMemoryDB(cityIndex)

		mlsNumber := "19016318"
		price := []*mlspb.PriceHistory{
			{
				Price:     10000,
				Timestamp: time.Now().Unix(),
			},
		}
		listings := map[string]*mlspb.Property{
			mlsNumber: {
				Address:       "1234 street|city, province A0B1C2",
				Bathrooms:     "1",
				Bedrooms:      "3 + 0",
				LandSize:      "0X",
				MlsId:         "1234",
				MlsNumber:     mlsNumber,
				MlsUrl:        "/abc/20552312/house",
				Parking:       []string{"None"},
				PhotoUrl:      []string{"https://picture/listings/high/456.jpg"},
				Price:         price,
				PublicRemarks: "HOUSE DESCRIPTION",
				Stories:       "1.5",
				PropertyType:  "House",
				ListTimestamp: 123456789,
				Latitude:      10.1234,
				Longitude:     20.9876,
				City:          "city",
				State:         "province",
				Zipcode:       "A0B1C2",
			},
		}
		if err := mDB.SaveNewListing(listings[mlsNumber]); err != nil {
			t.Errorf("Failed to save the new listing: %v", err)
		}

		if mDB.Mls[mlsNumber].mlsID != listings[mlsNumber].MlsId {
			t.Errorf("mlsID incorrectly saved, expected %d, got %s", listings[mlsNumber].MlsId, mDB.Mls[mlsNumber].mlsID)
		}

		if mDB.Mls[mlsNumber].mlsURL != listings[mlsNumber].MlsUrl {
			t.Errorf("mlsURL incorrectly saved, expected %s, got %s", listings[mlsNumber].MlsUrl, mDB.Mls[mlsNumber].mlsURL)
		}

		if mDB.Property[mlsNumber].address != listings[mlsNumber].Address {
			t.Errorf("property address incorrectly saved, expected %s, got %s", listings[mlsNumber].Address, mDB.Property[mlsNumber].address)
		}

		if mDB.PriceHistory[mlsNumber][0].price != price[0].Price {
			t.Errorf("price incorrectly saved, expected %d, got %d", price[0].Price, mDB.PriceHistory[mlsNumber][0].price)
		}

		if mDB.Property[mlsNumber].longitude != listings[mlsNumber].Longitude {
			t.Errorf("longitude incorrectly saved, expected %f, got %f", listings[mlsNumber].Longitude, mDB.Property[mlsNumber].longitude)
		}

		if mDB.Property[mlsNumber].city != listings[mlsNumber].City {
			t.Errorf("city incorrectly saved, expected %s, got %s", listings[mlsNumber].City, mDB.Property[mlsNumber].city)
		}

		if _, ok := mDB.CityIndex["city,province"].MlsNumber[mlsNumber]; !ok {
			t.Errorf("city index incorrectly saved, expected mlsNumber to exist under city,province, but got %t", ok)
		}
	})

	t.Run("save same listing should reject", func(t *testing.T) {
		cityIndex := map[string]*City{
			"city,province": {
				Name:      "City",
				State:     "Province",
				MlsNumber: make(map[string]bool),
			},
		}
		mDB, _ := NewMemoryDB(cityIndex)

		mlsNumber := "19016319"
		price := []*mlspb.PriceHistory{
			{
				Price:     10000,
				Timestamp: time.Now().Unix(),
			},
		}
		listings := map[string]*mlspb.Property{
			mlsNumber: {
				Address:       "1234 street|city, province A0B1C2",
				Bathrooms:     "1",
				Bedrooms:      "3 + 0",
				LandSize:      "0X",
				MlsId:         "1234",
				MlsNumber:     mlsNumber,
				MlsUrl:        "/abc/20552312/house",
				Parking:       []string{"None"},
				PhotoUrl:      []string{"https://picture/listings/high/456.jpg"},
				Price:         price,
				PublicRemarks: "HOUSE DESCRIPTION",
				Stories:       "1.5",
				PropertyType:  "House",
				ListTimestamp: 123456789,
				Latitude:      10.1234,
				Longitude:     20.9876,
				City:          "city",
				State:         "province",
				Zipcode:       "A0B1C2",
			},
		}
		if err := mDB.SaveNewListing(listings[mlsNumber]); err != nil {
			t.Errorf("Failed to save the new listing: %v", err)
		}

		if err := mDB.SaveNewListing(listings[mlsNumber]); err == nil {
			t.Errorf("Save the same listing should fail: %v", err)
		}
	})
}

func TestReadListings(t *testing.T) {
	t.Run("read a saved listing", func(t *testing.T) {
		cityIndex := map[string]*City{
			"city,province": {
				Name:      "City",
				State:     "Province",
				MlsNumber: make(map[string]bool),
			},
		}
		mDB, _ := NewMemoryDB(cityIndex)

		mlsNumber := "19016320"
		price := []*mlspb.PriceHistory{
			{
				Price:     10000,
				Timestamp: time.Now().Unix(),
			},
		}
		listings := map[string]*mlspb.Property{
			mlsNumber: {
				Address:       "1234 street|city, province A0B1C2",
				Bathrooms:     "1",
				Bedrooms:      "3 + 0",
				LandSize:      "0X",
				MlsId:         "1234",
				MlsNumber:     mlsNumber,
				MlsUrl:        "/abc/20552312/house",
				Parking:       []string{"None"},
				PhotoUrl:      []string{"https://picture/listings/high/456.jpg"},
				Price:         price,
				PublicRemarks: "HOUSE DESCRIPTION",
				Stories:       "1.5",
				PropertyType:  "House",
				ListTimestamp: 123456789,
				Latitude:      10.1234,
				Longitude:     20.9876,
				City:          "city",
				State:         "province",
				Zipcode:       "A0B1C2",
			},
		}
		if err := mDB.SaveNewListing(listings[mlsNumber]); err != nil {
			t.Errorf("Failed to save the new listing: %v", err)
		}
		results, err := mDB.ReadListings()
		if err != nil {
			t.Errorf("Failed to read the saved listing: %v", err)
		}

		if results.Property[0].MlsId != listings[mlsNumber].MlsId {
			t.Errorf("mlsID incorrectly saved, expected %d, got %s", listings[mlsNumber].MlsId, results.Property[0].MlsId)
		}

		if results.Property[0].MlsUrl != listings[mlsNumber].MlsUrl {
			t.Errorf("mlsURL incorrectly saved, expected %s, got %s", listings[mlsNumber].MlsUrl, results.Property[0].MlsUrl)
		}

		if results.Property[0].Address != listings[mlsNumber].Address {
			t.Errorf("property address incorrectly saved, expected %s, got %s", listings[mlsNumber].Address, results.Property[0].Address)
		}

		if results.Property[0].Price[0].Price != price[0].Price {
			t.Errorf("price incorrectly saved, expected %d, got %d", price[0].Price, results.Property[0].Price[0].Price)
		}

		if results.Property[0].Longitude != listings[mlsNumber].Longitude {
			t.Errorf("longitude incorrectly saved, expected %f, got %f", listings[mlsNumber].Longitude, results.Property[0].Longitude)
		}
	})
}
