package station

import (
	"testing"

	"go4.org/sort"
)

func TestUpdate(t *testing.T) {
	station := NewStation("A", "A", "A")
	station.Update("B", "B", "B")
	if station.Brand != "B" ||
		station.Name != "B" ||
		station.Place != "B" {
		t.Error("Error on testing Station.Update()")
	}
}

func TestAddPrice(t *testing.T) {
	station := NewStation("A", "A", "A")
	var ts int64 = 1

	station.AddPrice(ts, "A", 1.0)
	station.AddPrice(ts, "B", 2.0)

	if price, ok := station.Prices["A"]; !ok {
		t.Error("Price addition failed for fuel type A")
	} else {
		if price[ts] != 1.0 {
			t.Error("Price addition failed for fuel type A, price not")
		}
	}

	if price, ok := station.Prices["B"]; !ok {
		t.Error("Price addition failed for fuel type B")
	} else {
		if price[ts] != 2.0 {
			t.Error("Price addition failed for fuel type B, price not")
		}
	}
}

func TestPricesSince(t *testing.T) {
	station := NewStation("A", "A", "A")

	for x := 10; x < 20; x++ {
		station.AddPrice(int64(x), "A", float64(x))
	}

	for x := 10; x < 20; x++ {
		pricesSince := station.PricesSince(int64(x), "A")
		if len(pricesSince) != (20 - x) {
			t.Error("PricesSince() failed, did not returned the right count of prices.")
		}

		sort.Float64s(pricesSince)
		if pricesSince[0] != float64(x) {
			t.Error("PricesSince() failed, did not returned the right prices.")
		}
	}
}

func TestLatestPrice(t *testing.T) {
	station := NewStation("A", "A", "A")

	_, _, err := station.LatestPrice("A")
	if err == nil {
		t.Error("LatestPrice() failed, should return an error.")
	}

	for x := 10; x < 20; x++ {
		station.AddPrice(int64(x), "A", float64(x))
	}

	ts, latestPrice, err := station.LatestPrice("A")
	if err != nil {
		t.Error("LatestPrice() failed, should not return an error.")
	}

	if ts != int64(19) {
		t.Error("LatestPrice() failed, did not returned the right timestamp.")
	}

	if latestPrice != 19 {
		t.Error("LatestPrice() failed, did not returned the latest price.")
	}
}
