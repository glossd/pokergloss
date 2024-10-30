package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/table/services"
	"github.com/glossd/pokergloss/table/services/paging"
	"github.com/glossd/pokergloss/table/services/player"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func PositionParams(c *gin.Context) (*player.PositionParams, error) {
	position, err := strconv.Atoi(c.Param("position"))
	if err != nil {
		return nil, err
	}
	return player.NewPositionParams(c.Request.Context(), c.Param("id"), position, auth.Id(c))
}

func PositionChipsParams(c *gin.Context, chips int64) (*player.ChipsParams, error) {
	posParams, err := PositionParams(c)
	if err != nil {
		return nil, err
	}
	params, err := player.ToChipsParams(posParams, chips)
	if err != nil {
		return nil, err
	}
	return params, nil
}

func IdenParams(c *gin.Context) (*services.IdenParams, error) {
	return services.NewIdenParams(c.Request.Context(), c.Param("id"), auth.Id(c))
}

func PageParams(c *gin.Context) (*paging.Params, error) {
	skipP := parseInt(c, "skip")
	skip := 0
	if skipP != nil {
		skip = *skipP
	}

	limitP := parseInt(c, "limit")
	limit := 20
	if limitP != nil {
		limit = *limitP
	}

	skipEmptyP := parseBool(c, "skipEmpty")
	skipEmpty := false
	if skipEmptyP != nil {
		skipEmpty = *skipEmptyP
	}

	skipFullP := parseBool(c, "skipFull")
	skipFull := false
	if skipFullP != nil {
		skipFull = *skipFullP
	}

	return paging.NewParams(int64(skip), int64(limit), skipEmpty, skipFull)
}

func parseInt(c *gin.Context, param string) *int {
	paramStr := c.Request.URL.Query().Get(param)
	if len(paramStr) > 0 {
		result, err := strconv.Atoi(paramStr)
		if err != nil {
			log.Warnf("Couldn't parse %s to int", param)
			return nil
		}
		return &result
	}
	return nil
}

func parseBool(c *gin.Context, param string) *bool {
	paramStr := c.Request.URL.Query().Get(param)
	if len(paramStr) > 0 {
		result, err := strconv.ParseBool(paramStr)
		if err != nil {
			log.Warnf("Couldn't parse %s to bool", param)
			return nil
		}
		return &result
	}
	return nil
}
