package station

// StationMap is a map of stations with the station ID as the key
type StationMap map[string]*Station

// NewStationMap returns an initialized StationMap
func NewStationMap() StationMap {
	return make(StationMap)
}

// Station represents a gas station with its price history
type Station struct {
	Brand  string                       // The brand of the station
	Name   string                       // The name of the station
	Place  string                       // The place of the station (city name)
	Prices map[string]map[int64]float64 // the price history for the fuel type
}

// NewStation returns a new Station
func NewStation(brand, name, place string) *Station {
	return &Station{
		Brand:  brand,
		Name:   name,
		Place:  place,
		Prices: make(map[string]map[int64]float64),
	}
}

// Update updates the data of a station
func (s *Station) Update(brand, name, place string) {
	s.Brand = brand
	s.Name = name
	s.Place = place
}

// AddPrice adds a fuel price to a station for the given unix epoch timestamp
func (s *Station) AddPrice(timestamp int64, fuel string, price float64) {
	if _, ok := s.Prices[fuel]; !ok {
		s.Prices[fuel] = make(map[int64]float64)
	}
	s.Prices[fuel][timestamp] = price
}

// PricesSince returns a slice of float64 with the prices since the given
// epoch timestamp and for the given fuel type
func (s *Station) PricesSince(since int64, fuel string) []float64 {
	prices := []float64{}
	for ts, price := range s.Prices[fuel] {
		if ts >= since {
			prices = append(prices, price)
		}
	}
	return prices
}

// LatestPrice returns the latest price for the given fuel type
func (s *Station) LatestPrice(fuel string) (timestamp int64, price float64) {
	var latestTimestamp int64
	var latestPrice float64

	for ts, p := range s.Prices[fuel] {
		if ts > latestTimestamp {
			latestTimestamp = ts
			latestPrice = p
		}
	}

	return latestTimestamp, latestPrice
}
