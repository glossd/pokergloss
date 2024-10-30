package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/bonus/services/daily"
	"net/http"
)

type DailyBonusResponse struct {
	IsBonusPresent bool        `json:"isBonusPresent"`
	Bonus          *DailyBonus `json:"bonus,omitempty"`
}

type DailyBonus struct {
	DayInARow int   `json:"dayInARow"`
	Chips     int64 `json:"chips"`
}

// @ID take daily bonus
// @Success 200 {object} DailyBonusResponse
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /daily-bonus [put]
func TakeDailyBonus(c *gin.Context) {
	bonus, err := daily.TakeDaily(c.Request.Context(), auth.Id(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	if bonus == nil {
		c.JSON(http.StatusOK, &DailyBonusResponse{IsBonusPresent: false})
	} else {
		c.JSON(http.StatusOK, &DailyBonusResponse{IsBonusPresent: true, Bonus: &DailyBonus{DayInARow: bonus.DayInARow, Chips: bonus.CalculateBonus()}})
	}
}
