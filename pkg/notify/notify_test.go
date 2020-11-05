package notify

import (
	"testing"

	"github.com/300481/pricenotifier/pkg/station"
)

type testClient struct{}

func TestUpdateBestStations(t *testing.T) {
	price := float64(1.0)
	stations := []string{"A", "B", "C"}
	fuels := []string{"E5", "Diesel"}

	n := NewNotifier(&testClient{})

	for _, fuel := range fuels {
		sm := station.NewStationMap()
		for _, s := range stations {
			sm[s] = station.NewStation(s, s, s)
			sm[s].AddPrice(1, fuel, price)
		}
		n.UpdateBestStations(fuel, price, sm)

		nStations := len(stations)

		nBestStations := len(n.CurrentBestStations[fueltype(fuel)])
		if nBestStations != nStations {
			t.Errorf("Count of best stations don't fits the right size. Got: %d must: %d", nBestStations, nStations)
		}

		nBestPrices := len(n.CurrentBestPrices[fueltype(fuel)])
		if nBestPrices != nStations {
			t.Errorf("Count of best prices don't fits the right size. Got: %d must: %d", nBestPrices, nStations)
		}
	}
}

// TestNotify tests the notifiers notify function
func TestNotify(t *testing.T) {
	testPrice := float64(1.0)
	testFuel := "Diesel"
	testStation := "A"

	n := NewNotifier(&testClient{})

	sm := station.NewStationMap()
	sm[testStation] = station.NewStation(testStation, testStation, testStation)
	sm[testStation].AddPrice(1, testFuel, testPrice)

	n.UpdateBestStations(testFuel, testPrice, sm)

	// will call the testclient, which checks if the notification message is ok
	if !n.Notify() {
		t.Error("Test Notify failed. Notification string not as expected.")
	}
}

// Notify is the testClient notification function which tests if the message is ok
func (t *testClient) Notify(message string) bool {
	testMessage := "Good price for Diesel : 1.000€\nBest price for Diesel : 1.000€ \nat A \nin A\n\n"
	return message == testMessage
}
