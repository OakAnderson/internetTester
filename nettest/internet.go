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
func MakeTest() (string, *Netdata, error) {
	var result Netdata
	err := result.execTest()
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf(
		"%s - download: %.2f Mbps - upload: %.2f Mbps",
		result.Datetime,
		result.Download.BandwidthMB,
		result.Upload.BandwidthMB,
	), &result, nil
}

// Save is
func (test Netdata) Save() error {
	db, err := database.ConnDatabase()
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
