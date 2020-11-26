package marketservice

import (
	"testing"

	"github.com/300481/pricenotifier/pkg/market"
	"github.com/300481/pricenotifier/pkg/market/mock"
	"gonum.org/v1/gonum/stat"
)

func TestAdd(t *testing.T) {
	stations := mock.Stations()
	m := NewMarketService()

	for _, station := range stations {
		m.Add(station)
		mstation := m.Stations[StationID(station.ID)]
		if station.Brand != mstation.Brand {
			t.Errorf("Brand for Station wrong. ID: %s", station.ID)
		}
		if station.ID != mstation.ID {
			t.Errorf("ID for Station wrong. ID: %s", station.ID)
		}
		if station.Name != mstation.Name {
			t.Errorf("Name for Station wrong. ID: %s", station.ID)
		}
		if station.Place != mstation.Place {
			t.Errorf("Place for Station wrong. ID: %s", station.ID)
		}
		if station.Lat != mstation.Lat {
			t.Errorf("Lat for Station wrong. ID: %s", station.ID)
		}
		if station.Lng != mstation.Lng {
			t.Errorf("Lng for Station wrong. ID: %s", station.ID)
		}
		ts := station.Timestamp
		if station.Price["Diesel"] != mstation.Price["Diesel"][ts] {
			t.Errorf("Price Diesel for Station wrong. ID: %s", station.ID)
		}
		if station.Price["E5"] != mstation.Price["E5"][ts] {
			t.Errorf("Price E5 for Station wrong. ID: %s", station.ID)
		}
		if station.IsOpen != mstation.IsOpen[ts] {
			t.Errorf("IsOpen for Station wrong. ID: %s", station.ID)
		}
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
			if !found {
				t.Errorf("Station %s not got back for fuel %s and for ID %s", src.Brand, string(fuelType), src.ID)
			}
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
			if !found {
				t.Errorf("Station %s not got back for fuel %s and for ID %s price %.3f good price %.3f", src.Brand, string(fuelType), src.ID, src.Price[fuelType], goodPrice[fuelType])
			}
		}
	}
}
