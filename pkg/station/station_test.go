package station

import (
	"testing"

	"go4.org/sort"
)

func TestUpdate(t *testing.T) {
	station := NewStation("A", "A", "A")
	station.Update("B", "B", "B")
	if station.brand != "B" ||
		station.name != "B" ||
		station.place != "B" {
		t.Error("Error on testing Station.Update()")
	}
}

func TestAddPrice(t *testing.T) {
	station := NewStation("A", "A", "A")
	var ts int64 = 1

	station.AddPrice(ts, "A", 1.0)
	station.AddPrice(ts, "B", 2.0)

	if price, ok := station.prices["A"]; !ok {
		t.Error("Price addition failed for fuel type A")
	} else {
		if price[ts] != 1.0 {
			t.Error("Price addition failed for fuel type A, price not")
		}
	}

	if price, ok := station.prices["B"]; !ok {
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

	for x := 10; x < 20; x++ {
		station.AddPrice(int64(x), "A", float64(x))
	}

	latestPrice := station.LatestPrice("A")

	if latestPrice != 19 {
		t.Error("LatestPrice() failed, did not returned the latest price.")
	}
}
