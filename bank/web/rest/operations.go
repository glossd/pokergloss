package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/bank/services"
	"github.com/glossd/pokergloss/bank/services/model"
	"net/http"
)

// @ID get operations
// @Param skip query int true "entities to skip"
// @Param limit query int true "number of entities to return"
// @Success 200 {array} model.Operation
// @Failure 500 {object} ErrorRes
// @Router /operations [get]
func GetOperations(c *gin.Context) {
	operations, err := services.ListOperations(c.Request.Context(), auth.Id(c), GetPage(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	c.JSON(http.StatusOK, model.ToOperations(operations))
}
