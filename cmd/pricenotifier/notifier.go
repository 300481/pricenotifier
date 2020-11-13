package main

import (
	"log"
	"os"

	"github.com/300481/pricenotifier/pkg/market"
	"github.com/300481/pricenotifier/pkg/notify"
)

func send(timestamp int64, m *market.Market) {
	log.Println("send notification.")
	client := notify.NewPushover(
		os.Getenv("PUSHOVER_TOKEN"),
		os.Getenv("PUSHOVER_USER"),
	)
	notifier := notify.NewNotifier(client)
	for _, fuel := range []string{"Diesel", "E5"} {
		goodPrice, bestStations := m.BestStations(timestamp, fuel)
		log.Println("Good Price for ", fuel, goodPrice)
		notifier.UpdateBestStations(fuel, goodPrice, bestStations)
		for ID, s := range bestStations {
			log.Println("Good station for", fuel, ID, s.Brand, s.Name, s.Place)
		}
	}
	notifier.Notify()
}
