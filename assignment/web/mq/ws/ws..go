package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/assignment/domain"
	"github.com/glossd/pokergloss/assignment/model"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/gomq/mqws"
	log "github.com/sirupsen/logrus"
)

type AssignmentEvent struct {
	Type    string      `json:"type" enums:"changedAssignment"`
	Payload interface{} `json:"payload"`
}

// @ID events from websocket
// @Success 200 {array} AssignmentEvent
// @Router /ws [get]
func UseWS(c *gin.Context) {}

func PublishAssignmentEvent(userID string, a *domain.UserAssignment) {
	if !conf.IsProd() {
		return
	}
	err := mqws.Publish(&mqws.Message{
		EntityType: mqws.Message_USER,
		EntityId:   userID,
		Events:     []*mqws.Event{{Type: "changedAssignment", Payload: gomq.M{"assignment": model.ToAssignment(a)}.JSON()}},
	})
	if err != nil {
		log.Errorf("Failed to publish exp event: %s", err)
	}
}
