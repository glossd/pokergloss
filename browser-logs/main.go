package browserlogs

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/browser-logs/web/router"
)

// @title Browser Logs API
// @schemes https
// @license.name Pokerblow
// @host pokerblow.com
// @BasePath /api/browser-logs
func Run(c *gin.Engine) func(context.Context) {
	auth.Init()
	router.New(c)
	return func(ctx context.Context) {}
}
