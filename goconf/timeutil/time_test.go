package timeutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDiff(t *testing.T) {
	now := time.Now()
	timeoutAt := add(now, time.Hour)
	d := diff(now, timeoutAt)
	assert.EqualValues(t, time.Hour, d)
}


func TestToTime(t *testing.T) {
	r := ToTime(1621778975592)
	assert.EqualValues(t, 9, r.Minute())
	assert.EqualValues(t, 35, r.Second())
	assert.EqualValues(t, 592000000, r.Nanosecond())
}
