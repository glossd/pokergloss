package timeutil

import "time"

func Now() int64 {
	return NowAdd(0)
}

func Duration(d time.Duration) int64 {
	return int64(d) / 1e6
}

func Time(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func ToTime(t int64) time.Time {
	return time.Unix(t/1000, t%1000*1e6)
}

func Add(time int64, add time.Duration) int64 {
	return time + Duration(add)
}

// Returns point in unix time of now plus duration in milliseconds
func NowAdd(d time.Duration) int64 {
	return add(time.Now(), d)
}

func add(t time.Time, d time.Duration) int64 {
	return t.Add(d).UnixNano() / 1e6
}

func MinusNow(timeoutAt int64) time.Duration {
	return diff(time.Now(), timeoutAt)
}

func NowMinus(millis int64) time.Duration {
	milliDiff := time.Now().UnixNano()/1e6 - millis
	return time.Duration(milliDiff * 1e6)
}

func diff(t time.Time, timeoutAt int64) time.Duration {
	milliDiff := timeoutAt - t.UnixNano()/1e6
	return time.Duration(milliDiff * 1e6)
}

func Midnight(t time.Time) int64 {
	return Time(t.Truncate(24*time.Hour))
}

func Multiply(d time.Duration, multiplier float64) time.Duration {
	return time.Duration(int64(float64(d)*multiplier))
}

func AfterTimeAt(at int64) <-chan time.Time {
	return time.After(MinusNow(at))
}
