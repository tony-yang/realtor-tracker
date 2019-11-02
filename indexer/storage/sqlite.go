package storage

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
)

var dbCreated = false

// DB creates the sqlite DB reference used to store data locally.
type SqliteDB struct {
	db *sql.DB
}

// NewDBStorage creates an instance of the sqlite database used to store
// the data locally.
func NewSqliteDB(dbPath string) (*SqliteDB, error) {
	// dbCreated is reset whenever a new sqlite DB is initiated.
	dbCreated = false
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new database: %v", err)
	}

	return &SqliteDB{db: database}, nil
}

func (d *SqliteDB) createListingStatusTable() error {
	sqlStatement := `CREATE TABLE IF NOT EXISTS listingStatus (
		statusId INTEGER PRIMARY KEY,
		status TEXT UNIQUE)`
	statement, err := d.db.Prepare(sqlStatement)
	if err != nil {
		return fmt.Errorf("error prepare the create listingStatus table: %v", err)
	}
	if _, err := statement.Exec(); err != nil {
		return fmt.Errorf("error execute %q: %v", sqlStatement, err)
	}

	sqlStatement = `INSERT INTO listingStatus (status)
			VALUES(?)`
	statement, err = d.db.Prepare(sqlStatement)
	if err != nil {
		return fmt.Errorf("error prepare the insert statement to listingStatus: %v", err)
	}
	if _, err := statement.Exec("Open"); err != nil {
		return fmt.Errorf("error execute %q with value Open: %v", sqlStatement, err)
	}
	if _, err := statement.Exec("Pending"); err != nil {
		return fmt.Errorf("error execute %q with value Pending: %v", sqlStatement, err)
	}
	if _, err := statement.Exec("Sold"); err != nil {
		return fmt.Errorf("error execute %q with value Sold: %v", sqlStatement, err)
	}
	if _, err := statement.Exec("Closed"); err != nil {
		return fmt.Errorf("error execute %q with value Closed: %v", sqlStatement, err)
	}
	return nil
}

func (d *SqliteDB) createCityTable() error {
	sqlStatement := `CREATE TABLE IF NOT EXISTS city (
		name TEXT NOT NULL,
		state TEXT NOT NULL,
		PRIMARY KEY (name, state))`
	statement, err := d.db.Prepare(sqlStatement)
	if err != nil {
		return fmt.Errorf("error prepare the create city table: %v", err)
	}
	if _, err := statement.Exec(); err != nil {
		return fmt.Errorf("error execute %q: %v", sqlStatement, err)
	}

	// statement, err = d.db.Prepare(`INSERT INTO city (name, state)
	// 		VALUES(?, ?)`)
	// if err != nil {
	// 	return fmt.Errorf("error prepare the insert statement to city: %v", err)
	// }
	// statement.Exec("Windsor", "Ontario")
	return nil
}

func (d *SqliteDB) createPropertyTable() error {
	sqlStatement := `CREATE TABLE IF NOT EXISTS property (
		address TEXT PRIMARY KEY,
		zipcode TEXT NOT NULL,
		latitude REAL,
		longitude REAL,
		city TEXT,
		state TEXT,
		FOREIGN KEY(city, state) REFERENCES city(name, state))`
	statement, err := d.db.Prepare(sqlStatement)
	if err != nil {
		return fmt.Errorf("error prepare the create property table: %v", err)
	}
	if _, err := statement.Exec(); err != nil {
		return fmt.Errorf("error execute %q: %v", sqlStatement, err)
	}
	return nil
}

func (d *SqliteDB) createMlsTable() error {
	sqlStatement := `CREATE TABLE IF NOT EXISTS mls (
		mlsNumber TEXT PRIMARY KEY,
		mlsId TEXT,
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
		source TEXT,
		address TEXT,
 		FOREIGN KEY(statusId) REFERENCES listingStatus(statusId),
		FOREIGN KEY(address) REFERENCES property(address))`
	statement, err := d.db.Prepare(sqlStatement)
	if err != nil {
		return fmt.Errorf("error prepare the create mls table: %v", err)
	}
	if _, err := statement.Exec(); err != nil {
		return fmt.Errorf("error execute %q: %v", sqlStatement, err)
	}
	return nil
}

func (d *SqliteDB) createPhotoTable() error {
	sqlStatement := `CREATE TABLE IF NOT EXISTS photo (
		photoUrl TEXT PRIMARY KEY,
		mlsNumber TEXT,
		FOREIGN KEY(mlsNumber) REFERENCES mls(mlsNumber))`
	statement, err := d.db.Prepare(sqlStatement)
	if err != nil {
		return fmt.Errorf("error prepare the create photo table: %v", err)
	}
	if _, err := statement.Exec(); err != nil {
		return fmt.Errorf("error execute %q: %v", sqlStatement, err)
	}
	return nil
}

