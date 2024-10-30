package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/messenger/service"
	"github.com/glossd/pokergloss/messenger/web/model"
	"net/http"
	"strconv"
)

type GetChatMessagesResponse struct {
	Messages []*model.Message `json:"messages"`
}

// @ID get chat messages
// @Param lastId query int false "from which message id to start searching messages"
// @Param limit query int false "number of ratings to return, max 20"
// @Param chatId path string true "chat id"
// @Success 200 {object} GetChatMessagesResponse
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /chats/{chatId}/messages [get]
func GetChatMessages(c *gin.Context) {
	chatID := c.Param("chatId")
	id := c.Query("lastId")
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 0
	}
	messages, err := service.GetChatMessages(c.Request.Context(), auth.Id(c), chatID, id, int64(limit))
	if err != nil {
		handleError(c, err)
		return
	}
	var result []*model.Message
	for _, message := range messages {
		result = append(result, model.ToMessage(message))
	}

	c.JSON(http.StatusOK, GetChatMessagesResponse{Messages: result})
}
