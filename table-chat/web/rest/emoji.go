package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/table-chat/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type PostEmojiInput struct {
	Emoji string `json:"emoji" enums:"joy,wink,cry,rage,like,scream,sunglasses,raisedEyebrow"`
}

// @ID post emoji
// @Param tableId path string true "Table ID"
// @Param input body PostEmojiInput true "New emoji"
// @Success 200 {object} OkRes
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{tableId}/emojis [post]
func PostEmoji(c *gin.Context) {
	var input PostEmojiInput
	err := c.BindJSON(&input)
	if err != nil {
		log.Warn("Couldn't post emoji: parse input: %", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = service.PostEmoji(c.Request.Context(), c.Param("tableId"), auth.Id(c), input.Emoji)
	if err != nil {
		if err == service.ErrEmpty || err == service.ErrInvalid {
			c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		}
		return
	}

	c.JSON(http.StatusCreated, Ok("emoji sent"))
}
