package paging

import (
	"fmt"
)

type Params struct {
	Skip      int64
	Limit     int64
	SkipEmpty bool
	SkipFull  bool
}

var DefaultParams = Params{
	Skip:      0,
	Limit:     20,
	SkipEmpty: false,
	SkipFull:  false,
}

const MaxLimit = 30

func NewParams(skip, limit int64, skipEmpty, skipFull bool) (*Params, error) {
	if limit > MaxLimit {
		return nil, fmt.Errorf("page limit can't be more than %d", MaxLimit)
	}

	if limit < 0 {
		limit = 0
	}
	if skip < 0 {
		skip = 0
	}

	return &Params{Skip: skip, Limit: limit, SkipEmpty: skipEmpty, SkipFull: skipFull}, nil
}

func NewParamsLimitOnly(limit int64) (*Params, error) {
	return NewParams(0, limit, false, false)
}
