package controllers

import (
	"encoding/json"
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
	bodyContent, err := l.ReadListing(subpath)
	statusCode := http.StatusOK
	if err != nil {
		base.Error("error fetch mls listing:", err)
		statusCode = http.StatusInternalServerError
	}

	body, err := json.Marshal(bodyContent)
	if err != nil {
		body = []byte{}
		statusCode = http.StatusInternalServerError
	}

	return &base.HttpResponse{
		Body:       string(body),
		StatusCode: statusCode,
	}
}
