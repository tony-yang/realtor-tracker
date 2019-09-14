package storage

import (
	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
)

type StorageInterfac interface {
	CreateStorage() error
	SaveNewListing(listings map[string]*mlspb.Property) error
	ReadListing(id string) (string, error)
	ReadListings() (string, error)
}
