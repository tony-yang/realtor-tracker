package storage

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
)

var (
	mlsDB *memoryDB
)

type listingStatus int

const (
	Open listingStatus = iota
	Pending
	Sold
	Closed
)

type mls struct {
	mlsID              int
	mlsURL             string
	bathrooms          string
	bedrooms           string
	landSize           string
	parking            string
	publicRemark       string
	stories            string
	propertyType       string
	availableTimestamp int
	status             listingStatus
}

type property struct {
	address string
}

type photo struct {
	photoURL string
}

type priceHistory struct {
	price     int
	timestamp int64
}

type memoryDB struct {
	lock         sync.Mutex
	mls          map[string]*mls
	property     map[string]*property
	photo        map[string][]*photo
	priceHistory map[string][]*priceHistory
}

func newMemoryDB() *memoryDB {
	return &memoryDB{
		mls:          make(map[string]*mls),
		property:     make(map[string]*property),
		photo:        make(map[string][]*photo),
		priceHistory: make(map[string][]*priceHistory),
	}
}

func (m *memoryDB) CreateStorage() *memoryDB {
	return m
}

func (m *memoryDB) SaveNewListing(listings map[string]*mlspb.Property) error {
	for mlsNumber, p := range listings {
		logrus.Info("######### mlsNumber =", mlsNumber, "listing", p)
		id, err := strconv.Atoi(p.MlsId)
		if err != nil {
			return err
		}
		timestamp, err := strconv.Atoi(p.ListTimestamp)
		if err != nil {
			return err
		}
		m.mls[mlsNumber] = &mls{
			mlsID:              id,
			mlsURL:             p.MlsUrl,
			bathrooms:          p.Bathrooms,
			bedrooms:           p.Bedrooms,
			landSize:           p.LandSize,
			parking:            p.Parking,
			publicRemark:       p.PublicRemarks,
			stories:            p.Stories,
			propertyType:       p.PropertyType,
			availableTimestamp: timestamp,
			status:             Open,
		}
		m.property[mlsNumber] = &property{address: p.Address}
		m.photo[mlsNumber] = []*photo{{photoURL: p.PhotoUrl}}
		price, err := strconv.Atoi(strings.ReplaceAll(strings.TrimLeft(p.Price, "$"), ",", ""))
		if err != nil {
			return err
		}
		m.priceHistory[mlsNumber] = []*priceHistory{
			{
				price:     price,
				timestamp: time.Now().Unix(),
			},
		}
	}
	return nil
}

func (m *memoryDB) ReadListing(id string) (string, error) {
	return "", nil
}

func (m *memoryDB) ReadListings() (string, error) {
	return "", nil
}
