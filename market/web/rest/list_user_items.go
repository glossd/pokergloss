package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/market/service"
	"github.com/glossd/pokergloss/market/web/rest/model"
	"net/http"
)

// @ID list user items
// @Param userId path string true "User ID"
// @Success 200 {array} model.UserItem
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /users/{userId}/items [get]
func ListUserItems(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, Efmt("userId can't be empty"))
		return
	}

	items, err := service.ListUserItems(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.ToUserItems(items))
}

// @ID list my items
// @Success 200 {array} model.UserItem
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /me/items [get]
func ListMyItems(c *gin.Context) {
	items, err := service.ListUserItems(c.Request.Context(), auth.Id(c).UserId)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.ToUserItems(items))
}
