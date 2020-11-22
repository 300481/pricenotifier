package marketservice

import (
	"errors"
	"time"

	"github.com/300481/pricenotifier/pkg/market"
	"gonum.org/v1/gonum/stat"
)

// MarketService represents the struct for all market information
type MarketService struct {
	Stations Stations
}

// Timestamp is a unix epoch timestamp
type Timestamp int64

// StationID is the unique ID of a station
type StationID string

// Stations is a map of stations
type Stations map[StationID]*Station

// Station represents the MarketService station information
type Station struct {
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
	// TODO: implement customer based station map
	for stationid, station := range ms.Stations {
		stations[stationid] = station
	}
	return stations
}

// marketStation transforms a Station of MarketService to a market.Station
func (ms *MarketService) marketStation(stationID StationID, src *Station) *market.Station {
	timestamp := src.LastUpdated
	dst := &market.Station{
		ID:        string(stationID),
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
func (ms *MarketService) goodPrice(customer *market.Customer, fuelType market.FuelType) (float64, Stations, error) {
	goodPrice := 1000.0
	prices := []float64{}
	maxAge := time.Now().Unix() - customer.MaxAge
	stations := ms.customerStations(customer)

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
		return 0.0, nil, errors.New("No prices available")
	}

	mean, dev := stat.MeanStdDev(prices, nil)
	goodPrice = mean - dev

	return goodPrice, stations, nil
}

// GetBestPriceStations the best price stations of interest for a customer and fuel type
func (ms *MarketService) GetBestPriceStations(customer *market.Customer, fuelType market.FuelType) (market.Stations, error) {
	// first get what is a good price
	goodPrice, customerStations, err := ms.goodPrice(customer, fuelType)
	if err != nil {
		return nil, err
	}

	var stations market.Stations

	for stationID, station := range customerStations {
		mStation := ms.marketStation(stationID, station)
		if !mStation.IsOpen {
			continue
		}
		if mStation.Price[fuelType] > goodPrice {
			continue
		}
		stations = append(stations, mStation)
	}

	return stations, nil
}
