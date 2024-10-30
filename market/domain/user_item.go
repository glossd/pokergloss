package domain

import "time"

type UserItem struct {
	ItemID ItemID
	// Unix Seconds
	ExpiresAt int64
	Forever bool

	selected bool
}

func NewUserItem(cmd IAddItemCommand) *UserItem {
	return &UserItem{ItemID: cmd.GetItemID(), ExpiresAt: time.Now().Add(cmd.Duration()).Unix()}
}

func newUserItemForever(item *Item) *UserItem {
	return &UserItem{ItemID: item.ID, Forever: true}
}

func (up *UserItem) IncreaseExpiration(d time.Duration) {
	if up.Forever {
		return
	}
	seconds := d.Seconds()
	if up.IsExpired() {
		up.ExpiresAt = Now().Add(d).Unix()
	} else {
		up.ExpiresAt += int64(seconds)
	}
}

func (up *UserItem) IsExpired() bool {
	if up.Forever {
		return false
	}
	return up.ExpiresAt < Now().Unix()
}

func (up *UserItem) IsSelected() bool {
	return up.selected
}
