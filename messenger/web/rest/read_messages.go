package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/messenger/service"
	"net/http"
)

type ReadMessagesInput struct {
	MessageIDs []string `json:"messageIds"`
}

// @ID read chat messages
// @Param chatId path string true "chat id"
// @Param input body ReadMessagesInput true "input body"
// @Success 201 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /chats/{chatId}/read/messages [put]
func ReadMessages(c *gin.Context) {
	var input ReadMessagesInput
	err := c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	err = service.ReadMessages(c.Request.Context(), auth.Id(c), c.Param("chatId"), input.MessageIDs)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, Ok("read"))
}
