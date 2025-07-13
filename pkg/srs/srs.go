package srs

import "time"

var srsIntervals = []int{1, 3, 7, 14, 30, 90}

func CalculateNextRevisionDate(srsLevel int) time.Time {
	index := srsLevel - 1
	var intervalDays int

	if index >= 0 && index < len(srsIntervals) {
		intervalDays = srsIntervals[index]
	} else if index >= len(srsIntervals) {
		intervalDays = srsIntervals[len(srsIntervals)-1]
	} else {
		intervalDays = srsIntervals[0]
	}

	nextDate := time.Now().AddDate(0, 0, intervalDays)
	return nextDate
}
