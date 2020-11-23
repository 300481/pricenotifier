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

// Timestamp is a unix epoch timestamp
type Timestamp int64

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
	LastUpdated Timestamp
	Price       map[market.FuelType]map[Timestamp]float64
	IsOpen      map[Timestamp]bool
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
	timestamp := Timestamp(station.Timestamp.Unix())
	id := StationID(station.ID)

	if _, ok := ms.Stations[id]; !ok {
		s := &Station{
			Price:  make(map[market.FuelType]map[Timestamp]float64),
			IsOpen: make(map[Timestamp]bool),
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
	s.LastUpdated = timestamp

	// add open state
	s.IsOpen[timestamp] = station.IsOpen

	// add price for fuel type and timestamp
	for fuelType, price := range station.Price {
		if _, ok := s.Price[fuelType]; !ok {
			s.Price[fuelType] = make(map[Timestamp]float64)
		}
		s.Price[fuelType][timestamp] = price
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
	timestamp := src.LastUpdated
	dst := &market.Station{
		ID:        src.ID,
		Brand:     src.Brand,
		Name:      src.Name,
		Place:     src.Place,
		Lat:       src.Lat,
		Lng:       src.Lng,
		IsOpen:    src.IsOpen[timestamp],
		Timestamp: time.Unix(int64(timestamp), 0),
	}
	for fuelType := range src.Price {
		dst.Price[fuelType] = src.Price[fuelType][timestamp]
	}
	return dst
}

// goodPrice returns the good price for a fuel type and the stations of customers interest
// returns -1 on error
func (ms *MarketService) goodPrice(stations Stations, fuelType market.FuelType, maxAge int64) float64 {
	goodPrice := 1000.0
	prices := []float64{}

	for stationID, station := range stations {
		for timestamp, isOpen := range station.IsOpen {
			if !isOpen {
				continue
			}
			if int64(timestamp) < maxAge {
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
			maxAge := time.Now().Unix() - customer.MaxAge
			goodPrice := make(map[market.FuelType]float64)
			for _, fuelType := range customer.Fuels {
				goodPrice[fuelType] = ms.goodPrice(customerStations[fuelType], fuelType, maxAge)
				if goodPrice[fuelType] == -1 {
					// on error clean stations and continue to next fuel type
					customerStations[fuelType] = nil
					continue
				}
				for stationID, station := range customerStations[fuelType] {
					timestamp := station.LastUpdated
					if station.IsOpen[timestamp] {
						continue
					}
					if station.Price[fuelType][timestamp] <= goodPrice[fuelType] {
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
