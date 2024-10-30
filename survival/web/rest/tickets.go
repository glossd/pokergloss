package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/survival/db"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type GetTicketsResponse struct {
	Tickets int64 `json:"tickets"`
}

// @ID get my tickets
// @Success 200 {object} GetTicketsResponse
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /my/tickets [get]
func GetTickets(c *gin.Context) {
	card, err := db.FindCard(c.Request.Context(), auth.Id(c).UserId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, &GetTicketsResponse{Tickets: 0})
			return
		}
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, &GetTicketsResponse{Tickets: card.Tickets})
}
