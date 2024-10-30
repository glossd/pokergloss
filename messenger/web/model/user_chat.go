package model

type Chat struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Picture string `json:"picture"`
	LastMessage *Message `json:"lastMessage"`
	IsPhantom bool `json:"isPhantom"`
}
