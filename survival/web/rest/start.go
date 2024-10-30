package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/survival/domain"
	"github.com/glossd/pokergloss/survival/service"
	"net/http"
)

type StartSurvivalResponse struct {
	TableID string `json:"tableId"`
}

// @ID start survival
// @Success 200 {object} StartSurvivalResponse
// @Success 201 {object} StartSurvivalResponse
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /start [post]
func StartSurvival(c *gin.Context) {
	res, err := service.Start(c.Request.Context(), auth.Id(c), domain.Params{})
	if err != nil {
		handleServiceError(c, err)
		return
	}
	if res.AlreadyStarted {
		c.JSON(http.StatusOK, StartSurvivalResponse{TableID: res.TableID})
	} else {
		c.JSON(http.StatusCreated, StartSurvivalResponse{TableID: res.TableID})
	}
}

// @ID start anonymous survival
// @Success 200 {object} StartSurvivalResponse
// @Success 201 {object} StartSurvivalResponse
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /anonymous/start [post]
func StartSurvivalAnonymous(c *gin.Context) {
	res, err := service.Start(c.Request.Context(), auth.Id(c), domain.Params{Anonymous: true})
	if err != nil {
		handleServiceError(c, err)
		return
	}
	if res.AlreadyStarted {
		c.JSON(http.StatusOK, StartSurvivalResponse{TableID: res.TableID})
	} else {
		c.JSON(http.StatusCreated, StartSurvivalResponse{TableID: res.TableID})
	}
}

// @ID start survival idle
// @Success 200 {object} StartSurvivalResponse
// @Success 201 {object} StartSurvivalResponse
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /start-idle [post]
func StartSurvivalIdle(c *gin.Context) {
	res, err := service.Start(c.Request.Context(), auth.Id(c), domain.Params{Idle: true})
	if err != nil {
		handleServiceError(c, err)
		return
	}
	if res.AlreadyStarted {
		c.JSON(http.StatusOK, StartSurvivalResponse{TableID: res.TableID})
	} else {
		c.JSON(http.StatusCreated, StartSurvivalResponse{TableID: res.TableID})
	}
}
