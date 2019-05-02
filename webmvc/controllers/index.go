package controllers

import (
	"fmt"
	"webmvc/base"
)

type Index struct {
	base.Controller
}

func (index *Index) Get(subpath string, queries map[string]string) *base.HttpResponse {
	var bodyContent string
	if val, ok := queries["query"]; ok {
		if val2, ok2 := queries["query2"]; ok2 {
			bodyContent = fmt.Sprintf("Hello Webmvc from Go %s with param = %s and %s", subpath, val, val2)
		} else {
			bodyContent = fmt.Sprintf("Hello Webmvc from Go %s with param = %s", subpath, val)
		}
	} else {
		bodyContent = fmt.Sprintf("Hello Webmvc from Go %s", subpath)
	}

	response := &base.HttpResponse{
		Body:       bodyContent,
		StatusCode: 200,
	}
	return response
}
