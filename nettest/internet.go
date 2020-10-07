package nettest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	testSpinner  int = 59
	nextSpinner  int = 67
	testThrottle int = 50
	nextThrottle int = 400
)

// Saver is an interface that is used to save data
type Saver interface {
	// Exec must save the args into something, sql.Stmt implements
	// this interface, just adjust your database into this
	// sql.Result is not used, is just here to implement sql.Stmt.Exec
	Exec(args ...interface{}) (sql.Result, error)
}

// Netdata a struct that keeps speedtest relevant results
type Netdata struct {
	Datetime   string
	Start      time.Time
	PacketLoss float64 `json:"packetLoss"`
	Type       string  `json:"type"`
	Message    string  `json:"message"`

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
		Hardware string `json:"name"`
	} `json:"interface"`

	Server struct {
		ID       int    `json:"id"`
		IP       string `json:"ip"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Host     string `json:"host"`
	} `json:"server"`
}

var speedtest string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	bin := map[string]string{"windows": "windows/speedtest.exe", "linux": "linux/speedtest"}
	speedtest = filepath.Join(filepath.Dir(filename), "../API/"+bin[runtime.GOOS])
}

func (test *Netdata) execTest(verbose bool) (err error) {
	var result []byte
	test.Start = time.Now()
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

	strResults := strings.TrimLeftFunc(
		string(results),
		func(r rune) bool {
			return r != '{'
		},
	)
	strResults = strings.TrimRightFunc(
		strResults,
		func(r rune) bool {
			return r != '}'
		},
	)

	err := json.Unmarshal([]byte(strResults), test)
	if err != nil {
		return err
	}

	if test.Type != "result" {
		return fmt.Errorf(test.Message)
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
func (test Netdata) Save(s Saver) error {
	_, err := s.Exec(
		test.Datetime,
		test.Ping.Latency,
		test.Ping.Jitter,
		test.Download.BandwidthMB,
		test.Upload.BandwidthMB,
		test.PacketLoss,
		test.Interface.Hardware,
		test.Server.ID,
		test.Server.IP,
		test.Server.Name,
		test.Server.Location,
		test.Server.Host,
	)

	return err
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
			close(c)
			return
		case <-time.After(time.Millisecond * 50):
			bar.Add(1)
			break
		}
	}
}

func execTestVerbose() (result []byte, err error) {
	c := make(chan bool)

	go func() {
		result, err = execTest()
		c <- true
	}()

	showProgressbar(c, "Executing test", testThrottle, testSpinner)

	return
}

func execTest() (result []byte, err error) {
	cmd := exec.Command(speedtest, "-f", "json")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	resultErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}
	result, err = ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()

	if err != nil {
		if len(resultErr) == 0 {
			return nil, err
		}
		return resultErr, nil
	}

	return
}

// MakeTest execute a single speedtest and return the results with a formated
// string and its struct
func MakeTest(verbose bool) (*Netdata, error) {
	var result Netdata

	if err := result.execTest(verbose); err != nil {
		return nil, err
	}

	if verbose {
		fmt.Println(result)
	}

	return &result, nil
}

// MultiTests execute n tests with an interval between them
func MultiTests(times int, verbose bool, save Saver, interval ...time.Duration) error {
	var waitInterval bool
	if len(interval) > 0 {
		for _, v := range interval {
			if v < time.Minute {
				return fmt.Errorf("an interval must be bigger than 1 minute, not %v", v)
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
		return count <= times
	}

	for nextTest() {
		nd, err := MakeTest(verbose)
		if err != nil {
			return err
		}

		if save != nil {
			if err = nd.Save(save); err != nil {
				return err
			}
		}

		if waitInterval && nextTest() {
			count--
			ticker := interval[(count-1)%len(interval)]
			ticker -= time.Since(nd.Start)

			if verbose {
				timerVerbose(ticker)
			} else {
				time.Sleep(ticker)
			}
		}
	}

	if verbose {
		fmt.Println("Done!")
	}

	return nil
}

func timerVerbose(t time.Duration) {
	c := make(chan bool)

	go func() {
		time.Sleep(t)
		c <- true
	}()

	showProgressbar(
		c,
		fmt.Sprintf("Next test in: %s",
			time.Now().Add(t).Format("2006-01-02 15:04:05"),
		),
		nextThrottle,
		nextSpinner,
	)
}
