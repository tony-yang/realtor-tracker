package indexer

import (
	"encoding/json"
	"net/http"
	"strings"

	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"

	"github.com/sirupsen/logrus"
)

type building struct {
	Bathrooms    string `json:"BathroomTotal"`
	Bedrooms     string `json:"Bedrooms"`
	Stories      string `json:"StoriesTotal"`
	BuildingType string `json:"Type"`
}

type address struct {
	Address string `json:"AddressText"`
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

func ParseJsonResults(responseBody []byte) (*mlspb.Listings, error) {
	listingResults := &mlspb.Listings{}

	var listings listings
	err := json.Unmarshal(responseBody, &listings)
	if err != nil {
		logrus.Errorf("cannot parse the mls listing json: %v", err)
		return listingResults, err
	}

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
			Price:         strings.TrimSpace(l.Property.Price),
			PublicRemarks: strings.TrimSpace(l.PublicRemarks),
			Stories:       strings.TrimSpace(l.Building.Stories),
			PropertyType:  houseType,
			ListTimestamp: "123456789",
		}
		listingResults.Property = append(listingResults.Property, house)
	}

	return listingResults, nil
}

func FetchMlsListing() (string, int, error) {
	statusCode := http.StatusOK

	// listingUrl := "https://api2.realtor.ca/Listing.svc/PropertySearch_Post"
	// data := url.Values{
	// 	"ZoomLevel":            {"11"},
	// 	"LatitudeMax":          {"42.3661983"},
	// 	"LongitudeMax":         {"-82.4784635"},
	// 	"LatitudeMin":          {"41.9947561"},
	// 	"LongitudeMin":         {"-83.1245969"},
	// 	"CurrentPage":          {"1"},
	// 	"Sort":                 {"1-A"},
	// 	"RecordsPerPage":       {"2"},
	// 	"PropertyTypeGroupID":  {"1"},
	// 	"PropertySearchTypeId": {"1"},
	// 	"TransactionTypeId":    {"2"},
	// 	"ApplicationId":        {"1"},
	// 	"CultureId":            {"1"},
	// 	"Version":              {"7.0"},
	// }
	//
	// resp, err := http.PostForm(listingUrl, data)
	// if err != nil {
	// 	base.Error("http post form error:", err)
	// 	statusCode = http.StatusInternalServerError
	// }
	//
	// defer resp.Body.Close()
	// bodyContent, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	base.Error("failed to read response:", err)
	// 	statusCode = http.StatusInternalServerError
	// }
	//
	// fmt.Println("bodyContent =", string(bodyContent))

	bodyContent := []byte(`{
	  "ErrorCode": {
	    "Id": 200,
	    "Description": "Success (hidden)",
	    "Status": "Status",
	    "ProductName": "Name",
	    "Version": "1.0.7047.17724"
	  },
	  "Paging": {
	    "RecordsPerPage": 2,
	    "CurrentPage": 1,
	    "TotalRecords": 10,
	    "MaxRecords": 10,
	    "TotalPages": 2,
	    "RecordsShowing": 10,
	    "Pins": 1
	  },
	  "Results": [{
	    "Id": "20552312",
	    "MlsNumber": "19016318",
	    "PublicRemarks": "HOUSE DESCRIPTION",
	    "Building": {
	      "BathroomTotal": "1",
	      "Bedrooms": "3 + 0",
	      "StoriesTotal": "1.5",
	      "Type": "House"
	    },
	    "Property": {
	      "Price": "$10,000",
	      "Type": "Single Family",
	      "Address": {
	        "AddressText": "1234 street|city, province A0B1C2",
	        "Longitude": "-12.345678",
	        "Latitude": "98.765432"
	      },
	      "Photo": [{
	        "SequenceId": "1",
	        "HighResPath": "https:\/\/picture\/listings\/high\/456.jpg",
	        "MedResPath": "https:\/\/picture\/listings\/med\/456.jpg",
	        "LowResPath": "https:\/\/picture\/listings\/low\/456.jpg",
	        "LastUpdated": "2019-05-04 12:34:56 PM"
	      }],
	      "Parking": [{
	        "Name": "None"
	      }],
	      "TypeId": "300",
	      "OwnershipType": "Freehold"
	    },
	    "Business": {},
	    "Land": {
	      "SizeTotal": "0X"
	    },
	    "PostalCode": "A0B1C2",
	    "RelativeDetailsURL": "\/abc.com\/20552312\/house",
	    "StatusId": "1",
	    "PhotoChangeDateUTC": "2019-05-04 12:34:56 PM",
	    "Distance": "",
	    "RelativeURLEn": "\/abc.com\/20552312\/house",
	    "RelativeURLFr": "\/abc.com\/20552312\/house"
	  }],
	  "Pins": [{
	    "key": "L8|42.01|-82.49",
	    "propertyId": "",
	    "count": 1,
	    "longitude": "-82.49",
	    "latitude": "42.01"
	  }],
	  "GroupingLevel": "8"
	}`)

	result, err := ParseJsonResults(bodyContent)
	if err != nil {
		logrus.Fatalf("failed to parse the json response into useful listing format: %v", err)
		statusCode = http.StatusInternalServerError
	}

	listings, err := json.Marshal(result)
	if err != nil {
		logrus.Fatalf("failed to create json listing: %v", err)
		statusCode = http.StatusInternalServerError
	}

	// database, err := sql.Open("sqlite3", "/tmp/realtor.db")
	// if err != nil {
	// 	base.Error("failed to create a new database:", err)
	// 	statusCode = http.StatusInternalServerError
	// }
	// err = CreateDB(database)
	// if err != nil {
	// 	fmt.Println("create db failed", err)
	// }
	// err = SaveNewListing(database, result)
	// if err != nil {
	// 	fmt.Println("save new listing failed", err)
	// }
	//
	// _, err = ReadAllListings(database)
	// if err != nil {
	// 	fmt.Println("read all listings failed:", err)
	// }

	return string(listings), statusCode, err
}
