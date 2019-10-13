// indexer is a spider that collects MLS data from listing sources, normalize
// the data, and serve it to other components for further analysis.
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tony-yang/realtor-tracker/indexer/collector"
)

func runCollectors() {
	// each collector will run daily
	for name, c := range collector.Collectors {
		logrus.Infof("Running the '%s' collector...", name)
		c.FetchListing()

		result, err := c.GetDB().ReadListings()
		if err != nil {
			logrus.Errorf("reading property listing failed: %v", err)
		}
		logrus.Println(result.String())
	}
}

func main() {
	logrus.Info("Indexer Main")

	for i := 0; i < 3; i++ {
		runCollectors()
	}

}
