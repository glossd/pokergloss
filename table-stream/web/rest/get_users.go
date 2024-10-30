package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/table-stream/web/model"
	"github.com/glossd/pokergloss/table-stream/web/ws"
	"net/http"
)

// @ID table users
// @Param id path string true "Table ID"
// @Success 200 {object} model.TableUsers
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /tables/{id}/users [get]
func GetTableUsers(c *gin.Context) {
	tableId := c.Param("id")
	users := ws.GetTableUsers(tableId)
	c.JSON(http.StatusOK, model.ToTableUsers(users))
}
