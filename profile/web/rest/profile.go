package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/service"
	"github.com/glossd/pokergloss/profile/web/model"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// @ID get profile
// @Param username path string true "Profile Username"
// @Success 200 {object} model.Profile
// @Failure 400 {object} ErrorRes
// @Failure 404 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /profiles/{username} [get]
func GetProfile(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, Efmt("username can't be empty"))
		return
	}
	user, err := service.GetProfile(c.Request.Context(), username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(http.StatusNotFound, Efmt("user with such username doesn't exist"))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	c.JSON(http.StatusOK, model.ToProfile(user))
}

// @ID search profiles
// @Param username path string true "Profile Username Search String"
// @Success 200 {array} model.Profile
// @Failure 400 {object} ErrorRes
// @Failure 404 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /profiles/{username}/search [get]
func SearchProfiles(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, Efmt("username can't be empty"))
		return
	}

	profiles, err := db.SearchProfiles(c.Request.Context(), username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(http.StatusNotFound, Efmt("users with such username don't exist"))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	c.JSON(http.StatusOK, model.ToProfiles(profiles))
}
