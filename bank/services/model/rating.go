package model

import "github.com/glossd/pokergloss/bank/domain"

type Rating struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Picture  string `json:"picture"`
	Chips    int64  `json:"chips"`
	Rank     int64  `json:"rank"`
}

func ToRatings(r []*domain.Rating) []*Rating {
	var ratings []*Rating
	for _, rating := range r {
		ratings = append(ratings, ToRating(rating))
	}
	return ratings
}

func ToRating(r *domain.Rating) *Rating {
	return &Rating{
		UserID:   r.UserID,
		Username: r.Username,
		Picture:  r.Picture,
		Chips:    r.Chips,
		Rank:     r.Rank,
	}
}
