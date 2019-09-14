package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
)

type DB struct {
	db *sql.DB
}

func NewDBStorage() (*DB, error) {
	database, err := sql.Open("sqlite3", "/tmp/realtor.db")
	if err != nil {
		logrus.Fatalf("failed to create a new database: %v", err)
		return nil, err
	}

	return &DB{db: database}, nil
}

func (d *DB) CreateStorage() error {
	statement, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS listingStatus (
		statusId INTEGER PRIMARY KEY,
		status TEXT)`)
	if err != nil {
		return fmt.Errorf("error prepare create the listingStatus table statement: %s", err)
	}
	statement.Exec()

	statement, err = d.db.Prepare(`INSERT INTO listingStatus (
			statusId, status)
			VALUES(?, ?)`)
	if err != nil {
		return fmt.Errorf("error prepare insert statement to photo: %s", err)
	}
	statement.Exec(1, "OPEN")
	statement.Exec(1, "SOLD")

	statement, err = d.db.Prepare(`CREATE TABLE IF NOT EXISTS mls (
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

	statement, err = d.db.Prepare(`CREATE TABLE IF NOT EXISTS property (
		address TEXT PRIMARY KEY,
		mlsNumber INTEGER,
		FOREIGN KEY(mlsNumber) REFERENCES mls(mlsNumber))`)
	if err != nil {
		return fmt.Errorf("error prepare create the property table statement: %s", err)
	}
	statement.Exec()

	statement, err = d.db.Prepare(`CREATE TABLE IF NOT EXISTS photo (
		photoUrl TEXT PRIMARY KEY,
		mlsNumber INTEGER,
		FOREIGN KEY(mlsNumber) REFERENCES mls(mlsNumber))`)
	if err != nil {
		return fmt.Errorf("error prepare create the photo table statement: %s", err)
	}
	statement.Exec()

	statement, err = d.db.Prepare(`CREATE TABLE IF NOT EXISTS priceHistory (
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

func (d *DB) SaveNewListing(listings map[string]*mlspb.Property) error {
	if err := d.CreateStorage(); err != nil {
		return fmt.Errorf("failed to create DB: %s", err)
	}

	for mlsNumber, listing := range listings {
		fmt.Println("######### mlsNumber =", mlsNumber, "listing", listing)
		statement, err := d.db.Prepare(`INSERT INTO mls (
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

		statement, err = d.db.Prepare(`INSERT INTO property (
				address, mlsNumber)
				VALUES(?, ?)`)
		if err != nil {
			return fmt.Errorf("error prepare insert statement to property: %s", err)
		}
		statement.Exec(listing.Address, listing.MlsNumber)

		statement, err = d.db.Prepare(`INSERT INTO photo (
				photoUrl, mlsNumber)
				VALUES(?, ?)`)
		if err != nil {
			return fmt.Errorf("error prepare insert statement to photo: %s", err)
		}
		statement.Exec(listing.PhotoUrl, listing.MlsNumber)

		statement, err = d.db.Prepare(`INSERT INTO priceHistory (
				mlsNumber, price, priceTimestamp)
				VALUES(?, ?, ?)`)
		if err != nil {
			return fmt.Errorf("error prepare insert statement to photo: %s", err)
		}
		statement.Exec(listing.MlsNumber, listing.Price, time.Now().Unix())
	}

	return nil
}

func (d *DB) ReadListing(id string) (string, error) {
	return "", nil
}

func (d *DB) ReadListings() (string, error) {
	return "", nil
}
