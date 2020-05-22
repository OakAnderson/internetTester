package main

import (
	"github.com/OakAnderson/internetTester/cmd"
	"github.com/OakAnderson/internetTester/nettest"
)

func main() {
	_, err := nettest.MakeTest(true)
	cmd.Check(err)
}
