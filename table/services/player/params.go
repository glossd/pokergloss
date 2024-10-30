package player

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/table/services"
	"github.com/glossd/pokergloss/table/services/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PositionParams struct {
	ctx      context.Context
	tableID  primitive.ObjectID
	position int
	iden     authid.Identity
}

func NewPositionParams(ctx context.Context, tableID string, position int, iden authid.Identity) (*PositionParams, error) {
	oid, err := primitive.ObjectIDFromHex(tableID)
	if err != nil {
		return nil, services.E(err)
	}
	return &PositionParams{ctx: ctx, tableID: oid, position: position, iden: iden}, nil
}

type ChipsParams struct {
	*PositionParams
	chips int64
}

func ToChipsParams(params *PositionParams, chips int64) (*ChipsParams, error) {
	if chips < 0 {
		return nil, services.ErrFormat("you must put positive number of chips, your chips=%d", chips)
	}

	return &ChipsParams{PositionParams: params, chips: chips}, nil
}

type IntentParams struct {
	*PositionParams
	model.Intent
}

func NewIntentParams(posParams *PositionParams, intent model.Intent) (*IntentParams, error) {
	if intent.Type == "" {
		return nil, services.ErrFormat("intent type is specified")
	}

	return &IntentParams{PositionParams: posParams, Intent: intent}, nil
}
