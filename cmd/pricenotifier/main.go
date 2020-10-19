// This main is just an example and first test for using the packages

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/300481/pricenotifier/pkg/notify"
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
			if station.E5 != nil {
				e5 = append(e5, station.E5.(float64))
			}
			if station.E10 != nil {
				e10 = append(e10, station.E10.(float64))
			}
			if station.Diesel != nil {
				diesel = append(diesel, station.Diesel.(float64))
			}
		}
	}

	ph.AddFuelPrices(pricehistory.Timestamp(ts), pricehistory.E5, e5)
	ph.AddFuelPrices(pricehistory.Timestamp(ts), pricehistory.E10, e10)
	ph.AddFuelPrices(pricehistory.Timestamp(ts), pricehistory.Diesel, diesel)

	ph.CleanHistory(pricehistory.E5)
	ph.CleanHistory(pricehistory.E10)
	ph.CleanHistory(pricehistory.Diesel)

	w := gcs.NewWriter()
	ph.Write(w)
	if err := w.Close(); err != nil {
		fmt.Println("error closing persistence")
	}

	for _, station := range stations {
		if station.IsOpen {
			if station.E5 != nil {
				if ph.GoodPrice(pricehistory.E5, station.E5.(float64)) {
					notify.Notify(station, pricehistory.E5)
					log.Printf("found good price for E5: %.3f€ at %s %s", station.E5.(float64), station.Brand, station.Place)
				}
			}
			// if station.E10 != nil {
			// 	if ph.GoodPrice(pricehistory.E10, station.E10.(float64)) {
			// 		notify.Notify(station, pricehistory.E10)
			// 	}
			// }
			if station.Diesel != nil {
				if ph.GoodPrice(pricehistory.Diesel, station.Diesel.(float64)) {
					notify.Notify(station, pricehistory.Diesel)
					log.Printf("found good price for Diesel: %.3f€ at %s %s", station.Diesel.(float64), station.Brand, station.Place)
				}
			}
		}
	}
}
