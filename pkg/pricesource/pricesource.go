package pricesource

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alexruf/tankerkoenig-go"
)

// GetStations returns a list of tankerkoenig Station
func GetStations() ([]tankerkoenig.Station, error) {
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

	return stations, nil
}
