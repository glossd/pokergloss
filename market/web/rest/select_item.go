package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/market/service"
	"net/http"
)

// @ID select item
// @Param itemId path string true "User's Item ID"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /items/{itemId}/select [put]
func SelectItem(c *gin.Context) {
	itemID := c.Param("itemId")
	err := service.SelectUserItem(c.Request.Context(), auth.Id(c).UserId, itemID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, Ok("selected"))
}
