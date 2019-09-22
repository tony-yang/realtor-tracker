package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
	"github.com/tony-yang/realtor-tracker/indexer/storage"

	"github.com/sirupsen/logrus"
)

const (
	source = "mls-canada"
)

func init() {
	RegisterCollector(source, NewMls(storage.NewMemoryDB(), &http.Client{}))
}

type mls struct {
	db     storage.DBInterface
	client *http.Client
}

// NewMls create a new client for the MLS Canada collector.
func NewMls(s storage.DBInterface, c *http.Client) *mls {
	if s == nil {
		s = storage.NewMemoryDB()
	}

	if c == nil {
		c = &http.Client{}
	}

	return &mls{
		db:     s,
		client: c,
	}
}

func formatListing(listings *listings) map[string]*mlspb.Property {
	properties := make(map[string]*mlspb.Property)
	for _, l := range listings.Listing {
		parkings := []string{}
		photos := []string{}
		for _, p := range l.Property.Parkings {
			parkings = append(parkings, strings.TrimSpace(p.Name))
		}
		for _, photo := range l.Property.Photos {
			if photo.HighRes != "" {
				photos = append(photos, strings.TrimSpace(photo.HighRes))
			} else if photo.MedRes != "" {
				photos = append(photos, strings.TrimSpace(photo.MedRes))
			} else if photo.LowRes != "" {
				photos = append(photos, strings.TrimSpace(photo.LowRes))
			}
		}

		mlsURL := strings.TrimSpace(l.URL)
		if mlsURL == "" {
			mlsURL = strings.TrimSpace(l.URLEn)
		}

		houseType := strings.TrimSpace(l.Building.BuildingType)
		if houseType == "" {
			houseType = strings.TrimSpace(l.Property.PropertyType)
		}
		p, err := strconv.Atoi(strings.ReplaceAll(strings.TrimLeft(l.Property.Price, "$"), ",", ""))
		if err != nil {
			p = -1
		}
		price := []*mlspb.PriceHistory{
			{
				Price:     int32(p),
				Timestamp: time.Now().Unix(),
			},
		}

		house := &mlspb.Property{
			Address:       strings.TrimSpace(l.Property.Address.Address),
			Bathrooms:     strings.TrimSpace(l.Building.Bathrooms),
			Bedrooms:      strings.TrimSpace(l.Building.Bedrooms),
			LandSize:      strings.TrimSpace(l.Land.Size),
			MlsId:         strings.TrimSpace(l.ID),
			MlsNumber:     strings.TrimSpace(l.MlsNumber),
			MlsUrl:        mlsURL,
			Parking:       parkings,
			PhotoUrl:      photos,
			Price:         price,
			PublicRemarks: strings.TrimSpace(l.PublicRemarks),
			Stories:       strings.TrimSpace(l.Building.Stories),
			PropertyType:  houseType,
			ListTimestamp: 123456789,
			Source:        source,
		}
		properties[house.MlsNumber] = house
	}
	return properties
}

// FetchListing retrieves the mls listing from MLS Canada.
func (m *mls) FetchListing() {
	listingUrl := "https://api2.realtor.ca/Listing.svc/PropertySearch_Post"
	data := url.Values{
		"ZoomLevel":            {"11"},
		"LatitudeMax":          {"42.3661983"},
		"LongitudeMax":         {"-82.4784635"},
		"LatitudeMin":          {"41.9947561"},
		"LongitudeMin":         {"-83.1245969"},
		"CurrentPage":          {"1"},
		"Sort":                 {"1-A"},
		"RecordsPerPage":       {"2"},
		"PropertyTypeGroupID":  {"1"},
		"PropertySearchTypeId": {"1"},
		"TransactionTypeId":    {"2"},
		"ApplicationId":        {"1"},
		"CultureId":            {"1"},
		"Version":              {"7.0"},
	}

	resp, err := m.client.PostForm(listingUrl, data)
	if err != nil {
		logrus.Error("http post form error:", err)
	}

	defer resp.Body.Close()
	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("failed to read response:", err)
	}

	var listings *listings
	err = json.Unmarshal(bodyContent, &listings)
	if err != nil {
		logrus.Fatalf("failed to parse the json response into listing: %v", err)
	}

	properties := formatListing(listings)
	m.db.SaveNewListing(properties)
}
