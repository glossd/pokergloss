package ws

type ReadEventType string
const (
	Typing ReadEventType = "typing"
)

type ReadEvent struct {
	Type ReadEventType `json:"type"`
	Payload TypingPayload `json:"payload"`
}

type TypingPayload struct {
	ChatID string `json:"chatId"`
	Text string `json:"text"`
}
