package indexer

import (
	"strings"
	"testing"

	mlspb "github.com/tony-yang/realtor-tracker/indexer/mls"
)

func TestParseJsonResults(t *testing.T) {
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
		result, _ := ParseJsonResults(respContent)
		wanted := mlspb.Listings{
			Property: []*mlspb.Property{
				{
					Address:       "1234 street|city, province A0B1C2",
					Bathrooms:     "1",
					Bedrooms:      "3 + 0",
					LandSize:      "0X",
					MlsId:         "20552312",
					MlsNumber:     "19016318",
					MlsUrl:        "/abc/20552312/house",
					Parking:       []string{"None"},
					PhotoUrl:      []string{"https://picture/listings/high/456.jpg"},
					Price:         "$10,000",
					PublicRemarks: "HOUSE DESCRIPTION",
					Stories:       "1.5",
					PropertyType:  "House",
				},
			},
		}
		AssertStringEqual(t, result.Property[0].Address, wanted.Property[0].Address)
		AssertStringEqual(t, result.Property[0].Bathrooms, wanted.Property[0].Bathrooms)
		AssertStringEqual(t, result.Property[0].Bedrooms, wanted.Property[0].Bedrooms)
		AssertStringEqual(t, result.Property[0].LandSize, wanted.Property[0].LandSize)
		AssertStringEqual(t, result.Property[0].MlsId, wanted.Property[0].MlsId)
		AssertStringEqual(t, result.Property[0].MlsNumber, wanted.Property[0].MlsNumber)
		AssertStringEqual(t, result.Property[0].MlsUrl, wanted.Property[0].MlsUrl)
		AssertArrayEqual(t, result.Property[0].Parking, wanted.Property[0].Parking)
		AssertArrayEqual(t, result.Property[0].PhotoUrl, wanted.Property[0].PhotoUrl)
		AssertStringEqual(t, result.Property[0].Price, wanted.Property[0].Price)
		AssertStringEqual(t, result.Property[0].PublicRemarks, wanted.Property[0].PublicRemarks)
		AssertStringEqual(t, result.Property[0].Stories, wanted.Property[0].Stories)
		AssertStringEqual(t, result.Property[0].PropertyType, wanted.Property[0].PropertyType)
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
		result, _ := ParseJsonResults(respContent)
		wanted := mlspb.Listings{
			Property: []*mlspb.Property{
				{
					Address:       "1234 street|city, province A0B1C2",
					Bathrooms:     "1",
					Bedrooms:      "3 + 0",
					LandSize:      "0X",
					MlsId:         "20552312",
					MlsNumber:     "19016318",
					MlsUrl:        "/abc/20552312/house",
					Parking:       []string{"None"},
					PhotoUrl:      []string{"https://picture/listings/high/456.jpg"},
					Price:         "$10,000",
					PublicRemarks: "HOUSE DESCRIPTION",
					Stories:       "1.5",
					PropertyType:  "Single Family",
				},
			},
		}
		AssertStringEqual(t, result.Property[0].Address, wanted.Property[0].Address)
		AssertStringEqual(t, result.Property[0].Bathrooms, wanted.Property[0].Bathrooms)
		AssertStringEqual(t, result.Property[0].Bedrooms, wanted.Property[0].Bedrooms)
		AssertStringEqual(t, result.Property[0].LandSize, wanted.Property[0].LandSize)
		AssertStringEqual(t, result.Property[0].MlsId, wanted.Property[0].MlsId)
		AssertStringEqual(t, result.Property[0].MlsNumber, wanted.Property[0].MlsNumber)
		AssertStringEqual(t, result.Property[0].MlsUrl, wanted.Property[0].MlsUrl)
		AssertArrayEqual(t, result.Property[0].Parking, wanted.Property[0].Parking)
		AssertArrayEqual(t, result.Property[0].PhotoUrl, wanted.Property[0].PhotoUrl)
		AssertStringEqual(t, result.Property[0].Price, wanted.Property[0].Price)
		AssertStringEqual(t, result.Property[0].PublicRemarks, wanted.Property[0].PublicRemarks)
		AssertStringEqual(t, result.Property[0].Stories, wanted.Property[0].Stories)
		AssertStringEqual(t, result.Property[0].PropertyType, wanted.Property[0].PropertyType)
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
		result, _ := ParseJsonResults(respContent)
		wanted := mlspb.Listings{
			Property: []*mlspb.Property{
				{
					Address:       "1234 street|city, province A0B1C2",
					Bathrooms:     "1",
					Bedrooms:      "3 + 0",
					LandSize:      "",
					MlsId:         "20552312",
					MlsNumber:     "19016318",
					MlsUrl:        "/abc/20552312/house",
					Parking:       []string{""},
					PhotoUrl:      []string{"https://picture/listings/low/456.jpg"},
					Price:         "",
					PublicRemarks: "HOUSE DESCRIPTION",
					Stories:       "",
					PropertyType:  "Single Family",
					ListTimestamp: "123456789",
				},
			},
		}
		AssertStringEqual(t, result.Property[0].Address, wanted.Property[0].Address)
		AssertStringEqual(t, result.Property[0].Bathrooms, wanted.Property[0].Bathrooms)
		AssertStringEqual(t, result.Property[0].Bedrooms, wanted.Property[0].Bedrooms)
		AssertStringEqual(t, result.Property[0].LandSize, wanted.Property[0].LandSize)
		AssertStringEqual(t, result.Property[0].MlsId, wanted.Property[0].MlsId)
		AssertStringEqual(t, result.Property[0].MlsNumber, wanted.Property[0].MlsNumber)
		AssertStringEqual(t, result.Property[0].MlsUrl, wanted.Property[0].MlsUrl)
		AssertArrayEqual(t, result.Property[0].Parking, wanted.Property[0].Parking)
		AssertArrayEqual(t, result.Property[0].PhotoUrl, wanted.Property[0].PhotoUrl)
		AssertStringEqual(t, result.Property[0].Price, wanted.Property[0].Price)
		AssertStringEqual(t, result.Property[0].PublicRemarks, wanted.Property[0].PublicRemarks)
		AssertStringEqual(t, result.Property[0].Stories, wanted.Property[0].Stories)
		AssertStringEqual(t, result.Property[0].PropertyType, wanted.Property[0].PropertyType)
	})
}

func AssertStringEqual(t *testing.T, got, want string) {
	t.Helper()
	if strings.TrimSpace(got) != strings.TrimSpace(want) {
		t.Errorf("got string '%s', want '%s'", strings.TrimSpace(got), strings.TrimSpace(want))
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
