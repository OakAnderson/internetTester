package main

import (
	"time"

	"github.com/OakAnderson/internetTester/nettest"
)

func main() {
	err := nettest.MultiTests(-1, true, true, time.Minute*10)
	if err != nil {
		panic(err)
	}
}
