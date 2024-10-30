package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/achievement/service"
	"github.com/glossd/pokergloss/auth"
	"net/http"
)

// @ID get my achievements
// @Success 200 {array} model.Achievement
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /achievements/me [get]
func GetMyAchievements(c *gin.Context) {
	achievements, err := service.FindSortedAchievements(c.Request.Context(), auth.Id(c).UserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}
	c.JSON(http.StatusOK, achievements)
}

// @ID get user achievements
// @Param userId path string true "User ID"
// @Success 200 {array} model.Achievement
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /users/{userId}/achievements [get]
func GetUserAchievements(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, Efmt("userId can't be empty"))
		return
	}
	achievements, err := service.FindSortedAchievements(c.Request.Context(), userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}
	c.JSON(http.StatusOK, achievements)
}
