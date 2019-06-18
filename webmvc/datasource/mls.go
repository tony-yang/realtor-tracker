package datasource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"webmvc/base"
	"webmvc/models"
)

func ParseJsonResults(responseBody []byte) ([]models.Listing, error) {
	var listingResults []models.Listing

	var resp interface{}
	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		base.Error("parseJson result err:", err)
	}

	if listings, ok := resp.(map[string]interface{})["Results"]; ok {
		for _, listing := range listings.([]interface{}) {
			listingDetail := listing.(map[string]interface{})

			var v interface{}

			property := listingDetail["Property"].(map[string]interface{})
			propertyAddress := (property["Address"].(map[string]interface{}))["AddressText"].(string)
			building := listingDetail["Building"].(map[string]interface{})
			bathroom := building["BathroomTotal"].(string)
			bedroom := building["Bedrooms"].(string)
			mlsId := (listingDetail["Id"]).(string)
			mlsNum := (listingDetail["MlsNumber"]).(string)

			lots := "NA"
			if v, ok = (listingDetail["Land"].(map[string]interface{}))["SizeTotal"]; ok {
				lots = v.(string)
			}

			url := "NA"
			if v, ok = listingDetail["RelativeDetailsURL"]; ok {
				url = strings.Replace(v.(string), "\\", "", -1)
			}

			parkingAvailable := "NA"
			if v, ok = ((property["Parking"].([]interface{}))[0].(map[string]interface{}))["Name"]; ok {
				parkingAvailable = v.(string)
			}

			photo := "NA"
			if v, ok = ((property["Photo"].([]interface{}))[0].(map[string]interface{}))["HighResPath"]; ok {
				photo = strings.Replace(v.(string), "\\", "", -1)
			} else if v, ok = ((property["Photo"].([]interface{}))[0].(map[string]interface{}))["MedResPath"]; ok {
				photo = strings.Replace(v.(string), "\\", "", -1)
			} else if v, ok = ((property["Photo"].([]interface{}))[0].(map[string]interface{}))["LowResPath"]; ok {
				photo = strings.Replace(v.(string), "\\", "", -1)
			}

			propertyPrice := "NA"
			if v, ok = property["Price"]; ok {
				propertyPrice = v.(string)
			}

			remarks := "NA"
			if v, ok = listingDetail["PublicRemarks"]; ok {
				remarks = v.(string)
			}

			numOfStories := "NA"
			if v, ok = building["StoriesTotal"]; ok {
				numOfStories = v.(string)
			}

			houseType := "NA"
			if v, ok = building["Type"]; ok {
				houseType = v.(string)
			} else if v, ok = property["Type"]; ok {
				houseType = v.(string)
			}

			house := models.Listing{
				Address:       propertyAddress,
				Bathrooms:     bathroom,
				Bedrooms:      bedroom,
				LandSize:      lots,
				MlsId:         mlsId,
				MlsNumber:     mlsNum,
				MlsUrl:        url,
				Parking:       parkingAvailable,
				PhotoUrl:      photo,
				Price:         propertyPrice,
				PublicRemarks: remarks,
				Stories:       numOfStories,
				PropertyType:  houseType,
				ListTimestamp: "123456789",
			}
			listingResults = append(listingResults, house)
		}
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

	bodyContent := `{
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

	result, err := ParseJsonResults(bodyContent)
	if err != nil {
		base.Error("failed to parse the json response into useful listing format:", err)
		statusCode = http.StatusInternalServerError
	}

	listings, err := json.Marshal(result)
	if err != nil {
		base.Error("failed to create json listing:", err)
		statusCode = http.StatusInternalServerError
	}

	return string(listings), statusCode, err
}
