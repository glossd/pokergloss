package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services"
	"github.com/glossd/pokergloss/table/web/client/bankclient"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// Defines http error response
type ErrorRes struct {
	Message string `json:"message"`
}

func E(err error) *ErrorRes {
	return &ErrorRes{Message: err.Error()}
}

func EStr(msg string) *ErrorRes {
	return &ErrorRes{Message: msg}
}

type OkRes struct {
	Message string `json:"message"`
}

func Ok() *OkRes {
	return &OkRes{Message: "ok"}
}

func handleServiceError(c *gin.Context, err error) {
	if _, ok := err.(*domain.PokerErr); ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	if _, ok := err.(*services.ServiceErr); ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		c.AbortWithStatusJSON(http.StatusNotFound, E(err))
		return
	}

	if errors.Is(err, db.ErrVersionNotMatch) {
		c.AbortWithStatusJSON(http.StatusConflict, EStr("dirty write"))
		return
	}

	if errors.Is(err, bankclient.ErrNotEnoughChips) {
		c.AbortWithStatusJSON(http.StatusPreconditionFailed, E(err))
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
}
