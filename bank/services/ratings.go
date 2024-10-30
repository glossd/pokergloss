package services

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/domain"
	"github.com/glossd/pokergloss/bank/services/model"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
)

func GetRating(ctx context.Context, userID string) (*domain.Rating, error) {
	return db.FindRating(ctx, userID)
}

func GetRatingPage(ctx context.Context, iden *authid.Identity, page model.PageRequest) (*model.PageRating, error) {
	if page.PageNumber == 0 {
		if iden != nil {
			return GetUserRatingPage(ctx, iden.UserId, page.PageSize)
		} else {
			return GateRatingPageByNumber(ctx, model.PageRequest{PageNumber: 1, PageSize: page.PageSize})
		}
	}
	return GateRatingPageByNumber(ctx, page)
}

func GateRatingPageByNumber(ctx context.Context, page model.PageRequest) (*model.PageRating, error) {
	pageSize := page.PageSize
	pageNumber := page.PageNumber
	pageRatings, err := db.FindRatings(ctx, (pageNumber-1)*pageSize, pageSize)
	if err != nil {
		return nil, err
	}
	count, err := db.CountRatings(ctx)
	if err != nil {
		return nil, err
	}
	return &model.PageRating{
		Ratings:     model.ToRatings(pageRatings),
		PageNumber:  page.PageNumber,
		PageCount:   getPageCount(count, pageSize),
		UserPageIdx: -1,
	}, nil
}

func GetUserRatingPage(ctx context.Context, userID string, pageSize int64) (*model.PageRating, error) {
	userRating, err := db.FindRating(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return GateRatingPageByNumber(ctx, model.PageRequest{PageNumber: 1, PageSize: pageSize})
		}
		return nil, err
	}
	count, err := db.CountRatings(ctx)
	if err != nil {
		return nil, err
	}

	pageNumber := (userRating.Rank-1)/pageSize + 1
	pageRatings, err := db.FindRatings(ctx, (pageNumber-1)*pageSize, pageSize)
	if err != nil {
		return nil, err
	}

	pageCount := getPageCount(count, pageSize)
	userPageIdx := (userRating.Rank - 1) % pageSize

	return &model.PageRating{
		Ratings:     model.ToRatings(pageRatings),
		PageNumber:  pageNumber,
		PageCount:   pageCount,
		UserPageIdx: userPageIdx,
	}, nil
}

func getPageCount(count int64, pageSize int64) int64 {
	pageCount := int64(math.Ceil(float64(count) / float64(pageSize)))
	return pageCount
}
