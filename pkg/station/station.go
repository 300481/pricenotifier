package station

// Station represents a gas station with its price history
type Station struct {
	brand  string                       // The brand of the station
	name   string                       // The name of the station
	place  string                       // The place of the station (city name)
	prices map[string]map[int64]float64 // the price history for the fuel type
}

// NewStation returns a new Station
func NewStation(brand, name, place string) *Station {
	return &Station{
		brand:  brand,
		name:   name,
		place:  place,
		prices: make(map[string]map[int64]float64),
	}
}

// Update updates the data of a station
func (s *Station) Update(brand, name, place string) {
	s.brand = brand
	s.name = name
	s.place = place
}

// AddPrice adds a fuel price to a station for the given unix epoch timestamp
func (s *Station) AddPrice(timestamp int64, fuel string, price float64) {
	if _, ok := s.prices[fuel]; !ok {
		s.prices[fuel] = make(map[int64]float64)
	}
	s.prices[fuel][timestamp] = price
}

// PricesSince returns a slice of float64 with the prices since the given
// epoch timestamp and for the given fuel type
func (s *Station) PricesSince(since int64, fuel string) []float64 {
	prices := []float64{}
	for ts, price := range s.prices[fuel] {
		if ts >= since {
			prices = append(prices, price)
		}
	}
	return prices
}

// LatestPrice returns the latest price for the given fuel type
func (s *Station) LatestPrice(fuel string) float64 {
	var latestTimestamp int64
	var latestPrice float64

	for timestamp, price := range s.prices[fuel] {
		if timestamp > latestTimestamp {
			latestTimestamp = timestamp
			latestPrice = price
		}
	}

	return latestPrice
}
