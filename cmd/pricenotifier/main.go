// This main is just an example and first test for using the packages

package main

import (
	"fmt"
	"time"

	"github.com/300481/pricenotifier/pkg/persistence"
	"github.com/300481/pricenotifier/pkg/pricehistory"
	"github.com/300481/pricenotifier/pkg/pricesource"
)

func main() {
	var ph *pricehistory.History

	gcs := persistence.NewGoogleCloudStorage("pricenotifier", "pricehistory.json")
	r, err := gcs.NewReader()
	if err != nil {
		ph = pricehistory.NewHistory()
	} else {
		ph = pricehistory.ReadHistory(r)
		r.Close()
	}

	stations, err := pricesource.GetStations()
	if err != nil {
		fmt.Println("error while getting gas stations", err)
		return
	}

	ts := time.Now().Unix()

	var e5, e10, diesel []float64

	for _, station := range stations {
		if station.IsOpen {
			e5 = append(e5, station.E5.(float64))
			e10 = append(e10, station.E10.(float64))
			diesel = append(diesel, station.Diesel.(float64))
		}
	}

	ph.AddFuelPrices(pricehistory.Timestamp(ts), pricehistory.Diesel, diesel)
	ph.AddFuelPrices(pricehistory.Timestamp(ts), pricehistory.E5, e5)
	ph.AddFuelPrices(pricehistory.Timestamp(ts), pricehistory.E10, e10)

	w := gcs.NewWriter()
	ph.Write(w)
	if err := w.Close(); err != nil {
		fmt.Println("error closing persistence")
	}

	// TODO add notification here
	if ph.GoodPrice(pricehistory.Diesel, 1.012) {
		fmt.Println("good")
	} else {
		fmt.Println("bad")
	}
}
