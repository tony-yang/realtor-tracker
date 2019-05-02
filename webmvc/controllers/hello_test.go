package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"webmvc"
	"webmvc/controllers"
	"webmvc/tester"
)

func TestHello(t *testing.T) {
	server := webmvc.CreateNewServer()
	server.Routes.RegisterRoute("/hello", &controllers.Hello{})

	t.Run("return Hello World from hello", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/hello", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "Hello World!"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, http.StatusOK)
	})

	t.Run("Post to hello returns POST message", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/hello", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "Post to Hello World!"
		tester.AssertStringEqual(t, got, want)
		tester.AssertIntEqual(t, response.Code, http.StatusOK)
	})
}
