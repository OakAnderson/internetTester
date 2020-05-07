package nettest

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// Netdata is
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

func pathAPI() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "../API/linux-x86_64/speedtest")
}

func (test *Netdata) execTest() error {
	result, err := exec.Command(pathAPI(), "-f", "json").Output()
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

// MakeTest is
func MakeTest() (string, error) {
	var result netdata
	err := result.execTest()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s - download: %.2f Mbps - upload: %.2f Mbps",
		result.Datetime,
		result.Download.BandwidthMB,
		result.Upload.BandwidthMB,
}
