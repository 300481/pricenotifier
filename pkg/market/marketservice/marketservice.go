package marketservice

import (
	"time"

	"github.com/300481/pricenotifier/pkg/market"
	"gonum.org/v1/gonum/stat"
)

// MarketService represents the struct for all market information
type MarketService struct {
	Stations Stations
	// TODO: implement a geolocation map
}

// StationID is the unique ID of a station
type StationID string

// Stations is a map of stations
type Stations map[StationID]*Station

// Station represents the MarketService station information
type Station struct {
	ID          string
	Brand       string
	Name        string
	Place       string
	Lat, Lng    float64
	LastUpdated time.Time
	Price       map[market.FuelType]map[time.Time]float64
	IsOpen      map[time.Time]bool
}

// NewMarketService returns an initialized *MarketService
func NewMarketService() *MarketService {
	return &MarketService{
		Stations: make(map[StationID]*Station),
	}
}

// Add a station information to the market
func (ms *MarketService) Add(station *market.Station) {
	// TODO: geolocation map for customer interest station selection
	id := StationID(station.ID)

	if _, ok := ms.Stations[id]; !ok {
		s := &Station{
			Price:  make(map[market.FuelType]map[time.Time]float64),
			IsOpen: make(map[time.Time]bool),
		}
		ms.Stations[id] = s
	}

	s := ms.Stations[id]

	// update station information
	s.ID = station.ID
	s.Brand = station.Brand
	s.Name = station.Name
	s.Place = station.Place
	s.Lat = station.Lat
	s.Lng = station.Lng
	s.LastUpdated = station.Timestamp

	// add open state
	s.IsOpen[station.Timestamp] = station.IsOpen

	// add price for fuel type and timestamp
	for fuelType, price := range station.Price {
		if _, ok := s.Price[fuelType]; !ok {
			s.Price[fuelType] = make(map[time.Time]float64)
		}
		s.Price[fuelType][station.Timestamp] = price
	}
}

// getCustomerStations the stations of interest for a customer with the latest price and latest open status
func (ms *MarketService) customerStations(customer *market.Customer) Stations {
	stations := make(map[StationID]*Station)
	// mock the stations
	// TODO: implement customer based station map by geolocation and radius
	for stationid, station := range ms.Stations {
		stations[stationid] = station
	}
	return stations
}

// marketStation transforms a Station of MarketService to a market.Station
func (ms *MarketService) marketStation(src *Station) *market.Station {
	dst := &market.Station{
		ID:        src.ID,
		Brand:     src.Brand,
		Name:      src.Name,
		Place:     src.Place,
		Lat:       src.Lat,
		Lng:       src.Lng,
		IsOpen:    src.IsOpen[src.LastUpdated],
		Price:     make(map[market.FuelType]float64),
		Timestamp: src.LastUpdated,
	}
	for fuelType := range src.Price {
		dst.Price[fuelType] = src.Price[fuelType][src.LastUpdated]
	}
	return dst
}

// goodPrice returns the good price for a fuel type and the stations of customers interest
// returns -1 on error
func (ms *MarketService) goodPrice(stations Stations, fuelType market.FuelType, maxAge time.Duration) float64 {
	goodPrice := 1000.0
	prices := []float64{}

	for stationID, station := range stations {
		for timestamp, isOpen := range station.IsOpen {
			if !isOpen {
				continue
			}

			oldestTime := time.Now().Add(-maxAge)
			if timestamp.Before(oldestTime) {
				continue
			}

			price := stations[stationID].Price[fuelType][timestamp]
			prices = append(prices, price)
		}
	}

	if len(prices) == 0 {
		return -1
	}

	mean, dev := stat.MeanStdDev(prices, nil)
	goodPrice = mean - dev

	return goodPrice
}

// Get returns a map of Fuel with the stations of customer interest,
// limited by the option
func (ms *MarketService) Get(customer *market.Customer, option market.GetOption) map[market.FuelType]market.Stations {
	customerStations := make(map[market.FuelType]Stations)
	stations := make(map[market.FuelType]market.Stations)

	for _, fuelType := range customer.Fuels {
		// get the stations of customer interest
		customerStations[fuelType] = ms.customerStations(customer)

		// remove stations which have not a good price, if wanted by option
		if option == market.GetBest {
			goodPrice := make(map[market.FuelType]float64)
			for _, fuelType := range customer.Fuels {
				goodPrice[fuelType] = ms.goodPrice(customerStations[fuelType], fuelType, customer.MaxAge)
				if goodPrice[fuelType] == -1 {
					// on error clean stations and continue to next fuel type
					customerStations[fuelType] = nil
					continue
				}
				for stationID, station := range customerStations[fuelType] {
					if station.IsOpen[station.LastUpdated] {
						continue
					}
					if station.Price[fuelType][station.LastUpdated] <= goodPrice[fuelType] {
						continue
					}
					delete(customerStations[fuelType], stationID)
				}
			}
		}

		// export Stations as market.Stations
		for _, station := range customerStations[fuelType] {
			mStation := ms.marketStation(station)
			stations[fuelType] = append(stations[fuelType], mStation)
		}
	}

	return stations
}

// GoodPrice returns a map of Fuel with the good price
func (ms *MarketService) GoodPrice(customer *market.Customer) map[market.FuelType]float64 {
	customerStations := make(map[market.FuelType]Stations)
	goodPrice := make(map[market.FuelType]float64)

	for _, fuelType := range customer.Fuels {
		customerStations[fuelType] = ms.customerStations(customer)
		for _, fuelType := range customer.Fuels {
			goodPrice[fuelType] = ms.goodPrice(customerStations[fuelType], fuelType, customer.MaxAge)
		}
	}

	return goodPrice
}
