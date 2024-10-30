package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/profile/domain"
	"github.com/glossd/pokergloss/profile/service"
	"net/http"
)

type CustomToken struct {
	Token string `json:"custom_token"`
}

type SignUpInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`

	IP      string `json:"ip"`
	Lang    string `json:"lang"`
	OS      string `json:"os"`
	Browser string `json:"browser"`
}

// @ID signup
// @Accept x-www-form-urlencoded
// @Produce  json
// @Param input body SignUpInput true "Sing un params"
// @Success 200 {object} CustomToken
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /signup [post]
func Signup(c *gin.Context) {
	r := c.Request

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	lang := r.FormValue("lang")
	os := r.FormValue("os")
	ip := r.FormValue("ip")
	browser := r.FormValue("browser")

	user, err := service.CreateUser(r.Context(), username, email, password, domain.TechInfo{
		IP:      ip,
		Lang:    lang,
		Browser: browser,
		OS:      os,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	token, err := service.GetCustomToken(r.Context(), user.UID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	c.JSON(http.StatusOK, &CustomToken{Token: token})
}
