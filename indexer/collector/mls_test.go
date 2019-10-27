package collector

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
	"github.com/tony-yang/realtor-tracker/indexer/storage"
)

type roundTripFunc func(r *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r), nil
}

func NewTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestFetchListing(t *testing.T) {
	respContent := `{
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
	}`
	t.Run("can format and save listing", func(t *testing.T) {
		c := NewTestClient(func(r *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(respContent)),
				Header:     make(http.Header),
			}
		})
		cityIndex := map[string]*storage.City{
			"city,province": {
				Name:      "City",
				State:     "Province",
				MlsNumber: make(map[string]bool),
			},
		}
		mDB, _ := storage.NewMemoryDB(cityIndex)
		m := NewMls(mDB, c)
		m.FetchListing()
		savedListings, _ := mDB.ReadListings()

		if savedListings.Property[0].MlsId != "20552312" {
			t.Errorf("mlsID incorrectly saved, expected %s, got %s", "20552312", savedListings.Property[0].MlsId)
		}
	})
}

func TestFormatListing(t *testing.T) {
	t.Run("can parse result properly", func(t *testing.T) {
		respContent := []byte(`{
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
        "RelativeDetailsURL": "\/abc\/20552312\/house",
        "StatusId": "1",
        "PhotoChangeDateUTC": "2019-05-04 12:34:56 PM",
        "Distance": "",
        "RelativeURLEn": "\/abc\/20552312\/house",
        "RelativeURLFr": "\/abc\/20552312\/house"
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
		var listings *listings
		err := json.Unmarshal(respContent, &listings)
		if err != nil {
			t.Errorf("failed to parse the json response into listing: %v", err)
		}
		result := formatListing(listings)
		mlsNumber := "19016318"
		price := result[mlsNumber].Price
		wanted := map[string]*mlspb.Property{
			mlsNumber: {
				Address:       "1234 street|city, province A0B1C2",
				Bathrooms:     "1",
				Bedrooms:      "3 + 0",
				LandSize:      "0X",
				MlsId:         "20552312",
				MlsNumber:     mlsNumber,
				MlsUrl:        "/abc/20552312/house",
				Parking:       []string{"None"},
				PhotoUrl:      []string{"https://picture/listings/high/456.jpg"},
				Price:         price,
				PublicRemarks: "HOUSE DESCRIPTION",
				Stories:       "1.5",
				PropertyType:  "House",
				Latitude:      98.765432,
				Longitude:     -12.345678,
				City:          "city",
				State:         "province",
				Zipcode:       "A0B1C2",
			},
		}
		AssertStringEqual(t, result[mlsNumber].Address, wanted[mlsNumber].Address)
		AssertStringEqual(t, result[mlsNumber].Bathrooms, wanted[mlsNumber].Bathrooms)
		AssertStringEqual(t, result[mlsNumber].Bedrooms, wanted[mlsNumber].Bedrooms)
		AssertStringEqual(t, result[mlsNumber].LandSize, wanted[mlsNumber].LandSize)
		AssertStringEqual(t, result[mlsNumber].MlsId, wanted[mlsNumber].MlsId)
		AssertStringEqual(t, result[mlsNumber].MlsNumber, wanted[mlsNumber].MlsNumber)
		AssertStringEqual(t, result[mlsNumber].MlsUrl, wanted[mlsNumber].MlsUrl)
		AssertArrayEqual(t, result[mlsNumber].Parking, wanted[mlsNumber].Parking)
		AssertArrayEqual(t, result[mlsNumber].PhotoUrl, wanted[mlsNumber].PhotoUrl)
		AssertStringEqual(t, result[mlsNumber].PublicRemarks, wanted[mlsNumber].PublicRemarks)
		AssertStringEqual(t, result[mlsNumber].Stories, wanted[mlsNumber].Stories)
		AssertStringEqual(t, result[mlsNumber].PropertyType, wanted[mlsNumber].PropertyType)
		AssertFloat64Equal(t, result[mlsNumber].Latitude, wanted[mlsNumber].Latitude)
		AssertFloat64Equal(t, result[mlsNumber].Longitude, wanted[mlsNumber].Longitude)
		AssertStringEqual(t, result[mlsNumber].City, wanted[mlsNumber].City)
		AssertStringEqual(t, result[mlsNumber].State, wanted[mlsNumber].State)
		AssertStringEqual(t, result[mlsNumber].Zipcode, wanted[mlsNumber].Zipcode)
	})

	t.Run("can parse result when optional properties not available", func(t *testing.T) {
		respContent := []byte(`{
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
          "StoriesTotal": "1.5"
        },
        "Property": {
          "Price": "$10,000",
          "Type": "Single Family",
          "Address": {
            "AddressText": "1234 street|city, province A0B1C2"
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
        "RelativeDetailsURL": "\/abc\/20552312\/house",
        "StatusId": "1",
        "PhotoChangeDateUTC": "2019-05-04 12:34:56 PM",
        "Distance": "",
        "RelativeURLEn": "\/abc\/20552312\/house",
        "RelativeURLFr": "\/abc\/20552312\/house"
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
		var listings *listings
		err := json.Unmarshal(respContent, &listings)
		if err != nil {
			t.Errorf("failed to parse the json response into listing: %v", err)
		}
		result := formatListing(listings)
		mlsNumber := "19016318"
		price := result[mlsNumber].Price
		wanted := map[string]*mlspb.Property{
			mlsNumber: {
				Address:       "1234 street|city, province A0B1C2",
				Bathrooms:     "1",
				Bedrooms:      "3 + 0",
				LandSize:      "0X",
				MlsId:         "20552312",
				MlsNumber:     mlsNumber,
				MlsUrl:        "/abc/20552312/house",
				Parking:       []string{"None"},
				PhotoUrl:      []string{"https://picture/listings/high/456.jpg"},
				Price:         price,
				PublicRemarks: "HOUSE DESCRIPTION",
				Stories:       "1.5",
				PropertyType:  "Single Family",
				City:          "city",
				State:         "province",
				Zipcode:       "A0B1C2",
			},
		}
		AssertStringEqual(t, result[mlsNumber].Address, wanted[mlsNumber].Address)
		AssertStringEqual(t, result[mlsNumber].Bathrooms, wanted[mlsNumber].Bathrooms)
		AssertStringEqual(t, result[mlsNumber].Bedrooms, wanted[mlsNumber].Bedrooms)
		AssertStringEqual(t, result[mlsNumber].LandSize, wanted[mlsNumber].LandSize)
		AssertStringEqual(t, result[mlsNumber].MlsId, wanted[mlsNumber].MlsId)
		AssertStringEqual(t, result[mlsNumber].MlsNumber, wanted[mlsNumber].MlsNumber)
		AssertStringEqual(t, result[mlsNumber].MlsUrl, wanted[mlsNumber].MlsUrl)
		AssertArrayEqual(t, result[mlsNumber].Parking, wanted[mlsNumber].Parking)
		AssertArrayEqual(t, result[mlsNumber].PhotoUrl, wanted[mlsNumber].PhotoUrl)
		AssertStringEqual(t, result[mlsNumber].PublicRemarks, wanted[mlsNumber].PublicRemarks)
		AssertStringEqual(t, result[mlsNumber].Stories, wanted[mlsNumber].Stories)
		AssertStringEqual(t, result[mlsNumber].PropertyType, wanted[mlsNumber].PropertyType)
		AssertFloat64Equal(t, result[mlsNumber].Latitude, wanted[mlsNumber].Latitude)
		AssertFloat64Equal(t, result[mlsNumber].Longitude, wanted[mlsNumber].Longitude)
		AssertStringEqual(t, result[mlsNumber].City, wanted[mlsNumber].City)
		AssertStringEqual(t, result[mlsNumber].State, wanted[mlsNumber].State)
		AssertStringEqual(t, result[mlsNumber].Zipcode, wanted[mlsNumber].Zipcode)
	})

	t.Run("can parse result when missing building type but have property type", func(t *testing.T) {
		respContent := []byte(`{
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
          "Bedrooms": "3 + 0"
        },
        "Property": {
          "Type": "Single Family",
          "Address": {
            "AddressText": "1234 street|city, province A0B1C2"
          },
          "Photo": [{
            "SequenceId": "1",
            "LowResPath": "https:\/\/picture\/listings\/low\/456.jpg",
            "LastUpdated": "2019-05-04 12:34:56 PM"
          }],
          "Parking": [{}],
          "TypeId": "300",
          "OwnershipType": "Freehold"
        },
        "Business": {},
        "Land": {},
        "PostalCode": "A0B1C2",
        "StatusId": "1",
        "PhotoChangeDateUTC": "2019-05-04 12:34:56 PM",
        "Distance": "",
        "RelativeURLEn": "\/abc\/20552312\/house",
        "RelativeURLFr": "\/abc\/20552312\/house"
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
		var listings *listings
		err := json.Unmarshal(respContent, &listings)
		if err != nil {
			t.Errorf("failed to parse the json response into listing: %v", err)
		}
		result := formatListing(listings)
		mlsNumber := "19016318"
		price := result[mlsNumber].Price
		wanted := map[string]*mlspb.Property{
			mlsNumber: {
				Address:       "1234 street|city, province A0B1C2",
				Bathrooms:     "1",
				Bedrooms:      "3 + 0",
				LandSize:      "",
				MlsId:         "20552312",
				MlsNumber:     mlsNumber,
				MlsUrl:        "/abc/20552312/house",
				Parking:       []string{""},
				PhotoUrl:      []string{"https://picture/listings/low/456.jpg"},
				Price:         price,
				PublicRemarks: "HOUSE DESCRIPTION",
				Stories:       "",
				PropertyType:  "Single Family",
				ListTimestamp: 123456789,
				City:          "city",
				State:         "province",
				Zipcode:       "A0B1C2",
			},
		}
		AssertStringEqual(t, result[mlsNumber].Address, wanted[mlsNumber].Address)
		AssertStringEqual(t, result[mlsNumber].Bathrooms, wanted[mlsNumber].Bathrooms)
		AssertStringEqual(t, result[mlsNumber].Bedrooms, wanted[mlsNumber].Bedrooms)
		AssertStringEqual(t, result[mlsNumber].LandSize, wanted[mlsNumber].LandSize)
		AssertStringEqual(t, result[mlsNumber].MlsId, wanted[mlsNumber].MlsId)
		AssertStringEqual(t, result[mlsNumber].MlsNumber, wanted[mlsNumber].MlsNumber)
		AssertStringEqual(t, result[mlsNumber].MlsUrl, wanted[mlsNumber].MlsUrl)
		AssertArrayEqual(t, result[mlsNumber].Parking, wanted[mlsNumber].Parking)
		AssertArrayEqual(t, result[mlsNumber].PhotoUrl, wanted[mlsNumber].PhotoUrl)
		AssertStringEqual(t, result[mlsNumber].PublicRemarks, wanted[mlsNumber].PublicRemarks)
		AssertStringEqual(t, result[mlsNumber].Stories, wanted[mlsNumber].Stories)
		AssertStringEqual(t, result[mlsNumber].PropertyType, wanted[mlsNumber].PropertyType)
		AssertStringEqual(t, result[mlsNumber].City, wanted[mlsNumber].City)
		AssertStringEqual(t, result[mlsNumber].State, wanted[mlsNumber].State)
		AssertStringEqual(t, result[mlsNumber].Zipcode, wanted[mlsNumber].Zipcode)
	})
}

func AssertStringEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got string '%s', want '%s'", got, want)
	}
}

func AssertFloat64Equal(t *testing.T, got, want float64) {
	t.Helper()
	if got != want {
		t.Errorf("got string '%f', want '%f'", got, want)
	}
}

func AssertArrayEqual(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("got array '%v', want array '%v'", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got array '%v', want array '%v'", got, want)
		}
	}
}
