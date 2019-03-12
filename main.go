package main

import (
	"time"

	"./beacon"
)

func main() {
	b := beacon.NewBeacon("localhost")
	period := time.Duration(10) * time.Second
	b.Loop(period)
	time.Sleep(time.Duration(30) * time.Second)
	b.Stop()
}
