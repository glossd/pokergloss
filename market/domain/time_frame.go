package domain

import "time"

type TimeFrame string
const (
	Day TimeFrame = "day"
	Week TimeFrame = "week"
	Month TimeFrame = "month"
)

func (tf TimeFrame) Duration() time.Duration {
	switch tf {
	case Day:
		return 24*time.Hour
	case Week:
		return 7*24*time.Hour
	case Month:
		return 30*24*time.Hour
	default:
		return 0
	}
}

func errUnknownTimeFrame(tf TimeFrame) error {
	return E("unknown timeframe %s", tf)
}
