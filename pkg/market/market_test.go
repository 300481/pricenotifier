package market

import (
	"sort"
	"testing"
)

const (
	mints = 1000
	maxts = 1010
)

var (
	teststations   = []string{"A", "B"}
	testfuels      = []string{"E5", "Diesel"}
	testprices     []float64
	testtimestamps []int64
)

func prepareTestMarket() *Market {
	m := NewMarket()

	for x := mints; x < maxts; x++ {
		testtimestamps = append(testtimestamps, maxAge+int64(x))
		testprices = append(testprices, float64(maxts-x))
	}

	for _, s := range teststations {
		m.UpsertStation(s, s, s, s)
		for _, fuel := range testfuels {
			for index, timestamp := range testtimestamps {
				m.AddPrice(timestamp, s, fuel, testprices[index])
			}
		}
	}

	return m
}

// TestAllPrices tests allPrices()
func TestAllPrices(t *testing.T) {
	m := prepareTestMarket()

	for _, fuel := range testfuels {
		stationsCount := len(teststations)
		pricesCount := len(testprices)
		fuelpricesCount := stationsCount * pricesCount

		fuelprices := m.allPrices(fuel)

		if len(fuelprices) != fuelpricesCount {
			t.Errorf("Error getting the right number of prices for fuel: %s got: %d should be: %d", fuel, len(fuelprices), fuelpricesCount)
		}

		sort.Float64s(testprices)
		for index, price := range testprices {
			indexA := index * 2
			indexB := indexA + 1

			priceA := fuelprices[indexA]
			priceB := fuelprices[indexB]

			if (priceA != price) || (priceB != price) {
				t.Errorf("Error getting the right prices for fuel: %s price A = %.2f B = %.2f should be %.2f", fuel, priceA, priceB, price)
			}
		}
	}

}

// TestBestStations tests BestStations()
func TestBestStations(t *testing.T) {
	m := prepareTestMarket()

	for _, fuel := range testfuels {
		bestStations := m.BestStations(fuel)
		if len(bestStations) != len(teststations) {
			t.Errorf("Error getting the best stations for fuel type: %s test stations: %v best stations: %v", fuel, teststations, bestStations)
		}
	}
}
