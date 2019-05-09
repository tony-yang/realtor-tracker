package controllers

import (
	"fmt"
	"net/http"
	"webmvc/base"
	"webmvc/models"
)

type Listing struct {
	base.Controller
	models.Listing
}

func (l *Listing) Get(subpath string, queries map[string]string) *base.HttpResponse {
	fmt.Println("subpath =", subpath, "queries =", queries)
	bodyContent, statusCode, err := l.FetchMlsListing(subpath)
	if err != nil {
		base.Error("error fetch mls listing:", err)
		if statusCode != http.StatusInternalServerError {
			statusCode = http.StatusInternalServerError
		}
	}

	response := &base.HttpResponse{
		Body:       bodyContent,
		StatusCode: statusCode,
	}
	return response
}
