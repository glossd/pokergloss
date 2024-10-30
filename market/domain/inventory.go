package domain

import (
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

var ErrNoItemSelected = E("no item is selected")
var ErrSelectItemExpired = E("selected item is expired")

type Inventory struct {
	UserID    string `bson:"_id"`
	Inventory map[ItemID]*UserItem
	SelectedItemID ItemID
	Version int64
}

func NewInventory(userID string) *Inventory {
	inv := &Inventory{UserID: userID, Inventory: make(map[ItemID]*UserItem)}
	inv.addItemForever(Invisible)
	inv.SelectedItemID = Invisible.ID
	return inv
}

func (i *Inventory) AddItem(cmd IAddItemCommand) {
	if cmd == nil {
		log.Errorf("domain.AddItem: command is nil")
		return
	}
	ui, ok := i.Inventory[cmd.GetItemID()]
	if ok {
		ui.IncreaseExpiration(cmd.Duration())
	} else {
		i.Inventory[cmd.GetItemID()] = NewUserItem(cmd)
	}
	i.SelectedItemID = cmd.GetItemID()
}

func (i *Inventory) addItemForever(item *Item) {
	i.Inventory[item.ID] = newUserItemForever(item)
}

func (i *Inventory) SelectItem(iID string) error {
	itemID := ItemID(iID)
	item, ok := i.Inventory[itemID]
	if !ok {
		return E("you don't have item %s", itemID)
	}
	if item.IsExpired() {
		return ErrSelectItemExpired
	}
	i.SelectedItemID = itemID
	return nil
}

func (i *Inventory) Reselect() bool {
	_, err := i.GetSelectedItem()
	if err == nil {
		log.Errorf("Failed to reselect, selected item: %s", err)
		return false
	}
	i.deleteExpired()
	if len(i.Inventory) == 0 {
		return false
	}
	i.SelectedItemID = i.coolestItem().ID
	return true
}

func (i *Inventory) coolestItem() *Item {
	store := make(map[SaleType][]*Item)
	for id := range i.Inventory {
		if item, ok := AnimatedItemsMap[id]; ok {
			store[ForCoins] = append(store[ForCoins], item)
			continue
		}
		if item, ok := ItemsNotForSaleMap[id]; ok {
			store[NotForSale] = append(store[NotForSale], item)
			continue
		}
		if item, ok := ItemsOnSaleMap[id]; ok {
			store[ForChips] = append(store[ForChips], item)
			continue
		}
	}
	if len(store[ForCoins]) > 0 {
		items := store[ForCoins]
		sort.Slice(items, func(i, j int) bool {
			return items[i].PriceList.Day > items[j].PriceList.Day
		})
		return items[0]
	}
	if len(store[NotForSale]) > 1 {
		items := store[NotForSale]
		sort.Slice(items, func(i, j int) bool {
			return ItemsNotForSaleOrder[items[i].ID] > ItemsNotForSaleOrder[items[j].ID]
		})
		return items[0]
	}
	if len(store[ForChips]) > 0 {
		items := store[ForChips]
		sort.Slice(items, func(i, j int) bool {
			return items[i].PriceList.Day > items[j].PriceList.Day
		})
		return items[0]
	}

	return Invisible
}

func (i *Inventory) deleteExpired() bool {
	var toDelete []ItemID
	for itemID, item := range i.Inventory {
		if item.IsExpired() {
			toDelete = append(toDelete, itemID)
		}
	}
	for _, id := range toDelete {
		delete(i.Inventory, id)
	}
	return len(toDelete) > 0
}

func (i *Inventory) AvailableItems() []*UserItem {
	items := make([]*UserItem, 0, len(i.Inventory))
	for _, item := range i.Inventory {
		if !item.IsExpired() {
			if item.ItemID == i.SelectedItemID {
				item.selected = true
			}
			items = append(items, item)
		}
	}
	return items
}

func (i *Inventory) GetSelectedItem() (*UserItem, error) {
	if i.SelectedItemID == "" {
		log.Errorf("GetSelectedItem: selected item is empty")
		return nil, ErrNoItemSelected
	}

	ui, ok := i.Inventory[i.SelectedItemID]
	if !ok {
		log.Errorf("GetSelectedItem: selected item is not in inventory")
		return nil, ErrNoItemSelected
	}
	if ui.IsExpired() {
		return nil, ErrSelectItemExpired
	}
	return ui, nil
}

func (i *Inventory) GetSelectedItemUnsafe() *UserItem {
	item, _ := i.Inventory[i.SelectedItemID]
	return item
}

type IAddItemCommand interface {
	GetItemID() ItemID
	Duration() time.Duration
}
