package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sort"
)

type UserChatList struct {
	UserID string                       `bson:"_id"`
	U2UChats map[string]*U2UChatForList // string-uid
	IsSupportCreated bool
}

type U2UChatForList struct {
	ID                primitive.ObjectID
	OtherUserID string
	UnreadCount int64
	// Sets in service
	lastMessage       *Message
}

func (u2u *U2UChatForList) SetLastMessage(m *Message) {
	u2u.lastMessage = m
}

func (u2u *U2UChatForList) GetUpdatedAt() int64 {
	if u2u == nil || u2u.lastMessage == nil {
		return 0
	}
	return u2u.lastMessage.UpdatedAt
}

func (u2u *U2UChatForList) GetLastMessage() *Message {
	return u2u.lastMessage
}

func NewUserChatList(userID string) *UserChatList {
	return &UserChatList{
		UserID: userID,
		U2UChats:  make(map[string]*U2UChatForList),
	}
}

func (ucl *UserChatList) SetChatWith(chatID primitive.ObjectID, otherUserID string) {
	ucl.U2UChats[otherUserID] = NewChatForList(chatID, otherUserID)
}

func NewChatForList(chatID primitive.ObjectID, otherUserID string) *U2UChatForList {
	return &U2UChatForList{
		ID:          chatID,
		OtherUserID: otherUserID,
	}
}

func (ucl *UserChatList) GetSortedChats() []*U2UChatForList {
	list := make([]*U2UChatForList, 0, len(ucl.U2UChats))
	for _, c := range ucl.U2UChats {
		list = append(list, c)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].GetUpdatedAt() > list[j].GetUpdatedAt()
	})
	return list
}
