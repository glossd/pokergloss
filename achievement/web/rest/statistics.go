package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/achievement/model"
	"github.com/glossd/pokergloss/achievement/service"
	"github.com/glossd/pokergloss/auth"
	"net/http"
)

// @ID get my statistics
// @Success 200 {object} model.Statistics
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /statistics/me [get]
func GetMyStatistics(c *gin.Context) {
	exp, err := service.GetUserStatistics(c.Request.Context(), auth.Id(c).UserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}
	c.JSON(http.StatusOK, model.ToStatistics(exp))
}

// @ID get user statistics
// @Param userId path string true "User ID"
// @Success 200 {object} model.Statistics
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /users/{userId}/statistics [get]
func GetUserStatistics(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, Efmt("userId can't be empty"))
		return
	}
	stat, err := service.GetUserStatistics(c.Request.Context(), userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}
	c.JSON(http.StatusOK, model.ToStatistics(stat))
}
