package timeutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeSeries(t *testing.T) {
	now := time.Now()
	dates, err := TimeSeriesOfDay(now, "15:15", "15:30")
	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(dates))
	assert.EqualValues(t, now.Day(), dates[0].Day())
	assert.EqualValues(t, 15, dates[0].Hour())
	assert.EqualValues(t, 15, dates[0].Minute())

	assert.EqualValues(t, 15, dates[1].Hour())
	assert.EqualValues(t, 30, dates[1].Minute())

}
