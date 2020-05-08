package nettest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // Mysql driver
	"github.com/schollz/progressbar/v3"
)

const (
	testSpinner  int = 69
	nextSpinner  int = 70
	testThrottle int = 50
	nextThrottle int = 400
)

// Netdata a struct that keeps speedtest relevant results
type Netdata struct {
	Datetime   string
	PacketLoss float64 `json:"packetLoss"`

	Ping struct {
		Latency float64 `json:"latency"`
		Jitter  float64 `json:"jitter"`
	} `json:"ping"`

	Download struct {
		Bandwidth   int `json:"bandwidth"`
		BandwidthMB float32
	} `json:"download"`

	Upload struct {
		Bandwidth   int `json:"bandwidth"`
		BandwidthMB float32
	} `json:"Upload"`

	Interface struct {
		Hardware string `json:"hardware"`
	} `json:"interface"`

	Server struct {
		ID       int    `json:"id"`
		Port     int    `json:"port"`
		IP       string `json:"ip"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Host     string `json:"host"`
	} `json:"server"`
}

var speedtest string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	speedtest = filepath.Join(filepath.Dir(filename), "../API/linux-x86_64/speedtest")
}

func (test *Netdata) execTest(verbose bool) (err error) {
	var result []byte
	if verbose {
		result, err = execTestVerbose()
	} else {
		result, err = execTest()
	}
	if err != nil {
		return err
	}

	return test.loadFields(result)
}

func (test *Netdata) loadFields(results []byte) error {
	test.Datetime = time.Now().Format("2006-01-02 15:04:05")
	err := json.Unmarshal(results, test)
	if err != nil {
		return err
	}
	test.Download.BandwidthMB = float32(test.Download.Bandwidth) / float32(1e5)
	test.Upload.BandwidthMB = float32(test.Upload.Bandwidth) / float32(1e5)
	return nil
}

func (test Netdata) String() string {
	return fmt.Sprintf(
		"%s - download: %.2f Mbps - upload: %.2f Mbps - ping: %.0f ms",
		test.Datetime,
		test.Download.BandwidthMB,
		test.Upload.BandwidthMB,
		test.Ping.Latency,
	)
}

// Save insert the results into the configurated database
func (test Netdata) Save() error {
	db, err := connDatabase()
	if err != nil {
		return err
	}

	insert, err := db.Prepare(
		"INSERT speedtest SET dt=?,latency=?,jitter=?,download=?,upload=?,packetLoss=?,hardware=?,serverId=?,port=?,ip=?,name=?,location=?,host=?",
	)
	if err != nil {
		return err
	}

	_, err = insert.Exec(
		test.Datetime,
		test.Ping.Latency,
		test.Ping.Jitter,
		test.Download.BandwidthMB,
		test.Upload.BandwidthMB,
		test.PacketLoss,
		test.Interface.Hardware,
		test.Server.ID,
		test.Server.Port,
		test.Server.Name,
		test.Server.IP,
		test.Server.Location,
		test.Server.Host,
	)

	return err
}

// connDatabase connect to mysql database and return it
func connDatabase() (db *sql.DB, err error) {
	user, dbname := "oak", "internetTester"
	return sql.Open("mysql", user+":@/"+dbname)
}

func showProgressbar(c chan bool, description string, ms, spinner int) {
	bar := progressbar.NewOptions(
		-1,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSpinnerType(spinner),
		progressbar.OptionThrottle(time.Millisecond*time.Duration(ms)),
	)

	for {
		select {
		case <-c:
			fmt.Printf("\r                                          \r")
			return
		default:
			bar.Add(1)
			break
		}
	}
}

func execTestVerbose() (result []byte, err error) {
	c := make(chan bool)

	go func() {
		result, err = exec.Command(speedtest, "-f", "json").Output()
		c <- true
	}()
	showProgressbar(c, "Executing test", testThrottle, testSpinner)
	return
}
// MakeTest execute a single speedtest and return the results with a formated
// string and its struct
func MakeTest(verbose bool) (*Netdata, error) {
	var result Netdata
	err := result.execTest(verbose)
	if err != nil {
		return nil, err
	}

	if verbose {
		fmt.Println(result)
	}

	return &result, nil
}

// MultiTests execute n tests with an interval between them
func MultiTests(times int, verbose, save bool, interval ...time.Duration) error {
	var waitInterval bool
	if len(interval) > 0 {
		for _, v := range interval {
			if v < time.Minute {
				return fmt.Errorf("an interval must be bigger than 1 minute")
			}
		}
		waitInterval = true
	}

	var count int
	nextTest := func() bool {
		if times < 0 {
			return true
		}
		count++
		return count-1 < times
	}

	for nextTest() {
		nd, err := MakeTest(verbose)
		if err != nil {
			return err
		}

		if save {
			if err = nd.Save(); err != nil {
				return err
			}
		}

		if waitInterval {
			if count-1 > len(interval) {
				break
			}

			ticker := interval[count%len(interval)]

			var c chan bool
			if verbose {
				c = make(chan bool)
				go func() {
					showProgressbar(
						c,
						fmt.Sprintf(
							"Next test in: %s",
							time.Now().Add(ticker).Format("2006-01-02 15:04:05")),
						nextThrottle,
						nextSpinner,
					)
				}()
			}
			time.Sleep(ticker)

			c <- true
		}
	}

	if verbose {
		fmt.Println("Done!")
	}

	return nil
}
