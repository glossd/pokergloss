package domain

import (
	"github.com/glossd/pokergloss/auth/authid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	TableID string

	CreatedBy authid.Identity
	CreatedAt int64

	Text string
}
