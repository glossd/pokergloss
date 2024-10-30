package mq

const ReselectTopicID = "pg.market.reselect"

type ReselectEvent struct {
	UserID string `json:"userId"`
}
