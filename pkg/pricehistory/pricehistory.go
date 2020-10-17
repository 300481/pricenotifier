package pricehistory

import (
	"encoding/json"
	"io"
	"log"

	"gonum.org/v1/gonum/stat"
)

// Timestamp represents a timestamp as seconds from the Unix epoch
type Timestamp int64

// Fuel specifies a type for fuel constants
type Fuel string

// Diesel, E10, E5 constant definitions
const (
	Diesel Fuel = "Diesel"
	E10    Fuel = "E10"
	E5     Fuel = "E5"
)

// History contains the history of mean price information
// and for a specified fuel type
type History struct {
	Items map[Fuel]map[Timestamp]*Price
}

// Price represents a mean price information for a specific time
type Price struct {
	Mean   float64 // the calculated mean price information
	StdDev float64 // the calculated standard deviation
	Count  int     // number of price information for the calculated Mean and StdDev
}

// NewHistory returns a new instantiated *History
func NewHistory() *History {
	return &History{
		Items: make(map[Fuel]map[Timestamp]*Price),
	}
}

// ReadHistory reads a History from an io.Reader and returns a *History
// returns an empty *History on error while decoding from JSON stream
func ReadHistory(r io.Reader) *History {
	var h History
	err := json.NewDecoder(r).Decode(&h)
	if err != nil {
		log.Println("error while decoding History from JSON")
		return NewHistory()
	}
	return &h
}

// Write writes a History to an io.Writer
func (h *History) Write(w io.Writer) {
	err := json.NewEncoder(w).Encode(h)
	if err != nil {
		log.Println("error while encoding History to JSON")
	}
	return
}

// AddFuelPrices adds new prices for the specified type of fuel and timestamp
func (h *History) AddFuelPrices(timestamp Timestamp, fuel Fuel, prices []float64) {
	// return when no prices given
	if len(prices) == 0 {
		return
	}

	// clean up timestamp if it would exists already
	if _, ok := h.Items[fuel][timestamp]; ok {
		delete(h.Items[fuel], timestamp)
	}

	// instantiate a new price information
	mean, std := stat.MeanStdDev(prices, nil)
	price := &Price{
		Mean:   mean,
		StdDev: std,
		Count:  len(prices),
	}
	h.Items[fuel][timestamp] = price

	return
}