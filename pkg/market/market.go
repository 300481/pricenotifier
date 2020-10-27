package market

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/300481/pricenotifier/pkg/station"
	"gonum.org/v1/gonum/stat"
)

// Market represents the whole station market
type Market struct {
	station map[string]*station.Station // map of station IDs and stations
}

// NewMarket returns an initialized *Market
func NewMarket() *Market {
	m := &Market{
		station: make(map[string]*station.Station),
	}
	return m
}

// BestStations returns the best stations as a map with station ID as key
func (m *Market) BestStations(fuel string) map[string]*station.Station {
	bestStations := make(map[string]*station.Station)

	prices := m.allPrices(fuel)
	mean, dev := stat.MeanStdDev(prices, nil)
	goodPrice := mean - dev

	for ID, s := range m.station {
		latestPrice := s.LatestPrice(fuel)
		if latestPrice < goodPrice {
			bestStations[ID] = s
		}
	}

	return bestStations
}

// allPrices returns all prices for a fuel type
func (m *Market) allPrices(fuel string) []float64 {
	var prices []float64
	for _, s := range m.station {
		stationPrices := s.PricesSince(maxAge, fuel)
		prices = append(prices, stationPrices...)
	}
	sort.Float64s(prices)
	return prices
}

// UpsertStation upserts a station to the market
func (m *Market) UpsertStation(ID, brand, name, place string) {
	if s, ok := m.station[ID]; !ok {
		m.station[ID] = station.NewStation(brand, name, place)
	} else {
		s.Update(brand, name, place)
	}
}

// AddPrice adds a fuel price for a station and the given Unix epoch timestamp
func (m *Market) AddPrice(timestamp int64, ID, fuel string, price float64) error {
	if _, ok := m.station[ID]; !ok {
		err := fmt.Sprintf("Station with ID: %s not existing.", ID)
		return errors.New(err)
	}
	m.station[ID].AddPrice(timestamp, fuel, price)
	return nil
}

// package initialization
const daySeconds int64 = 86400

var maxAge int64

func init() {
	days, err := strconv.Atoi(os.Getenv("DAYS"))
	if err != nil {
		maxAge = 7 * daySeconds
	} else {
		maxAge = int64(days) * daySeconds
	}
}
