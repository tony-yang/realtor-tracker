package webmvc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"webmvc"
	"webmvc/base"
	"webmvc/tester"
)

func TestServer(t *testing.T) {
	server := webmvc.CreateNewServer()

	t.Run("GET base controller with non-pre-defined HTTP code should return method not allowed", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/base", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "Not Found"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, 404)
	})

	t.Run("GET base controller should return 405 with method not allowed", func(t *testing.T) {
		server.Routes.RegisterRoute("/base", &base.Controller{})

		request, _ := http.NewRequest(http.MethodGet, "/base", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "GET Method Not Allowed"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, 405)
	})

	t.Run("POST base controller should return 405 with method not allowed", func(t *testing.T) {
		server.Routes.RegisterRoute("/base", &base.Controller{})

		request, _ := http.NewRequest(http.MethodPost, "/base", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "POST Method Not Allowed"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, 405)
	})

	t.Run("PUT base controller should return 405 with method not allowed", func(t *testing.T) {
		server.Routes.RegisterRoute("/base", &base.Controller{})

		request, _ := http.NewRequest(http.MethodPut, "/base", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "PUT Method Not Allowed"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, 405)
	})

	t.Run("DELETE base controller should return 405 with method not allowed", func(t *testing.T) {
		server.Routes.RegisterRoute("/base", &base.Controller{})

		request, _ := http.NewRequest(http.MethodDelete, "/base", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "DELETE Method Not Allowed"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, 405)
	})

	t.Run("Unknown Method base controller should return 405 with method not allowed", func(t *testing.T) {
		server.Routes.RegisterRoute("/base", &base.Controller{})

		request, _ := http.NewRequest("UNKNOWN", "/base", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "GET Method Not Allowed"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, 405)
	})
}
