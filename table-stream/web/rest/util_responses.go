package rest

import (
	"github.com/gin-gonic/gin"
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
	c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
}
