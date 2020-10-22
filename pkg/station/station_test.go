package station

import (
	"testing"
)

func TestUpdate(t *testing.T) {
	station := NewStation("ARAL", "ARAL", "Berlin")
	station.Update("ESSO", "ESSO", "Stuttgart")
	if station.Brand != "ESSO" ||
		station.Name != "ESSO" ||
		station.Place != "Stuttgart" {
		t.Error("Error on testing Station.Update()")
	}
}

func TestAddFuelPrices(t *testing.T) {
	station := NewStation("ARAL", "ARAL", "Berlin")
	var ts int64 = 1

	station.AddFuelPrice(ts, "A", 1.0)
	station.AddFuelPrice(ts, "B", 2.0)

	if price, ok := station.Prices["A"]; !ok {
		t.Error("Price addition failed for fuel type A")
	} else {
		if price[ts] != 1.0 {
			t.Error("Price addition failed for fuel type A, price not")
		}
	}

	if price, ok := station.Prices["B"]; !ok {
		t.Error("Price addition failed for fuel type A")
	} else {
		if price[ts] != 2.0 {
			t.Error("Price addition failed for fuel type A, price not")
		}
	}
}

func TestPricesSince(t *testing.T) {
	station := NewStation("ARAL", "ARAL", "Berlin")

	for x := 10; x < 20; x++ {
		station.AddFuelPrice(int64(x), "A", float64(x))
	}

	for x := 10; x < 20; x++ {
		pricesSince := station.PricesSince(int64(x), "A")
		if len(pricesSince) != (20 - x) {
			t.Error("PricesSince() failed, did not returned the right count of prices.")
		}
	}
}
