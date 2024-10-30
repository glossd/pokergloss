package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PurchaseItemCommand struct {
	ID        primitive.ObjectID `bson:"_id"`
	AddItemCommand `bson:",inline"`

	UnitPrice int64
	CreatedAt int64
	item      *Item
}

type AddItemCommand struct {
	UserID string
	ItemID    ItemID
	Units     int64
	TimeFrame TimeFrame
}

func NewGiftItemCommand(userID string, itemID ItemID, units int64, tf TimeFrame) (*AddItemCommand, error) {
	_, ok := ItemsOnSaleMap[itemID]
	if !ok {
		_, ok = ItemsNotForSaleMap[itemID]
		if !ok {
			return nil, errNoSuchItem(itemID)
		}
	}

	return &AddItemCommand{
		UserID:    userID,
		ItemID:    itemID,
		Units:     units,
		TimeFrame: tf,
	}, nil
}

func NewPurchaseItemCommand(userID string, itemID ItemID, units int64, tf TimeFrame) (*PurchaseItemCommand, error) {
	item, ok := ItemsOnSaleMap[itemID]
	if !ok {
		return nil, errNoSuchItem(itemID)
	}
	price, err := item.GetPrice(tf)
	if err != nil {
		return nil, err
	}
	return &PurchaseItemCommand{
		ID:        primitive.NewObjectID(),
		AddItemCommand: AddItemCommand{
			UserID:    userID,
			ItemID:    itemID,
			Units:     units,
			TimeFrame: tf,
		},
		UnitPrice: price,
		CreatedAt: time.Now().Unix(),
		item:      item,
	}, nil
}

func (pic AddItemCommand) Duration() time.Duration {
	return time.Duration(pic.Units) * pic.TimeFrame.Duration()
}

func (pic AddItemCommand) GetItemID() ItemID {
	return pic.ItemID
}

func (pic AddItemCommand) IsForCoins() bool {
	return ItemsOnSaleMap[pic.ItemID].SaleType == ForCoins
}

func (pic *PurchaseItemCommand) TotalPrice() int64 {
	return pic.Units * pic.UnitPrice
}

func (pic *PurchaseItemCommand) GetItem() *Item {
	return pic.item
}

