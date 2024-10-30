package metrics

import (
	"github.com/gin-gonic/gin"
)

func ExposeMetrics(c *gin.Context) {
	c.JSON(200, gin.H{
		"up": 1,
	})
}
