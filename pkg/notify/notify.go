package notify

import (
	"fmt"
	"log"
	"strings"

	"github.com/300481/pricenotifier/pkg/station"
	"github.com/gregdel/pushover"
)

type fueltype string
type stationID string

// Notifier represents a notifier struct with information about best stations and last notifications
type Notifier struct {
	BestStations map[fueltype]station.StationMap
	lastNotified map[fueltype]station.StationMap
	GoodPrice    map[fueltype]float64
	client       Client
}

// NewNotifier returns an initialized *Notifier
func NewNotifier(client Client) *Notifier {
	return &Notifier{
		BestStations: make(map[fueltype]station.StationMap),
		GoodPrice:    make(map[fueltype]float64),
		client:       client,
	}
}

// UpdateBestStations updates the best stations of the Notifier
func (n *Notifier) UpdateBestStations(fuel string, goodPrice float64, bestStations station.StationMap) {
	n.lastNotified = n.BestStations
	n.BestStations[fueltype(fuel)] = bestStations
	n.GoodPrice[fueltype(fuel)] = goodPrice
}

// SetClient sets the client for a notifier loaded from a persistent storage
func (n *Notifier) SetClient(client Client) {
	n.client = client
}

// Notify will send a notification if there are best stations available currently
// Returns if a message was send
func (n *Notifier) Notify() bool {
	var msg string
	// for each fuel
	for fuel, sm := range n.BestStations {
		msg += fmt.Sprintf("Good price for %s : %.3f€\n", fuel, n.GoodPrice[fuel])
		// for each station of the best stations
		for _, s := range sm {
			price := s.Latest[string(fuel)].Value
			msg += fmt.Sprintf(
				"Best price for %s : %.3f€ \nat %s \nin %s\n\n",
				string(fuel), price, s.Brand, s.Place,
			)
		}
	}
	// if there is a best station
	x := len(strings.Split(msg, "\n"))
	if x > 3 {
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
