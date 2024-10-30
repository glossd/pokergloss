package model

import (
	"strconv"
)

type PageRequest struct {
	PageNumber int64
	PageSize int64
}

type PageRating struct {
	Ratings []*Rating `json:"ratings"`
	PageNumber int64 `json:"pageNumber"`
	PageCount int64 `json:"pageCount"`
	UserPageIdx int64 `json:"userPageIdx"`
}

func NewPageRequest(pageSize, pageNumber string) (*PageRequest, error) {
	var p PageRequest
	ps, err := strconv.Atoi(pageSize)
	if err != nil {
		return nil, err
	}
	p.PageSize = int64(ps)

	pn, err := strconv.Atoi(pageNumber)
	if err != nil {
		return &p, nil
	}
	p.PageNumber = int64(pn)

	return &p, nil
}
