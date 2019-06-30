package main

import (
	"net/http"

	"github.com/tony-yang/realtor-tracker/webmvc"
	"github.com/tony-yang/realtor-tracker/webmvc/base"
)

func main() {
	base.Debug("Starting the WebMVC Go Framework")
	addr := ":80"

	server := webmvc.CreateNewServer()
	ConfigRoutes(server)

	if err := http.ListenAndServe(addr, server); err != nil {
		base.Critical("The WebMVC Go Framework failed on port 80:", err)
	}
}
