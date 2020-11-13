// This main is just an example and first test for using the packages

package main

import (
	"time"
)

func main() {
	m := loadMarket()
	ts := time.Now().Unix()
	updateMarket(ts, m)
	saveMarket(m)
	send(ts, m)
}
