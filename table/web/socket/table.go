package socket

import (
	"github.com/gin-gonic/gin"
)

// @Summary Websocket stream of table changes.
// @Success 200 {object} events.TableEvent
// @ID use websocket instead
// @Param id path string true "Table ID"
// @Param token query string true "Verification token"
// @Router /ws/tables/{id} [get]
func TableChangeStream(c *gin.Context) {
	// for clients
}
