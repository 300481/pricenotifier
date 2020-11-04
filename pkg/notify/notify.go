package notify

import (
	"fmt"
	"log"

	"github.com/300481/pricenotifier/pkg/station"
	"github.com/gregdel/pushover"
)

type fueltype string
type stationID string

// Notifier represents a notifier struct with information about best stations and last notifications
type Notifier struct {
	CurrentBestStations  map[fueltype]station.StationMap
	CurrentBestPrices    map[fueltype]map[stationID]float64
	NotifiedBestStations map[fueltype]station.StationMap
	NotifiedBestPrices   map[fueltype]map[stationID]float64
	client               Client
}

// NewNotifier returns an initialized *Notifier
func NewNotifier(client Client) *Notifier {
	return &Notifier{
		CurrentBestStations:  make(map[fueltype]station.StationMap),
		CurrentBestPrices:    make(map[fueltype]map[stationID]float64),
		NotifiedBestStations: make(map[fueltype]station.StationMap),
		NotifiedBestPrices:   make(map[fueltype]map[stationID]float64),
		client:               client,
	}
}

// UpdateBestStations updates the best stations of the Notifier
func (n *Notifier) UpdateBestStations(fuel string, bestStations station.StationMap) {
	n.CurrentBestStations[fueltype(fuel)] = bestStations

	// make prices map if needed
	if _, ok := n.CurrentBestPrices[fueltype(fuel)]; !ok {
		n.CurrentBestPrices[fueltype(fuel)] = make(map[stationID]float64)
	}

	// cleanup prices map
	for ID := range n.CurrentBestPrices[fueltype(fuel)] {
		delete(n.CurrentBestPrices[fueltype(fuel)], ID)
	}

	// save current best prices
	for ID, s := range n.CurrentBestStations[fueltype(fuel)] {
		n.CurrentBestPrices[fueltype(fuel)][stationID(ID)] = s.LatestPrice(fuel)
	}
}

// Notify will send a notification if there are best stations available currently
// Returns if a message was send
func (n *Notifier) Notify() bool {
	var msg string
	// for each fuel
	for fuel, sm := range n.CurrentBestStations {
		// for each station of the best stations
		for ID, s := range sm {
			price := n.CurrentBestPrices[fuel][stationID(ID)]
			msg += fmt.Sprintf(
				"Best price for %s : %.3f€ \nat %s \nin %s\n\n",
				string(fuel), price, s.Brand, s.Place,
			)
		}
	}
	// if there is a best station
	if len(msg) > 0 {
		return n.client.Notify(msg)
	}
	return false
}

// Client is an interface for a notification client
type Client interface {
	Notify(message string) bool
}

// PushoverClient represents a notification client for Pushover
type PushoverClient struct {
	Token     string
	User      string
	App       *pushover.Pushover
	Recipient *pushover.Recipient
}

// NewPushover returns a client interface for a pushover client
func NewPushover(token, user string) Client {
	app := pushover.New(token)
	recipient := pushover.NewRecipient(user)
	return &PushoverClient{
		Token:     token,
		User:      user,
		App:       app,
		Recipient: recipient,
	}
}

// Notify will send a notification via Pushover
// returns if the message was send without failure
func (p *PushoverClient) Notify(message string) bool {
	msg := pushover.NewMessage(message)
	response, err := p.App.SendMessage(msg, p.Recipient)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println(response)
	return true
}
