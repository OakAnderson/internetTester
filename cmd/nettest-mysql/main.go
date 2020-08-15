package main

import (
	"flag"
	"fmt"

	"github.com/OakAnderson/internetTester/cmd"
	"github.com/OakAnderson/internetTester/internal"
	"github.com/OakAnderson/internetTester/nettest"

	_ "github.com/go-sql-driver/mysql"
)

var (
	user, pswd, db string
	verbose        bool
	tests          int
	interval       []string
)

func init() {
	flag.StringVar(&user, "u", "", "the database user")
	flag.StringVar(&user, "user", "", "the database user")
	flag.StringVar(&pswd, "p", "", "the password of database user")
	flag.StringVar(&pswd, "password", "", "the password of database user")
	flag.StringVar(&db, "d", "", "the database name")
	flag.StringVar(&db, "database", "", "the database name")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.BoolVar(&verbose, "verbose", false, "verbose")
	flag.IntVar(&tests, "t", 1, "a number of tests to be executed. Pass -1 to ilimited tests")
	flag.IntVar(&tests, "tests", 1, "a number of tests to be executed. Pass -1 to ilimited tests")

	flag.Parse()

	interval = flag.Args()
}

func main() {
	if user == "" || db == "" {
		cmd.Check(fmt.Errorf("the database name and user must be passed. Call 'nettest-mysql -h' for usage"))
	}

	tickers, err := cmd.GetIntervals(interval)
	cmd.Check(err)

	saver, err := internal.SaveMysql(user, pswd, db)
	cmd.Check(err)

	err = nettest.MultiTests(tests, verbose, saver, tickers...)
	cmd.Check(err)
}
