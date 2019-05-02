package main

import (
	"webmvc"
	"webmvc/controllers"
)

// ConfigRoutes configures the routes with the corresponding controller
func ConfigRoutes(s *webmvc.NewServer) {
	s.Routes.RegisterRoute("/index", &controllers.Index{})
	s.Routes.RegisterRoute("/hello", &controllers.Hello{})
}
