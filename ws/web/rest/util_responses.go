package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
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

func handleError(c *gin.Context, err error) {
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.AbortWithStatusJSON(http.StatusNotFound, E(err))
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
}
