package storage

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
)

type listingStatus int

const (
	Open listingStatus = iota
	Pending
	Sold
	Closed
)

type City struct {
	Name      string
	State     string
	MlsNumber map[string]bool
}

type mls struct {
	mlsID              string
	mlsURL             string
	bathrooms          string
	bedrooms           string
	landSize           string
	parking            []string
	publicRemark       string
	stories            string
	propertyType       string
	availableTimestamp int64
	status             listingStatus
	source             string
}

type property struct {
	address   string
	zipcode   string
	latitude  float64
	longitude float64
	city      string
	state     string
}

type photo struct {
	photoURL []string
}

type priceHistory struct {
	price     int32
	timestamp int64
}

// MemoryDB creates the in-memory data structure to hold the collected data.
type MemoryDB struct {
	Lock         sync.Mutex
	Mls          map[string]*mls
	Property     map[string]*property
	Photo        map[string]*photo
	PriceHistory map[string][]*priceHistory
	CityIndex    map[string]*City
}

// NewMemoryDB creates a instance of all the in-memory data structure used to
// hold the collected data.
// cityIndex is in the format of map[string]*City
// ie. map[string]*City{
// 	 "windsor,ontario": &City{
// 		 Name: "Windsor",
// 		 State: "Ontario",
// 		 MlsNumber: make(map[string]bool)
// 	 }
// }
func NewMemoryDB(cityIndex map[string]*City) *MemoryDB {
	return &MemoryDB{
		Mls:          make(map[string]*mls),
		Property:     make(map[string]*property),
		Photo:        make(map[string]*photo),
		PriceHistory: make(map[string][]*priceHistory),
		CityIndex:    cityIndex,
	}
}

// CreateStorage for in-memory DB is a placeholder to comply with the DBInterface.
func (m *MemoryDB) CreateStorage() error {
	return nil
}

// UpdateListing appends new pricing information for an existing listing record
func (m *MemoryDB) UpdateListing(mlsNumber string, p *mlspb.Property) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	logrus.Debugf("update listing: mlsNumber = %s listing %v\n", mlsNumber, p)
	if _, ok := m.PriceHistory[mlsNumber]; !ok {
		return fmt.Errorf("Listing %s does not exist", mlsNumber)
	}

	for _, p := range p.Price {
		price := &priceHistory{
			price:     p.Price,
			timestamp: p.Timestamp,
		}
		m.PriceHistory[mlsNumber] = append(m.PriceHistory[mlsNumber], price)
	}
	return nil
}

// SaveNewListing saves the data collected into the in-memory data structure.
func (m *MemoryDB) SaveNewListing(mlsNumber string, p *mlspb.Property) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	logrus.Debugf("Save Listing: mlsNumber = %s listing %v\n", mlsNumber, p)
	if _, ok := m.Mls[mlsNumber]; ok {
		return fmt.Errorf("Listing %s exists", mlsNumber)
	}
	cityKey := fmt.Sprintf("%s,%s", strings.ToLower(p.City), strings.ToLower(p.State))
	// logrus.Infof("### city key = %s", cityKey)
	// logrus.Infof("		city index = %v", m.CityIndex[cityKey])
	if c, ok := m.CityIndex[cityKey]; ok != true {
		m.CityIndex[cityKey] = &City{
			Name:      p.City,
			State:     p.State,
			MlsNumber: map[string]bool{mlsNumber: true},
		}
	} else {
		c.MlsNumber[mlsNumber] = true
	}

	m.Mls[mlsNumber] = &mls{
		mlsID:              p.MlsId,
		mlsURL:             p.MlsUrl,
		bathrooms:          p.Bathrooms,
		bedrooms:           p.Bedrooms,
		landSize:           p.LandSize,
		parking:            p.Parking,
		publicRemark:       p.PublicRemarks,
		stories:            p.Stories,
		propertyType:       p.PropertyType,
		availableTimestamp: p.ListTimestamp,
		status:             Open,
		source:             p.Source,
	}
	m.Property[mlsNumber] = &property{
		address:   p.Address,
		zipcode:   p.Zipcode,
		latitude:  p.Latitude,
		longitude: p.Longitude,
		city:      p.City,
		state:     p.State,
	}
	m.Photo[mlsNumber] = &photo{photoURL: p.PhotoUrl}
	m.PriceHistory[mlsNumber] = []*priceHistory{}
	for _, p := range p.Price {
		price := &priceHistory{
			price:     p.Price,
			timestamp: p.Timestamp,
		}
		m.PriceHistory[mlsNumber] = append(m.PriceHistory[mlsNumber], price)
	}
	return nil
}

// ReadListing reads a listing by listing ID from the in-memory data structure.
func (m *MemoryDB) ReadListing(id string) (string, error) {
	return "", nil
}

// ReadListings reads all MLS listings collected from the in-memory data structure.
func (m *MemoryDB) ReadListings() (*mlspb.Listings, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	listings := &mlspb.Listings{}
	for mlsNumber, mls := range m.Mls {
		price := []*mlspb.PriceHistory{}
		for _, p := range m.PriceHistory[mlsNumber] {
			price = append(price, &mlspb.PriceHistory{
				Price:     p.price,
				Timestamp: p.timestamp,
			})
		}
		p := &mlspb.Property{
			Address:       m.Property[mlsNumber].address,
			Bathrooms:     mls.bathrooms,
			Bedrooms:      mls.bedrooms,
			LandSize:      mls.landSize,
			MlsId:         mls.mlsID,
			MlsNumber:     mlsNumber,
			MlsUrl:        mls.mlsURL,
			Parking:       mls.parking,
			PhotoUrl:      m.Photo[mlsNumber].photoURL,
			Price:         price,
			PublicRemarks: mls.publicRemark,
			Stories:       mls.stories,
			PropertyType:  mls.propertyType,
			ListTimestamp: mls.availableTimestamp,
			Source:        mls.source,
			Latitude:      m.Property[mlsNumber].latitude,
			Longitude:     m.Property[mlsNumber].longitude,
			City:          m.Property[mlsNumber].city,
			State:         m.Property[mlsNumber].state,
			Zipcode:       m.Property[mlsNumber].zipcode,
		}
		listings.Property = append(listings.Property, p)
	}
	return listings, nil
}
