package model

import (
	"github.com/glossd/pokergloss/bank/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Operation struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Type        domain.OperationType `json:"type" enums:"deposit,withdraw"`
	Reason      domain.Reason        `json:"reason" enums:"bonus,admin,cacheGame"`
	Chips       int64                `json:"chips"`
	UserID      string               `json:"userId"`
	Description string               `json:"description"`
	CreatedAt   int64                `json:"createdAt"`
}

func ToOperations(ops []*domain.Operation) []*Operation {
	result := make([]*Operation, 0, len(ops))
	for _, op := range ops {
		result = append(result, ToOperation(op))
	}
	return result
}

func ToOperation(op *domain.Operation) *Operation {
	return &Operation{
		ID:          op.ID,
		Type:        op.Type,
		Reason:      op.Reason,
		Chips:       op.Chips,
		UserID:      op.UserID,
		Description: op.Description,
		CreatedAt:   op.CreatedAt,
	}
}
