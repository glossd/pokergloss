package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/table-chat/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type PostMessageInput struct {
	Text string `json:"text"`
}

// @ID post message
// @Param tableId path string true "Table ID"
// @Param input body PostMessageInput true "New message"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{tableId}/messages [post]
func PostMessage(c *gin.Context) {
	var input PostMessageInput
	err := c.BindJSON(&input)
	if err != nil {
		log.Warn("Couldn't post message: parse input: %", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = service.PostMessage(c.Request.Context(), c.Param("tableId"), auth.Id(c), input.Text)
	if err != nil {
		if err == service.ErrEmpty {
			c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
			return
		} else if err == service.ErrUserBlackListed {
			c.AbortWithStatusJSON(http.StatusForbidden, E(err))
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
			return
		}
	}

	c.JSON(http.StatusCreated, Ok("message sent"))
}
