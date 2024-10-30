package timeutil

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

var ErrWrongTimeFormat = errors.New("wrong time format, e.g. 17:15")

type hourMinute struct {
	hours int
	minutes int
}

func newHourMinute(str string) (*hourMinute, error) {
	arr := strings.Split(str, ":")
	if len(arr) != 2 {
		return nil, ErrWrongTimeFormat
	}
	hours, err := strconv.Atoi(arr[0])
	if err != nil {
		return nil, ErrWrongTimeFormat
	}
	mins, err := strconv.Atoi(arr[1])
	if err != nil {
		return nil, ErrWrongTimeFormat
	}

	return &hourMinute{hours: hours, minutes: mins}, nil
}

func TimeSeriesOne(t string) (time.Time, error) {
	result, err := TimeSeries(t)
	if err != nil {
		return time.Now(), err
	}
	return result[0], nil
}

func TimeSeriesOneUnsafe(t string) time.Time {
	result, err := TimeSeries(t)
	if err != nil {
		log.Panicf("wrong time series format: %s", err)
	}
	return result[0]
}

func TimeSeries(series...string) ([]time.Time, error) {
	return TimeSeriesOfDay(time.Now(), series...)
}

// E.g. TimeSeriesOfDay("17:00", "17:15")
func TimeSeriesOfDay(dayOf time.Time, series...string) ([]time.Time, error) {
	hourMins := make([]*hourMinute, 0, len(series))
	for _, s := range series {
		hm, err := newHourMinute(s)
		if err != nil {
			return nil, err
		}
		hourMins = append(hourMins, hm)
	}

	result := make([]time.Time, 0, len(hourMins))
	for _, hm := range hourMins {
		result = append(result, time.Date(dayOf.Year(), dayOf.Month(), dayOf.Day(), hm.hours, hm.minutes, 0, 0, dayOf.Location()))
	}
	return result, nil
}

func Contains(series []time.Time, t time.Time) bool {
	for _, s := range series {
		if s.Unix() == t.Unix() {
			return true
		}
	}
	return false
}
