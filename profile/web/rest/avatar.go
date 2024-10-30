package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/profile/service"
	"net/http"
)

// max avatar size https://meta.stackexchange.com/a/305279/457459
const maxAvatarSize = 1 << 18 // max 256kb
const maxAvatarSizeKb = maxAvatarSize / 1024

type UploadAvatarOutput struct {
	PhotoURL string `json:"photoURL"`
}

// @ID upload avatar
// @Accept mpfd
// @Produce json
// @Param avatar formData file true "avatar image"
// @Success 200 {object} UploadAvatarOutput
// @Failure 400 {object} ErrorRes
// @Router /upload-avatar [post]
func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	// todo check on format jpeg, png ...

	if file.Size > maxAvatarSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(fmt.Errorf("avatar image size can't be more than %dkb", maxAvatarSizeKb)))
		return
	}

	photoURL, err := service.UpdateUserAvatar(c.Request.Context(), file, auth.Id(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	c.JSON(http.StatusOK, UploadAvatarOutput{PhotoURL: photoURL})
}
