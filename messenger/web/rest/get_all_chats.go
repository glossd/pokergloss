package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/messenger/service"
	"net/http"
)

// @ID get all chats
// @Success 200 {array} model.Chat
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /chats [get]
func GetAllChats(c *gin.Context) {
	chats, err := service.GetUserChats(c.Request.Context(), auth.Id(c))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, chats)
}
