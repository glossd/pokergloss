package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/bank/db"
	"github.com/glossd/pokergloss/bank/services"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type GetBalanceOutput struct {
	Chips int64 `json:"chips"`
	Coins int64 `json:"coins"`
}

// @ID get balance
// @Success 200 {object} GetBalanceOutput
// @Failure 400 {object} ErrorRes
// @Router /balance [get]
func GetBalance(c *gin.Context) {
	iden := auth.Id(c)
	balance, err := db.FindBalance(c.Request.Context(), iden.UserId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusOK, GetBalanceOutput{})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	c.JSON(http.StatusOK, GetBalanceOutput{Chips: balance.Chips, Coins: balance.Coins})
}

// @ID get user balance with rank
// @Param userId path string true "User ID"
// @Success 200 {object} model.BalanceWithRank
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /balances/{userId} [get]
func GetUserBalanceWithRank(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, Efmt("userId can't be empty"))
		return
	}
	b, err := services.GetBalanceWithRank(c.Request.Context(), userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	c.JSON(http.StatusOK, b)
}
