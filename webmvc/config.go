package main

import (
	"github.com/tony-yang/realtor-tracker/webmvc/controllers"
	"github.com/tony-yang/realtor-tracker/webmvc/server"
)

// ConfigRoutes configures the routes with the corresponding controller
func ConfigRoutes(s *server.NewServer) {
	s.Routes.RegisterRoute("/index", &controllers.Index{})
	s.Routes.RegisterRoute("/hello", &controllers.Hello{})
	s.Routes.RegisterRoute("/listings", &controllers.Listing{})
}
