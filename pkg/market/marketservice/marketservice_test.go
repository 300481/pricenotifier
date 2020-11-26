package marketservice

import (
	"math"
	"testing"

	"github.com/300481/pricenotifier/pkg/market"
	"github.com/300481/pricenotifier/pkg/market/mock"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/stat"
)

func TestAdd(t *testing.T) {
	stations := mock.Stations()
	m := NewMarketService()

	for _, station := range stations {
		m.Add(station)
		mstation := m.Stations[StationID(station.ID)]

		assert.Equalf(t, station.ID, mstation.ID, "ID for Station wrong. ID: %s", station.ID)
		assert.Equalf(t, station.Lat, mstation.Lat, "Lat for Station wrong. ID: %s", station.ID)
		assert.Equalf(t, station.Lng, mstation.Lng, "Lng for Station wrong. ID: %s", station.ID)
		assert.Equalf(t, station.Name, mstation.Name, "Name for Station wrong. ID: %s", station.ID)
		assert.Equalf(t, station.Brand, mstation.Brand, "Brand for Station wrong. ID: %s", station.ID)
		assert.Equalf(t, station.Place, mstation.Place, "Place for Station wrong. ID: %s", station.ID)

		ts := station.Timestamp
		assert.Equalf(t, station.Price["E5"], mstation.Price["E5"][ts], "Price E5 for Station wrong. ID: %s", station.ID)
		assert.Equalf(t, station.Price["Diesel"], mstation.Price["Diesel"][ts], "Price Diesel for Station wrong. ID: %s", station.ID)
		assert.Equalf(t, station.IsOpen, mstation.IsOpen[ts], "IsOpen for Station wrong. ID: %s", station.ID)
	}
}

func TestGet(t *testing.T) {
	stations := mock.Stations()
	m := NewMarketService()
	c := mock.Customer()
	prices := make(map[market.FuelType][]float64)
	goodPrice := make(map[market.FuelType]float64)

	// add station to market, fill prices slice
	for _, station := range stations {
		m.Add(station)
		for fuelType, price := range station.Price {
			prices[fuelType] = append(prices[fuelType], price)
		}
	}

	// test get all stations
	getStations := m.Get(c, "all")

	for _, fuelType := range c.Fuels {
		for _, src := range stations {
			found := false
			for _, dst := range getStations[fuelType] {
				if dst.ID == src.ID {
					found = true
				}
			}
			assert.Truef(t, found, "Station %s not got back for fuel %s and for ID %s", src.Brand, string(fuelType), src.ID)
		}
	}

	// calculate good price
	for fuelType, fuelPrices := range prices {
		mean, dev := stat.MeanStdDev(fuelPrices, nil)
		goodPrice[fuelType] = mean - dev
	}

	// find good stations for the fuels
	goodStations := make(map[market.FuelType]market.Stations)
	for _, fuelType := range c.Fuels {
		for _, station := range stations {
			if station.Price[fuelType] <= goodPrice[fuelType] {
				goodStations[fuelType] = append(goodStations[fuelType], station)
			}
		}
	}

	// test get best stations
	getStations = m.Get(c, "best")
	for _, fuelType := range c.Fuels {
		for _, src := range goodStations[fuelType] {
			found := false
			for _, dst := range getStations[fuelType] {
				if dst.ID == src.ID {
					found = true
				}
			}
			assert.Truef(t, found, "Station %s not got back for fuel %s and for ID %s price %.3f good price %.3f", src.Brand, string(fuelType), src.ID, src.Price[fuelType], goodPrice[fuelType])
		}
	}
}

func TestGoodPrice(t *testing.T) {
	stations := mock.Stations()
	m := NewMarketService()
	c := mock.Customer()
	prices := make(map[market.FuelType][]float64)

	// add station to market, fill prices slice
	for _, station := range stations {
		m.Add(station)
		for fuelType, price := range station.Price {
			prices[fuelType] = append(prices[fuelType], price)
		}
	}

	// get the good price from market
	goodPrice := m.GoodPrice(c)

	// calculate good price
	for fuelType, fuelPrices := range prices {
		mean, dev := stat.MeanStdDev(fuelPrices, nil)
		actual := math.Round(goodPrice[fuelType]*1000) / 1000
		expected := math.Round((mean-dev)*1000) / 1000
		assert.Equalf(t, expected, actual, "GoodPrice for fuel %s is wrong. Got %.3f should be %.3f", string(fuelType), actual, expected)
	}
}
