package main

import (
	"flag"

	"github.com/OakAnderson/internetTester/cmd"
	"github.com/OakAnderson/internetTester/internal"
	"github.com/OakAnderson/internetTester/nettest"
)

var (
	file     string
	verbose  bool
	tests    int
	interval []string
)

func init() {
	flag.StringVar(&file, "f", "nettest.csv", "a file name")
	flag.StringVar(&file, "file", "nettest.csv", "a file name")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.BoolVar(&verbose, "verbose", false, "verbose")
	flag.IntVar(&tests, "t", 1, "a number of tests to be executed. Pass -1 to ilimited tests")
	flag.IntVar(&tests, "tests", 1, "a number of tests to be executed. Pass -1 to ilimited tests")

	flag.Parse()

	interval = flag.Args()
}

func main() {
	tickers, err := cmd.GetIntervals(interval)
	cmd.Check(err)

	saver, err := internal.SaveCSV(file, "dt,latency,jitter,download,upload,packetLoss,hardware,serverId,ip,name,location,host\n")
	cmd.Check(err)

	err = nettest.MultiTests(tests, verbose, saver, tickers...)
	cmd.Check(err)
}
