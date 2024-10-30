package service

import (
	"context"
	"github.com/glossd/pokergloss/survival/db"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func DecCardTickets(ctx context.Context, userID string) (bool, error) {
	err := db.CardDecTicket(ctx, userID)
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			for _, writeError := range err.WriteErrors {
				if writeError.Message == "Document failed validation" {
					return false, nil
				}
			}
		}
		log.Errorf("Failed to decrement card tickets: %s", err)
		return false, err
	}
	return true, nil
}
