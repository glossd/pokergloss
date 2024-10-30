package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/survival/db"
	"github.com/glossd/pokergloss/survival/web/model"
	"net/http"
)

// @ID get user score
// @Param userId path string true "User ID"
// @Success 200 {object} model.Score
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /users/{userId}/score [get]
func GetUserScore(c *gin.Context) {
	score, err := db.FindScoreOrDefault(c.Request.Context(), c.Param("userId"))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.Score{Level: score.Level})
}
