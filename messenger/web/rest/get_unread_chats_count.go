package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/messenger/service"
	"net/http"
)

type GetUnreadChatsCountResponse struct {
	Count int64 `json:"count"`
}

// @ID get unread chats count
// @Success 200 {object} GetUnreadChatsCountResponse
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /unread-chats-count [get]
func GetUnreadChatsCount(c *gin.Context) {
	count, err := service.GetUnreadChatsCount(c.Request.Context(), auth.Id(c))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, &GetUnreadChatsCountResponse{Count: count})
}
