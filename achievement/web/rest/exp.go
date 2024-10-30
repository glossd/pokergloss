package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/achievement/model"
	"github.com/glossd/pokergloss/achievement/service"
	"github.com/glossd/pokergloss/auth"
	"net/http"
)

// @ID get my points
// @Success 200 {object} model.Exp
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /points/me [get]
func GetMyExp(c *gin.Context) {
	exp, err := service.GetUserExp(c.Request.Context(), auth.Id(c).UserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}
	c.JSON(http.StatusOK, model.ToExp(exp))
}

// @ID get user points
// @Param userId path string true "User ID"
// @Success 200 {object} model.Exp
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /users/{userId}/points [get]
func GetUserExp(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, Efmt("userId can't be empty"))
		return
	}
	exp, err := service.GetUserExp(c.Request.Context(), userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}
	c.JSON(http.StatusOK, model.ToExp(exp))
}
