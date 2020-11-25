package marketservice

import (
	"testing"

	"github.com/300481/pricenotifier/pkg/market/mock"
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
		ts := Timestamp(station.Timestamp.Unix())
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

	for _, station := range stations {
		m.Add(station)
	}

	gstations := m.Get(c, "all")

	for _, fuelType := range c.Fuels {
		for _, src := range stations {
			found := false
			for _, dst := range gstations[fuelType] {
				if dst.ID == src.ID {
					found = true
				}
			}
			if !found {
				t.Errorf("Station not got back for fuel %s and for ID %s", string(fuelType), src.ID)
			}
		}
	}

	// TODO implement test for "best"
}
