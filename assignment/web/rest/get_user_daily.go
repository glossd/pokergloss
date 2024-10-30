package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/assignment/model"
	"github.com/glossd/pokergloss/assignment/service"
	"github.com/glossd/pokergloss/auth"
	"net/http"
)

// @ID get my daily assignments
// @Success 200 {array} model.Assignment
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /my/daily/assignments [get]
func GetMyDailyAssignments(c *gin.Context) {
	ud, err := service.FindUserDaily(c.Request.Context(), auth.Id(c).UserId)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.FromUserDaily(ud))
}
