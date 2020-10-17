package notify

import (
	"fmt"
	"log"
	"os"

	"github.com/300481/pricenotifier/pkg/pricehistory"
	"github.com/alexruf/tankerkoenig-go"
	"github.com/gregdel/pushover"
)

// Notify will handle the cloud function call
func Notify(station tankerkoenig.Station, fuel pricehistory.Fuel) {
	token := os.Getenv("PUSHOVER_TOKEN")
	user := os.Getenv("PUSHOVER_USER")
	app := pushover.New(token)
	recipient := pushover.NewRecipient(user)

	var msg string

	switch fuel {
	case "Diesel":
		msg = fmt.Sprintf("Price for Diesel: %.3f€ at %s in %s", station.Diesel.(float64), station.Brand, station.Place)
	case "E10":
		msg = fmt.Sprintf("Price for E10: %.3f€ at %s in %s", station.E10.(float64), station.Brand, station.Place)
	case "E5":
		msg = fmt.Sprintf("Price for E5: %.3f€ at %s in %s", station.E5.(float64), station.Brand, station.Place)
	}

	message := pushover.NewMessage(msg)
	response, err := app.SendMessage(message, recipient)
	if err != nil {
		log.Println(err)
	}
	log.Println(response)
}
