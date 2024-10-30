package router

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/profile/conf"
	"github.com/glossd/pokergloss/profile/web/rest"
	"net/http"
)

const BasePath = "/api/profile"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)

	base.GET("/status", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"up": 1}) })

	base.POST("/signup", rest.Signup)
	base.GET("/users/check-username", rest.CheckUsernameUniqueness)
	base.GET("/profiles/:username/check", rest.CheckUsernameUniquenessV2)
	base.GET("/profiles/:username", rest.GetProfile)
	base.GET("/profiles/:username/search", rest.SearchProfiles)

	auth.InitCustomSetup(conf.AuthClient)
	authenticated := base.Group("/", auth.Middleware)
	authenticated.POST("/upload-avatar", rest.UploadAvatar)
	authenticated.PUT("/users/me/username/change", rest.ChangeUsername)
	authenticated.GET("/djsl-test-ip-headers", rest.TestHeaders)

	return r
}
