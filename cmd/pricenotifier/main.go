// This main is just an example and first test for using the packages

package main

import (
	"fmt"
	"log"
	"time"

	"encoding/json"

	"github.com/300481/pricenotifier/pkg/market"
	"github.com/300481/pricenotifier/pkg/persistence"
	"github.com/300481/pricenotifier/pkg/pricesource"
)

func loadMarket() *market.Market {
	var m *market.Market = &market.Market{}

	gcs := persistence.NewGoogleCloudStorage("pricenotifier", "pricenotifier.json")
	r, err := gcs.NewReader()
	if err != nil {
		log.Printf("Error reading market from persistence: %+v", err)
		m = market.NewMarket()
	} else {
		defer r.Close()
		err := json.NewDecoder(r).Decode(m)
		if err != nil {
			log.Printf("Error reading market from persistence: %+v", err)
			m = market.NewMarket()
		}
	}
	return m
}

func saveMarket(m *market.Market) {
	gcs := persistence.NewGoogleCloudStorage("pricenotifier", "pricenotifier.json")
	w := gcs.NewWriter()
	defer w.Close()

	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Printf("Error writing market to persistence: %+v", err)
	}
}

func updateMarket(m *market.Market) {
	stations, err := pricesource.GetStations()
	if err != nil {
		fmt.Println("error while getting gas stations", err)
		return
	}

	ts := time.Now().Unix()

	for _, station := range stations {
		// insert/update station in market
		m.UpsertStation(
			station.Id,
			station.Brand,
			station.Name,
			station.Place,
		)

		// add price information if station is opened
		if station.IsOpen {
			if station.E5 != nil {
				m.AddPrice(ts, station.Id, "E5", station.E5.(float64))
			}
			if station.Diesel != nil {
				m.AddPrice(ts, station.Id, "Diesel", station.Diesel.(float64))
			}
		}
	}
}

func main() {
	m := loadMarket()
	updateMarket(m)
	saveMarket(m)
}