func (d *SqliteDB) createPriceHistoryTable() error {
	sqlStatement := `CREATE TABLE IF NOT EXISTS priceHistory (
		mlsNumber TEXT,
		price INTEGER,
		priceTimestamp INTEGER,
		FOREIGN KEY(mlsNumber) REFERENCES mls(mlsNumber))`
	statement, err := d.db.Prepare(sqlStatement)
	if err != nil {
		return fmt.Errorf("error prepare the create priceHistory table: %v", err)
	}
	if _, err := statement.Exec(); err != nil {
		return fmt.Errorf("error execute %q: %v", sqlStatement, err)
	}
	return nil
}

// CreateStorage for sqlite DB to create all the tables during module first use.
func (d *SqliteDB) CreateStorage() error {
	if dbCreated {
		return nil
	}
	if err := d.createListingStatusTable(); err != nil {
		return err
	}
	if err := d.createCityTable(); err != nil {
		return err
	}
	if err := d.createPropertyTable(); err != nil {
		return err
	}
	if err := d.createMlsTable(); err != nil {
		return err
	}
	if err := d.createPhotoTable(); err != nil {
		return err
	}
	if err := d.createPriceHistoryTable(); err != nil {
		return err
	}
	dbCreated = true
	return nil
}

// UpdateListing appends new pricing information for an existing listing record.
func (d *SqliteDB) UpdateListing(p *mlspb.Property) error {
	logrus.Debugf("update listing: mlsNumber = %s listing %v\n", p.MlsNumber, p)

	if err := d.insertPriceHistory(p); err != nil {
		return fmt.Errorf("failed to insert a price history with err: %v", err)
	}

	return nil
}

func (d *SqliteDB) listingExisted(mlsNumber string) bool {
	rows, err := d.db.Query(`SELECT * FROM mls WHERE mlsNumber = $1`, mlsNumber)
	if err != nil {
		return false
	}
	return rows.Next()
}

func (d *SqliteDB) cityExisted(city, state string) bool {
	rows, err := d.db.Query(`SELECT * FROM city WHERE name = $1 AND state = $2`, city, state)
	if err != nil {
		return false
	}
	return rows.Next()
}

func (d *SqliteDB) insertCity(city, state string) error {
	sqlStatement := `INSERT INTO city (name, state)
			VALUES(?, ?)`
	statement, err := d.db.Prepare(sqlStatement)
	if err != nil {
		return fmt.Errorf("error prepare the insert city: %v", err)
	}
	if _, err := statement.Exec(city, state); err != nil {
		return fmt.Errorf("error execute %q: %v", sqlStatement, err)
	}
	return nil
}

func (d *SqliteDB) insertProperty(p *mlspb.Property) error {
	sqlStatement := `INSERT INTO property (
			address, zipcode, latitude, longitude, city, state)
			VALUES(?, ?, ?, ?, ?, ?)`
	statement, err := d.db.Prepare(sqlStatement)
	if err != nil {
		return fmt.Errorf("error prepare the insert property: %v", err)
	}
	if _, err := statement.Exec(p.Address, p.Zipcode, p.Latitude, p.Longitude, p.City, p.State); err != nil {
		return fmt.Errorf("error execute %q: %v", sqlStatement, err)
	}
	return nil
}

func (d *SqliteDB) insertMls(p *mlspb.Property) error {
	sqlStatement := `INSERT INTO mls (
			mlsNumber, mlsId, mlsUrl, bathrooms, bedrooms, landSize, parking,
			publicRemark, stories, propertyType, availableTimestamp, statusId, source, address)
			VALUES(?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?)`
	statement, err := d.db.Prepare(sqlStatement)
	if err != nil {
		return fmt.Errorf("error prepare the insert mls: %v", err)
	}
	if _, err := statement.Exec(
		p.MlsNumber, p.MlsId, p.MlsUrl, p.Bathrooms, p.Bedrooms, p.LandSize, strings.Join(p.Parking, ";"),
		p.PublicRemarks, p.Stories, p.PropertyType, p.ListTimestamp, 1, p.Source, p.Address); err != nil {
		return fmt.Errorf("error execute %q: %v", sqlStatement, err)
	}
	return nil
}

func (d *SqliteDB) insertPhoto(p *mlspb.Property) error {
	for _, ph := range p.PhotoUrl {
		sqlStatement := `INSERT INTO photo (
				photoUrl, mlsNumber)
				VALUES(?, ?)`
		statement, err := d.db.Prepare(sqlStatement)
		if err != nil {
			return fmt.Errorf("error prepare the insert photo: %v", err)
		}
		if _, err := statement.Exec(ph, p.MlsNumber); err != nil {
			return fmt.Errorf("error execute %q: %v", sqlStatement, err)
		}
	}
	return nil
}

