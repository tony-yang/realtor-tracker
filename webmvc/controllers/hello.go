package controllers

import (
	"webmvc/base"
)

type Hello struct {
	base.Controller
}

func (h *Hello) Get(subpath string, queries map[string]string) *base.HttpResponse {
	return &base.HttpResponse{
		Body:       "Hello World!",
		StatusCode: 200,
	}
}

func (h *Hello) Post(subpath string, queries map[string]string) *base.HttpResponse {
	return &base.HttpResponse{
		Body:       "Post to Hello World!",
		StatusCode: 200,
	}
}
