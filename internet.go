package main

import (
    "fmt"
    "os"
    "time"
    "encoding/json"
    "os/exec"
    "database/sql"
    "runtime"

    _ "github.com/mattn/go-sqlite3"
)

type Internet struct {
    latency float64
    jitter float64
    download float64
    upload float64
    packageLoss int
    datetime string
    hardware string
}

type Server struct {
    id int
    port int
    ip string
    name string
    location string
    host string
}

var server Server
var internet Internet


func treatError(err error) {
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func connectDataBase() (db *sql.DB) {
    db, err := sql.Open("sqlite3", "data/internet.db")

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    return db
}

func RecoverData() {
    db := connectDataBase()

    rows, err := db.Query("SELECT * FROM results ORDER BY id DESC")
    treatError(err)

    for rows.Next() {
        var id int
        rows.Scan(&id, &internet.datetime, &internet.latency,
            &internet.jitter, &internet.download, &internet.upload,
            &internet.packageLoss, &internet.hardware, &server.id,
            &server.port, &server.ip, &server.name, &server.location,
            &server.host)

        fmt.Printf("%d | %s - download: %.2f Mbps - upload: %.2f Mbps\n", id, internet.datetime, internet.download, internet.upload)
    }
}

func insertIntoDataBase() {
    db := connectDataBase()
    defer db.Close()

    statement, err := db.Prepare("INSERT INTO results (datetime, latency, jitter, download, upload, packageLoss, hardware, serverId, port, ip, name, location, host) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

    treatError(err)

    statement.Exec(internet.datetime, internet.latency, internet.jitter,
        internet.download, internet.upload, internet.packageLoss,
        internet.hardware, server.id, server.port, server.ip,
        server.name, server.location, server.host)
}

func fillInternetData(resultMap map[string]interface {}) {
    internet.latency = resultMap["ping"].(map[string]interface{})["latency"].(float64)
    internet.jitter = resultMap["ping"].(map[string]interface{})["jitter"].(float64)
    internet.download = resultMap["download"].(map[string]interface{})["bytes"].(float64)/1000000.0
    internet.upload = resultMap["upload"].(map[string]interface{})["bytes"].(float64)/1000000.0
    internet.packageLoss = int(resultMap["packetLoss"].(float64))
    internet.hardware = resultMap["interface"].(map[string]interface {})["name"].(string)
    internet.datetime = time.Now().Format("2006-01-02 15:04:05")
}

func fillServerData(resultMap map[string]interface{}) {
    server.id = int(resultMap["server"].(map[string]interface{})["id"].(float64))
    server.port = int(resultMap["server"].(map[string]interface{})["port"].(float64))
    server.ip = resultMap["server"].(map[string]interface{})["ip"].(string)
    server.name = resultMap["server"].(map[string]interface{})["name"].(string)
    server.location = resultMap["server"].(map[string]interface{})["location"].(string)
    server.host = resultMap["server"].(map[string]interface{})["host"].(string)
}

func FillDataOnStructs() {
    resultMap := GenerateResultMap()

    fillInternetData(resultMap)
    fillServerData(resultMap)

    fmt.Printf("%s - download: %.2f Mbps - upload: %.2f Mbps\n", internet.datetime, internet.download, internet.upload)

    insertIntoDataBase()
}

func GenerateResultMap() (resultMap map[string] interface {}) {
    result, err := ExecuteTest()

    treatError(err)

    json.Unmarshal(result, &resultMap)

    return resultMap
}

func ExecuteTest() (result []byte, err error) {
    if runtime.GOOS == "linux" {
        result, err = exec.Command("./speedtestAPI/linux-x86_64/speedtest", "-f", "json").Output()
    } else if runtime.GOOS == "windows" {
        result, err = exec.Command("speedtestAPI\\windows-x86_64\\speedtest.exe", "-f", "json").Output()
    }

    return
}