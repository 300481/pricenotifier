package tankerkoenig

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/300481/pricenotifier/pkg/market"
	"github.com/alexruf/tankerkoenig-go"
)

// Tankerkoenig represents the struct for a tankerkoenig pricesource
type Tankerkoenig struct{}

// NewTankerkoenig returns an initialized *Tankerkoenig
func NewTankerkoenig() *Tankerkoenig {
	return &Tankerkoenig{}
}

// GetStations returns a list of Market Station
func (t *Tankerkoenig) GetStations() (market.Stations, error) {
	apiKey := os.Getenv("API_KEY")

	lat, err := strconv.ParseFloat(os.Getenv("LAT"), 64)
	if err != nil {
		fmt.Println("error converting LAT to float", err)
	}

	lng, err := strconv.ParseFloat(os.Getenv("LNG"), 64)
	if err != nil {
		fmt.Println("error converting LNG to float", err)
	}

	rad, err := strconv.Atoi(os.Getenv("RAD"))
	if err != nil {
		fmt.Println("error converting RAD to integer", err)
	}

	client := tankerkoenig.NewClient(apiKey, nil)
	stations, _, err := client.Station.List(lat, lng, rad)

	if err != nil {
		return nil, err
	}

	mStations := market.Stations{}

	for _, station := range stations {
		mStation := &market.Station{
			ID:        station.Id,
			Brand:     station.Brand,
			Name:      station.Name,
			Place:     station.Place,
			Lat:       station.Lat,
			Lng:       station.Lng,
			Price:     make(map[market.FuelType]float64),
			IsOpen:    station.IsOpen,
			Timestamp: time.Now(),
		}
		mStation.Price["E5"] = station.E5.(float64)
		mStation.Price["Diesel"] = station.Diesel.(float64)
		mStations = append(mStations, mStation)
	}

	return mStations, nil
}
