package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/messenger/service"
	"github.com/glossd/pokergloss/messenger/web/model"
	"net/http"
)

type SendChatMessageInput struct {
	Text string `json:"text"`
}

// @ID send chat message
// @Param chatId path string true "chat id"
// @Param input body SendChatMessageInput true "input body"
// @Success 201 {object} model.Message
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /chats/{chatId}/messages [post]
func SendChatMessage(c *gin.Context) {
	var input SendChatMessageInput
	err := c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	msg, err := service.SendMessageToChat(c.Request.Context(), auth.Id(c), c.Param("chatId"), input.Text)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, model.ToMessage(msg))
}
