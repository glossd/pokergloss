package services

import (
	"context"
"github.com/glossd/pokergloss/auth/authid"
"go.mongodb.org/mongo-driver/bson/primitive"
)
type IdenParams struct {
	ctx  context.Context
	id   primitive.ObjectID
	iden authid.Identity
}

func NewIdenParams(ctx context.Context, tableID string, iden authid.Identity) (*IdenParams, error) {
	oid, err := primitive.ObjectIDFromHex(tableID)
	if err != nil {
		return nil, E(err)
	}
	return &IdenParams{ctx: ctx, id: oid, iden: iden}, nil
}

func (p *IdenParams) GetCtx() context.Context {
	return p.ctx
}

func (p *IdenParams) GetID() primitive.ObjectID {
	return p.id
}

func (p *IdenParams) GetIden() authid.Identity {
	return p.iden
}