func (d *SqliteDB) insertPriceHistory(p *mlspb.Property) error {
	for _, pr := range p.Price {
		sqlStatement := `INSERT INTO priceHistory (
				mlsNumber, price, priceTimestamp)
				VALUES(?, ?, ?)`
		statement, err := d.db.Prepare(sqlStatement)
		if err != nil {
			return fmt.Errorf("error prepare insert statement to photo: %s", err)
		}
		if _, err := statement.Exec(p.MlsNumber, pr.Price, pr.Timestamp); err != nil {
			return fmt.Errorf("error execute %q: %v", sqlStatement, err)
		}
	}
	return nil
}

// SaveNewListing saves the data collected into the in-memory data structure.
func (d *SqliteDB) SaveNewListing(p *mlspb.Property) error {
	if err := d.CreateStorage(); err != nil {
		return fmt.Errorf("failed to create DB: %s", err)
	}

	logrus.Debugf("save mls: %q", p.MlsNumber)
	if d.listingExisted(p.MlsNumber) {
		return fmt.Errorf("listing %s exists", p.MlsNumber)
	}
	city := strings.ToLower(p.City)
	state := strings.ToLower(p.State)
	if !d.cityExisted(city, state) {
		if err := d.insertCity(city, state); err != nil {
			return fmt.Errorf("failed to insert a new city '%s, %s' with err: %v", city, state, err)
		}
	}

	if err := d.insertProperty(p); err != nil {
		return fmt.Errorf("failed to insert a new property with err: %v", err)
	}

	if err := d.insertMls(p); err != nil {
		return fmt.Errorf("failed to insert a new mls listing %q with err: %v", p.MlsNumber, err)
	}

	if err := d.insertPhoto(p); err != nil {
		return fmt.Errorf("failed to insert a listing photo with err: %v", err)
	}

	if err := d.insertPriceHistory(p); err != nil {
		return fmt.Errorf("failed to insert a price history with err: %v", err)
	}

	return nil
}

func (d *SqliteDB) ReadListing(mlsNumber string) (string, error) {
	return "", nil
}

func (d *SqliteDB) ReadListings() (*mlspb.Listings, error) {
	listings := &mlspb.Listings{}
	rows, err := d.db.Query(`SELECT mlsNumber, mlsId, mlsUrl, bathrooms, bedrooms, landSize, publicRemark, stories, propertyType, availableTimestamp, status, source, mls.address, zipcode, city, state, parking, latitude, longitude
		FROM mls
		INNER JOIN property ON mls.address = property.address
		INNER JOIN listingStatus ON mls.statusId = listingStatus.statusId
		WHERE status = "Open" LIMIT 10`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			mlsNumber, mlsID, mlsURL, bathrooms, bedrooms, landSize, publicRemark, stories, propertyType, status, source, address, zipcode, city, state, parking string
			availableTimestamp                                                                                                                                   int64
			latitude, longitude                                                                                                                                  float64
		)
		if err := rows.Scan(&mlsNumber, &mlsID, &mlsURL, &bathrooms, &bedrooms, &landSize, &publicRemark, &stories, &propertyType, &availableTimestamp, &status, &source, &address, &zipcode, &city, &state, &parking, &latitude, &longitude); err != nil {
			return nil, err
		}
		parkings := []string{parking}

		photoRows, err := d.db.Query(`SELECT photoUrl FROM photo WHERE mlsNumber = $1`, mlsNumber)
		if err != nil {
			return nil, err
		}
		var photos []string
		for photoRows.Next() {
			var photoUrl string
			if err := photoRows.Scan(&photoUrl); err != nil {
				return nil, err
			}
			photos = append(photos, photoUrl)
		}

		priceRows, err := d.db.Query(`SELECT price, priceTimestamp FROM priceHistory WHERE mlsNumber = $1`, mlsNumber)
		if err != nil {
			return nil, err
		}
		prices := []*mlspb.PriceHistory{}
		for priceRows.Next() {
			var p int32
			var t int64
			if err := priceRows.Scan(&p, &t); err != nil {
				return nil, err
			}
			prices = append(prices, &mlspb.PriceHistory{
				Price:     p,
				Timestamp: t,
			})
		}

		p := &mlspb.Property{
			Address:       address,
			Bathrooms:     bathrooms,
			Bedrooms:      bedrooms,
			LandSize:      landSize,
			MlsId:         mlsID,
			MlsNumber:     mlsNumber,
			MlsUrl:        mlsURL,
			Parking:       parkings,
			PhotoUrl:      photos,
			Price:         prices,
			PublicRemarks: publicRemark,
			Stories:       stories,
			PropertyType:  propertyType,
			ListTimestamp: availableTimestamp,
			Source:        source,
			Latitude:      latitude,
			Longitude:     longitude,
			City:          city,
			State:         state,
			Zipcode:       zipcode,
			Status:        status,
		}
		listings.Property = append(listings.Property, p)
	}
	return listings, nil
}
