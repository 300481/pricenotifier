package pricehistory

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"time"

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
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	err := enc.Encode(h)
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

	// instantiate a new price information
	mean, std := stat.MeanStdDev(prices, nil)
	price := &Price{
		Mean:   mean,
		StdDev: std,
		Count:  len(prices),
	}

	// instatiate the fuelmap if not existing
	fuelmap, ok := h.Items[fuel]
	if !ok {
		fuelmap = make(map[Timestamp]*Price)
		h.Items[fuel] = fuelmap
	}
	fuelmap[timestamp] = price

	return
}

// GoodPrice returns true when given price is below Mean-StdDev
func (h *History) GoodPrice(fuel Fuel, price float64) bool {
	// TODO: check only for a given range in the past
	var mean, std []float64

	lastInHistory := time.Now().Unix() - maxAge

	for timestamp, price := range h.Items[fuel] {
		// if price record older then the defined maximum age
		if int64(timestamp) < lastInHistory {
			continue
		}

		for x := 0; x < price.Count; x++ {
			mean = append(mean, price.Mean)
			std = append(std, price.StdDev)
		}
	}

	meanAll := stat.Mean(mean, nil)
	stdAll := stat.Mean(std, nil)

	good := meanAll - stdAll

	return price < good
}

// CleanHistory removes prices older than the max period given in days
func (h *History) CleanHistory(fuel Fuel) {
	lastInHistory := time.Now().Unix() - maxAge
	for timestamp := range h.Items[fuel] {
		// if price record older then the defined max period
		if int64(timestamp) < lastInHistory {
			delete(h.Items[fuel], timestamp)
		}
	}
}

// initialization
const daySeconds int64 = 86400

var maxAge int64

func init() {
	days, err := strconv.Atoi(os.Getenv("DAYS"))
	if err != nil {
		maxAge = 7 * daySeconds
	} else {
		maxAge = int64(days) * daySeconds
	}
}
