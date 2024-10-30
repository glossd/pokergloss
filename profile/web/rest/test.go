package rest

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func TestHeaders(c *gin.Context) {
	r := c.Request
	log.Infof("X-Forwarded-For %s, X-Real-Ip %s, RemoteAddr %s", r.Header.Get("X-Forwarded-For"), r.Header.Get("X-Real-Ip"), r.RemoteAddr)
}
