package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/messenger/service"
	"net/http"
)

type PostU2UChatInput struct {
	UserID string `json:"userId"`
}

// @ID post chat
// @Param input body PostU2UChatInput true "input body"
// @Success 201 {object} model.Chat
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /u2u/chats [post]
func PostU2UChat(c *gin.Context) {
	var input PostU2UChatInput
	err := c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	chat, err := service.CreateU2UChat(c.Request.Context(), auth.Id(c), input.UserID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, chat)
}
