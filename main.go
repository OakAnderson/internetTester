package main

import (
    "os"
    "time"
    "strconv"
)

func main() {
    interval := 5.0
    if len(os.Args) > 1 {
        var err error
        interval, err = strconv.ParseFloat(os.Args[1], 64)
        if err != nil {
            interval = 5.0
        }
    }

    FillDataOnStructs()
    count := time.Now()
    for {
        if float64(time.Since(count).Minutes()) > interval {
            FillDataOnStructs()
            count = time.Now()
        }
    }
}