package router

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/table/web/client/metrics"
	"github.com/glossd/pokergloss/table/web/rest"
)

const BasePath = "/api/table"

func New(r *gin.Engine) *gin.Engine {
	base := r.Group(BasePath)

	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	ginpprof.WrapGroup(base)

	base.GET("/metrics", metrics.ExposeMetrics)

	base.GET("/tables", rest.GetTables)

	base.GET("/sit-n-go/lobbies", rest.GetSitAndGoLobbies)
	base.GET("/multi/lobbies", rest.GetMultiLobbies)

	base.GET("/tables/:id", rest.GetFullTable)
	base.GET("/sit-n-go/lobbies/:id", rest.GetFullSitngo)
	base.GET("/multi/lobbies/:id", rest.GetFullLobbyMulti)

	verified := base.Group("/")
	verified.Use(auth.EmailVerifiedMiddleware)
	authenticated := base.Group("/")
	authenticated.Use(auth.Middleware)

	verified.POST("/tables", rest.PostTable)

	verified.POST("/tables/:id/seats/:position/reserve", rest.ReserveSeat)
	verified.PUT("/tables/:id/seats/:position/buy-in", rest.BuyInV2)
	verified.DELETE("/tables/:id/seats/:position/stand", rest.Stand)

	// for anonymous survival users
	authenticated.PUT("/tables/:id/seats/:position/actions/:action", rest.MakeActionV2)
	authenticated.PUT("/tables/:id/seats/:position/sit-back", rest.SitBack)

	verified.PUT("/tables/:id/seats/:position/showdown-actions/:action", rest.MakeShowDownAction)
	verified.PUT("/tables/:id/seats/:position/add-chips", rest.AddChips)
	verified.DELETE("/tables/:id/seats/:position/cancel", rest.CancelSeatReservation)

	verified.GET("/tables/:id/seats/:position/config", rest.GetConfig)
	verified.PUT("/tables/:id/seats/:position/configs/auto-muck", rest.SetAutoMuckPlayerConfig)
	verified.PUT("/tables/:id/seats/:position/configs/auto-top-up", rest.SetAutoTopUpPlayerConfig)
	verified.PUT("/tables/:id/seats/:position/configs/auto-rebuy", rest.SetAutoRebuyPlayerConfig)

	verified.GET("/configs", rest.GetPlayerAutoConfig)
	verified.PUT("/configs", rest.SetPlayerAutoConfig)

	verified.PUT("/tables/:id/seats/:position/intent", rest.PutIntent)
	verified.DELETE("/tables/:id/seats/:position/intent", rest.DeleteIntent)

	verified.POST("/sit-n-go/lobbies", rest.PostSitAndGoLobby)
	verified.PUT("/sit-n-go/lobbies/:id/register", rest.RegisterForSitAndGo)
	verified.PUT("/sit-n-go/lobbies/:id/unregister", rest.UnregisterFromSitAndGo)

	verified.PUT("/multi/lobbies/:id/register", rest.RegisterMulti)
	verified.PUT("/multi/lobbies/:id/unregister", rest.UnregisterMulti)
	return r
}
