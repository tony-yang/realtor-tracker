package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"webmvc"
	"webmvc/controllers"
	"webmvc/tester"
)

func TestIndex(t *testing.T) {
	server := webmvc.CreateNewServer()
	server.Routes.RegisterRoute("/index", &controllers.Index{})

	t.Run("return Hello Webmvc from index", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/index", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "Hello Webmvc from Go"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, http.StatusOK)
	})

	t.Run("return Hello Webmvc from index with subpath", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/index/1", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "Hello Webmvc from Go 1"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, http.StatusOK)
	})

	t.Run("return Hello Webmvc from index with param", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/index/1?query=value", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "Hello Webmvc from Go 1 with param = value"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, http.StatusOK)
	})

	t.Run("return Hello Webmvc from index with multiple params", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/index/1?query=value&query2=value2", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "Hello Webmvc from Go 1 with param = value and value2"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, http.StatusOK)
	})

	t.Run("POST to index not allowed", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/index", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "POST Method Not Allowed"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, 405)
	})
}
