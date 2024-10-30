package rest

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/survival/db"
	"github.com/glossd/pokergloss/survival/domain"
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

func Efmt(msg string, a ...interface{}) *ErrorRes {
	return &ErrorRes{Message: fmt.Sprintf(msg, a...)}
}

type OkRes struct {
	Message string `json:"message"`
}

func Ok(msg string) *OkRes {
	return &OkRes{Message: msg}
}

func handleServiceError(c *gin.Context, err error) {
	if _, ok := err.(*domain.Err); ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		c.AbortWithStatusJSON(http.StatusNotFound, E(err))
		return
	}

	if errors.Is(err, db.ErrVersionNotMatch) {
		c.AbortWithStatusJSON(http.StatusConflict, Efmt("dirty write"))
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
}
