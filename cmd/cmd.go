package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

// Check panics if the error is not nil
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// GetIntervals return time durations from a string list
func GetIntervals(intervals []string) ([]time.Duration, error) {
	metrics := map[string]time.Duration{
		"h": time.Hour,
		"m": time.Minute,
		"s": time.Second,
		"":  time.Duration(1),
	}

	duration, _ := regexp.Compile("^([0-9])*([msh])?$")
	var timeIntervals []time.Duration
	for _, t := range intervals {
		args := duration.FindStringSubmatch(t)
		if args == nil {
			return nil, fmt.Errorf("usage of intervals: 10 (seconds), 10m (minutes), 10s (seconds), 10h (hours)")
		}

		num, _ := strconv.Atoi(args[1])
		timeIntervals = append(timeIntervals, metrics[args[2]]*time.Duration(num))
	}

	return timeIntervals, nil
}
