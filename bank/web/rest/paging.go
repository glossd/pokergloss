package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/bank/db/paging"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func GetPage(c *gin.Context) paging.Page {
	page := paging.DefaultPage()
	skip, ok := parseInt(c, "skip")
	if ok {
		page.Skip = skip
	}
	limit, ok := parseInt(c, "limit")
	if ok {
		page.Limit = limit
	}
	return page
}

func parseInt(c *gin.Context, param string) (int64, bool) {
	paramStr := c.Request.URL.Query().Get(param)
	if len(paramStr) > 0 {
		result, err := strconv.Atoi(paramStr)
		if err != nil {
			log.Warnf("Couldn't parse %s to int", param)
			return -1, false
		}
		return int64(result), true
	}
	return -1, false
}
