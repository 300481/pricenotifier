package pricesource

import (
	"github.com/300481/pricenotifier/pkg/market"
)

// PriceSource defines the interface of a PriceSource
type PriceSource interface {
	GetStations() (market.Stations, error)
}
