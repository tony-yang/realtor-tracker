// package collector fetches MLS data from source as implemented by individual
// mls collector.
package collector

import (
	"fmt"

	"github.com/tony-yang/realtor-tracker/indexer/storage"
)

var (
	Collectors = make(map[string]Collector)
)

// RegisterCollector registers the individual mls collector for each MLS data source.
func RegisterCollector(name string, c Collector) error {
	if _, ok := Collectors[name]; ok {
		return fmt.Errorf("Collector '%s' already existed. Register the new collector with a different name", name)
	}
	Collectors[name] = c
	return nil
}

type building struct {
	Bathrooms    string `json:"BathroomTotal"`
	Bedrooms     string `json:"Bedrooms"`
	Stories      string `json:"StoriesTotal"`
	BuildingType string `json:"Type"`
}

type address struct {
	Address   string `json:"AddressText"`
	Latitude  string `json:"Latitude"`
	Longitude string `json:"Longitude"`
}

type photo struct {
	SequenceId  string `json:"SequenceId"`
	HighRes     string `json:"HighResPath"`
	MedRes      string `json:"MedResPath"`
	LowRes      string `json:"LowResPath"`
	LastUpdated string `json:"LastUpdated"`
}

type parking struct {
	Name string `json:"Name"`
}

type property struct {
	Price        string    `json:"Price"`
	PropertyType string    `json:"Type"`
	Address      address   `json:"Address"`
	Photos       []photo   `json:"Photo"`
	Parkings     []parking `json:"Parking"`
}

type land struct {
	Size string `json:"SizeTotal"`
}

type listing struct {
	ID            string   `json:"Id"`
	MlsNumber     string   `json:"MlsNumber"`
	PublicRemarks string   `json:"PublicRemarks"`
	Building      building `json:"Building"`
	Property      property `json:"Property"`
	Land          land     `json:"Land"`
	ZipCode       string   `json:"PostalCode"`
	URL           string   `json:"RelativeDetailsURL"`
	URLEn         string   `json:"RelativeURLEn"`
}

type listings struct {
	Listing []listing `json:"Results"`
}

// Collector defines the interface for individual collector implementation.
type Collector interface {
	// FetchListing retrieves from source listing and saves to DB
	FetchListing()
	// GetDB retrieves the DB instance
	GetDB() storage.DBInterface
}
