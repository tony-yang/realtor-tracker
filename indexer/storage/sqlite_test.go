package storage

import (
	"os"
	"testing"
	"time"

	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
)

func cleanSqliteDB(dbPath string) error {
	return os.Remove(dbPath)
}

func TestSqliteSaveNewListing(t *testing.T) {
	t.Run("save a new listing", func(t *testing.T) {
		var dbPath = "/tmp/realtor.db"
		db, err := NewSqliteDB(dbPath)
		if err != nil {
			t.Error(err)
		}

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
		if err := db.SaveNewListing(listings[mlsNumber]); err != nil {
			t.Errorf("Failed to save the new listing: %v", err)
		}

		if err := cleanSqliteDB(dbPath); err != nil {
			t.Errorf("Failed to cleanup the test sqlite db: %v", err)
		}
	})

	t.Run("save same listing should reject", func(t *testing.T) {
		var dbPath = "/tmp/realtor2.db"
		db, err := NewSqliteDB(dbPath)
		if err != nil {
			t.Error(err)
		}

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
		if err := db.SaveNewListing(listings[mlsNumber]); err != nil {
			t.Errorf("Failed to save the new listing: %v", err)
		}

		if err := db.SaveNewListing(listings[mlsNumber]); err == nil {
			t.Errorf("Save the same listing should fail: %v", err)
		}
		if err := cleanSqliteDB(dbPath); err != nil {
			t.Errorf("Failed to cleanup the test sqlite db: %v", err)
		}
	})
}

func TestSqliteReadListings(t *testing.T) {
	t.Run("read a saved listing", func(t *testing.T) {
		var dbPath = "/tmp/realtor3.db"
		db, err := NewSqliteDB(dbPath)
		if err != nil {
			t.Error(err)
		}

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
		if err := db.SaveNewListing(listings[mlsNumber]); err != nil {
			t.Errorf("Failed to save the new listing: %v", err)
		}
		results, err := db.ReadListings()
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

		if err := cleanSqliteDB(dbPath); err != nil {
			t.Errorf("Failed to cleanup the test sqlite db: %v", err)
		}
	})
}
