package market

import "time"

// Market is the interface for market interaction
type Market interface {
	Add(station *Station)                                                         // Add a station information to the market
	GetBestPriceStations(customer *Customer, fuelType FuelType) (Stations, error) // Get the best price stations of interest for a customer and fuel type
	Get(customer *Customer, option GetOption)
}

// FuelType is a type representing a fuel
type FuelType string

// GetOption defines the range of returned stations for the Get function
type GetOption string

const (
	// GetAll returns all stations of customer interest
	GetAll GetOption = "all"
	// GetBest returns only the stations of customer interest with the best price
	GetBest GetOption = "best"
)

// Station represents station information
type Station struct {
	ID        string               // The unique ID
	Brand     string               // The brand
	Name      string               // The name
	Place     string               // The place (city)
	Lat, Lng  float64              // The geo-coordinates
	Price     map[FuelType]float64 // The price for a type of fuel
	IsOpen    bool                 // The open state of the station
	Timestamp time.Time            // The timestamp from when the data is
}

// Stations is a collection of stations
type Stations []*Station

// Customer represents a stations customer
type Customer struct {
	MaxAge   int64       // Maximum Age of the prices for calculating a good price
	Location Geolocation // The customer geo location
}

// Geolocation represents geo coordinates
type Geolocation struct {
	Lat, Lng float64
}
