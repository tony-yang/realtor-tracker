// web is the server entry point to the MLS data collected by the Collectors.
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tony-yang/realtor-tracker/web/server"
)

func main() {
	logrus.Info("Web Server Main")
	server.StartServer()
}
