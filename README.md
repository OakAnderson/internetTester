# Internet Tester

This is a project that make speed tests of internet using [speedtest](https://www.speedtest.net/apps/cli).

You can install the cmd programs that implements, single test or multitests that will be saved.

To install you may first clone this repository:

```bash
git clone https://github.com/OakAnderson/internetTester.git
```

And execute the make file:

```bash
cd internetTester
make install
```

The following programs will be available

```bash
nettest
```

It will make a single test and show the Download and Upload speed and latency as ping.

```bash
nettest-csv -t 3 -v 2m
```

It will make 3 in a 2 minute interval and print the results in terminal

```bash
nettest-mysql -user=<user> -password=<password> -database=<database> -t -1 10m
```

The above command will make a test every 10 minute and save the result into a configured mysql database. You can create a database and copy the [db.sql](database/db.sql) content to create the table.

It can be imported and used to make a single test or N tests. Use the command bellow to download as package:

```bash
go get -u github.com/OakAnderson/internetTester/
```

The next example makes a single test and show the print a resumed result

```go
package main

import (
    "github.com/OakAnderson/internetTester/nettest"
)

func main (
    nettest.MakeTest(true)
)
```

It can also execute infinit tests or N tests in every interval:

```go
package main

import (
    "time"

    "github.com/OakAnderson/internetTester/nettest"
)

func main (
    nettest.MultiTests(-1, true, nil, time.Minute*5, time.Minute*10)
)
```

The above example will execute tests and wait 5 minutes until the next test and wait again for 10 minutes, repeatedly.
