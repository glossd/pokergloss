package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/ws/db"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type CheckOrUpdateNotificationTokenInput struct {
	Token    string `json:"token"`
	Platform string `json:"platform" enums:"web"`
}

type CheckOrUpdateNotificationTokenResponse struct {
	Status TokenStatus `json:"status" enums:"added,replaced,same"`
}

type TokenStatus string

const (
	Added    TokenStatus = "added"
	Replaced TokenStatus = "replaced"
	Same     TokenStatus = "same"
)

// @ID checkOrUpdateNotificationToken
// @Param input body CheckOrUpdateNotificationTokenInput true "input body"
// @Success 200 {object} CheckOrUpdateNotificationTokenResponse
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /notification/tokens [put]
func CheckOrUpdateNotificationToken(c *gin.Context) {
	var input CheckOrUpdateNotificationTokenInput
	err := c.BindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}
	if input.Platform != "web" {
		c.AbortWithStatusJSON(http.StatusBadRequest, EStr("only web platform is supported"))
		return
	}

	if input.Token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, EStr("token can't be empty"))
		return
	}

	ctx := c.Request.Context()
	iden := auth.Id(c)
	n, err := db.FindNotification(ctx, iden.UserId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			n = &db.Notification{UserID: iden.UserId}
		} else {
			handleError(c, err)
			return
		}
	}

	if n.Web.Token == input.Token {
		c.JSON(http.StatusOK, &CheckOrUpdateNotificationTokenResponse{Status: Same})
		return
	}

	var status TokenStatus
	if n.Web.Token == "" {
		status = Added
	} else {
		status = Replaced
	}

	n.Web.Token = input.Token
	err = db.UpsertNotification(ctx, n)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, &CheckOrUpdateNotificationTokenResponse{Status: status})
}
