// indexer is a spider that collects MLS data from listing sources, normalize
// the data, and serve it to other components for further analysis.
package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tony-yang/realtor-tracker/indexer/collector"
	"github.com/tony-yang/realtor-tracker/indexer/server"
)

func runCollectors() {
	// each collector will run daily
	for i := 0; i < 3; i++ {
		for name, c := range collector.Collectors {
			logrus.Infof("Running the '%s' collector...", name)
			c.FetchListing()
			time.Sleep(5 * time.Second)
		}
	}
}

func main() {
	logrus.Info("Indexer Main")
	go runCollectors()

	server.StartServer()
}
