package models

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/tony-yang/realtor-tracker/webmvc/base"

	_ "github.com/mattn/go-sqlite3"
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

func ReadAllListings(database *sql.DB) (map[string]Listing, error) {
	listings := make(map[string]Listing)

	rows, err := database.Query("SELECT * FROM mls")
	if err != nil {
	  return nil, fmt.Errorf("error query the mls table: %s", err)
	}
	var (
	  mlsNumber int
	  mlsId int
	  mlsUrl string
	  bathrooms string
	  bedrooms string
	  landSize string
	  parking string
	  publicRemark string
	  stories string
	  propertyType string
	  listTimestamp int
		statusId int
	)
	for rows.Next() {
	  rows.Scan(&mlsNumber, &mlsId, &mlsUrl, &bathrooms, &bedrooms, &landSize,
	      &parking, &publicRemark, &stories, &propertyType, &listTimestamp, &statusId)
		listings[strconv.Itoa(mlsNumber)] = Listing{
			Bathrooms: bathrooms,
			Bedrooms: bedrooms,
			LandSize: landSize,
			MlsId: strconv.Itoa(mlsId),
			MlsUrl: mlsUrl,
			Parking: parking,
			PublicRemarks: publicRemark,
			Stories: stories,
			PropertyType: propertyType,
			ListTimestamp: strconv.Itoa(listTimestamp),
		}
	  fmt.Println("#### listings", listings[strconv.Itoa(mlsNumber)])
	}

	rows, err = database.Query("SELECT * FROM property")
	if err != nil {
	  return nil, fmt.Errorf("error query the property table: %s", err)
	}
	var (
	  address string
	)
	for rows.Next() {
	  rows.Scan(&address, &mlsNumber)
		// listings[strconv.Itoa(mlsNumber)].Address = address
	  fmt.Printf("property: mlsNum = %s address = %s\n", strconv.Itoa(mlsNumber), address)
	}

	rows, err = database.Query("SELECT * FROM photo")
	if err != nil {
	  return nil, fmt.Errorf("error query the photo table: %s", err)
	}
	var (
	  photoUrl string
	)
	for rows.Next() {
	  rows.Scan(&photoUrl, &mlsNumber)
		// listings[strconv.Itoa(mlsNumber)].PhotoUrl = photoUrl
	  fmt.Printf("photo: mlsNum = %s photoUrl = %s\n", strconv.Itoa(mlsNumber), photoUrl)
	}

	rows, err = database.Query("SELECT * FROM priceHistory")
	if err != nil {
	  return nil, fmt.Errorf("error query the priceHistory table: %s", err)
	}
	var (
	  price string
		priceTimestamp int
	)
	for rows.Next() {
	  rows.Scan(&mlsNumber, &price, &priceTimestamp)
		// listings[strconv.Itoa(mlsNumber)].Price = price
	  fmt.Printf("price: mlsNum = %s price =  %s time = %d\n", strconv.Itoa(mlsNumber), price, priceTimestamp)
	}

	return listings, nil
}
