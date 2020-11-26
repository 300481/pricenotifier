package mock

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/300481/pricenotifier/pkg/market"
	"github.com/google/uuid"
)

// Customer returns a mocked Customer
func Customer() *market.Customer {
	day := 86400 * time.Second
	maxAge := 14 * day
	return &market.Customer{
		MaxAge: maxAge,
		Location: market.Geolocation{
			Lat: 49.833458,
			Lng: 8.052952,
		},
		Radius: 10,
		Fuels:  []market.FuelType{"Diesel", "E5"},
	}
}

// Station returns a mocked Station
func Station(id string, mockname string, prices map[market.FuelType]float64) *market.Station {
	return &market.Station{
		ID:    id,
		Brand: mockname,
		Name:  mockname,
		Place: mockname,
		Lat:   49.83962267405107,
		Lng:   8.12647068659747,
		Price: map[market.FuelType]float64{
			"Diesel": prices["Diesel"],
			"E5":     prices["E5"],
		},
		IsOpen:    true,
		Timestamp: time.Now(),
	}
}

// Stations returns mocked Stations
func Stations() market.Stations {
	stations := market.Stations{}
	for _, mockshortname := range strings.Split("A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T", ",") {
		id := uuid.New().String()
		mockname := fmt.Sprintf("Station %s", mockshortname)
		prices := make(map[market.FuelType]float64)
		for _, fuelType := range []string{"Diesel", "E5"} {
			prices[market.FuelType(fuelType)] = randFloat(0.9, 1.20)
		}
		station := Station(id, mockname, prices)
		stations = append(stations, station)
	}
	return stations
}

// generate a random float
func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
