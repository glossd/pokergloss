package wsstore

import "github.com/gin-gonic/gin"

type MessengerEventType string
const (
	NewMessage MessengerEventType = "messengerNewMessage"
	MessageStatus MessengerEventType = "messengerMessageStatus"
	NewChat MessengerEventType = "messengerNewChat"
	Typing MessengerEventType = "messengerTyping"
)

type MessengerEvent struct {
	Type MessengerEventType `json:"type" enums:"messengerNewMessage,messengerMessageStatus,messengerNewChat"`
	Payload map[string]interface{} `json:"payload"`
}

// @ID use websocket
// @Success 200 {object} MessengerEvent
// @Router /u2u/chats [post]
func UseWS(c *gin.Context) {}
