// package storage stores the MLS listing data collected from the collectors.
package storage

import (
	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
)

// DBInterface defines the common interface for all types of storage implemented.
type DBInterface interface {
	CreateStorage() error
	SaveNewListing(listings map[string]*mlspb.Property) error
	ReadListing(id string) (string, error)
	ReadListings() (*mlspb.Listings, error)
}
