package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/service"
	"net/http"
)

// Deprecated
// @ID checkUsername
// @Description if username does not exist it will return ok status with message username already taken
// @Produce  json
// @Param username query string true "Username to check"
// @Success 200 {object} OkRes
// @Failure 500 {object} ErrorRes
// @Router /users/check-username [get]
func CheckUsernameUniqueness(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	existsUsername, err := db.ExistsUsername(c.Request.Context(), username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	if existsUsername {
		c.JSON(http.StatusOK, Ok("username already exists"))
		return
	} else {
		c.JSON(http.StatusOK, Ok("username is unique"))
		return
	}
}

// @ID checkUsernameV2
// @Description if username does not exist it will return ok status with message username already taken
// @Param username path string true "Username to check"
// @Success 200 {object} OkRes
// @Failure 500 {object} ErrorRes
// @Router /profiles/{username}/check [get]
func CheckUsernameUniquenessV2(c *gin.Context) {
	username := c.Param("username")
	existsUsername, err := db.ExistsUsername(c.Request.Context(), username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	if existsUsername {
		c.JSON(http.StatusOK, Ok("username already exists"))
		return
	} else {
		c.JSON(http.StatusOK, Ok("username is unique"))
		return
	}
}

type ChangeUsernameInput struct {
	Username string
}

// @ID change user username
// @Produce json
// @Param input body ChangeUsernameInput true "input body"
// @Success 200 {object} OkRes
// @Failure 500 {object} ErrorRes
// @Router /users/me/username/change [put]
func ChangeUsername(c *gin.Context) {
	var input ChangeUsernameInput
	err := c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err = service.ChangeUsername(c.Request.Context(), input.Username, auth.Id(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	c.JSON(http.StatusOK, Ok("changed username"))
}
