package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/survival/service"
	"net/http"
)

// @ID release survival
// @Success 201 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /release [delete]
func ReleaseSurvival(c *gin.Context) {
	err := service.Release(c.Request.Context(), auth.Id(c))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, Ok("released"))
}
