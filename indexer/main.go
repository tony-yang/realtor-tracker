// indexer is a spider that collects MLS data from listing sources, normalize
// the data, and serve it to other components for further analysis.
package main

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tony-yang/realtor-tracker/indexer/collector"
)

func runCollectors(wg *sync.WaitGroup) {
	defer wg.Done()
	for name, c := range collector.Collectors {
		logrus.Infof("Running the %q collector...", name)
		c.FetchListing()
		time.Sleep(5 * time.Second)
		logrus.Infof("%q finished collection.", name)
	}
}

func main() {
	var wg sync.WaitGroup

	logrus.Info("Indexer Main")
	wg.Add(1)
	go runCollectors(&wg)
	wg.Wait()
	logrus.Info("Indexer collection cycle finished successfully.")
}
