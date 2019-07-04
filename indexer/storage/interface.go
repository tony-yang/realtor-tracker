package storage

type StorageInterfac interface {
  CreateStorage() error
  SaveNewListing(listings map[string]Listing) error
  ReadListing(id string) (string, error)
  ReadListings() (string, error)
}
