package main

import (
	"net/http"

	"github.com/tony-yang/realtor-tracker/webmvc/base"
	"github.com/tony-yang/realtor-tracker/webmvc/server"
)

func main() {
	base.Debug("Starting the WebMVC Go Framework")
	addr := ":80"

	s := server.CreateNewServer()
	ConfigRoutes(s)

	if err := http.ListenAndServe(addr, s); err != nil {
		base.Critical("The WebMVC Go Framework failed on port 80:", err)
	}
}
