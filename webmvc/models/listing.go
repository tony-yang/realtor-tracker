package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"webmvc/base"

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

func ReadListing(database *sql.DB, mlsNumber int) (Listing, error) {
	return Listing{}, nil
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

func SaveNewListing(database *sql.DB, listings map[string]Listing) error {
	err := CreateDB(database)
	if err != nil {
		return fmt.Errorf("failed to create DB: %s", err )
	}

	for mlsNumber, listing := range listings {
		fmt.Println("######### mlsNumber =", mlsNumber, "listing", listing)
		statement, err := database.Prepare(`INSERT INTO mls (
				mlsNumber, mlsId, mlsUrl, bathrooms, bedrooms, landSize, parking,
				publicRemark, stories, propertyType, availableTimestamp, statusId)
				VALUES(?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?)`)
		if err != nil {
			return fmt.Errorf("error prepare insert statement to mls: %s", err)
		}
		statement.Exec(
			listing.MlsNumber, listing.MlsId, listing.MlsUrl, listing.Bathrooms,
			listing.Bedrooms, listing.LandSize, listing.Parking,
			listing.PublicRemarks, listing.Stories, listing.PropertyType,
			time.Now().Unix(), 1)

		statement, err = database.Prepare(`INSERT INTO property (
				address, mlsNumber)
				VALUES(?, ?)`)
		if err != nil {
			return fmt.Errorf("error prepare insert statement to property: %s", err)
		}
		statement.Exec(listing.Address, listing.MlsNumber)

		statement, err = database.Prepare(`INSERT INTO photo (
				photoUrl, mlsNumber)
				VALUES(?, ?)`)
		if err != nil {
			return fmt.Errorf("error prepare insert statement to photo: %s", err)
		}
		statement.Exec(listing.PhotoUrl, listing.MlsNumber)

		statement, err = database.Prepare(`INSERT INTO priceHistory (
				mlsNumber, price, priceTimestamp)
				VALUES(?, ?, ?)`)
		if err != nil {
			return fmt.Errorf("error prepare insert statement to photo: %s", err)
		}
		statement.Exec(listing.MlsNumber, listing.Price, time.Now().Unix())
	}

	return nil
}

func CreateDB(database *sql.DB) error {
	statement, err := database.Prepare(`CREATE TABLE IF NOT EXISTS listingStatus (
		statusId INTEGER PRIMARY KEY,
		status TEXT)`)
	if err != nil {
		return fmt.Errorf("error prepare create the listingStatus table statement: %s", err)
	}
	statement.Exec()

	statement, err = database.Prepare(`INSERT INTO listingStatus (
			statusId, status)
			VALUES(?, ?)`)
	if err != nil {
		return fmt.Errorf("error prepare insert statement to photo: %s", err)
	}
	statement.Exec(1, "OPEN")
	statement.Exec(1, "SOLD")

	statement, err = database.Prepare(`CREATE TABLE IF NOT EXISTS mls (
		mlsNumber INTEGER PRIMARY KEY,
		mlsId INTEGER,
		mlsUrl TEXT,
		bathrooms TEXT,
		bedrooms TEXT,
		landSize TEXT,
		parking TEXT,
		publicRemark TEXT,
		stories TEXT,
		propertyType TEXT,
		availableTimestamp INTEGER,
		statusId INTEGER,
		FOREIGN KEY(statusId) REFERENCES listingStatus(statusId))`)
	if err != nil {
		return fmt.Errorf("error prepare create the mls table statement: %s", err)
	}
	statement.Exec()

	statement, err = database.Prepare(`CREATE TABLE IF NOT EXISTS property (
		address TEXT PRIMARY KEY,
		mlsNumber INTEGER,
		FOREIGN KEY(mlsNumber) REFERENCES mls(mlsNumber))`)
	if err != nil {
		return fmt.Errorf("error prepare create the property table statement: %s", err)
	}
	statement.Exec()

	statement, err = database.Prepare(`CREATE TABLE IF NOT EXISTS photo (
		photoUrl TEXT PRIMARY KEY,
		mlsNumber INTEGER,
		FOREIGN KEY(mlsNumber) REFERENCES mls(mlsNumber))`)
	if err != nil {
		return fmt.Errorf("error prepare create the photo table statement: %s", err)
	}
	statement.Exec()

	statement, err = database.Prepare(`CREATE TABLE IF NOT EXISTS priceHistory (
		mlsNumber INTEGER,
		price INTEGER,
		priceTimestamp INTEGER,
		FOREIGN KEY(mlsNumber) REFERENCES mls(mlsNumber))`)
	if err != nil {
		return fmt.Errorf("error prepare create the priceHistory table statement: %s", err)
	}
	statement.Exec()

	return nil
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
		ListTimestamp: "123456789",
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

	database, err := sql.Open("sqlite3", "/tmp/realtor.db")
	if err != nil {
		base.Error("failed to create a new database:", err)
		statusCode = http.StatusInternalServerError
	}
	err = CreateDB(database)
	if err != nil {
		fmt.Println("create db failed", err)
	}
	err = SaveNewListing(database, result)
	if err != nil {
		fmt.Println("save new listing failed", err)
	}

	_, err = ReadAllListings(database)
	if err != nil {
		fmt.Println("read all listings failed:", err)
	}

	return string(listings), statusCode, err
}
