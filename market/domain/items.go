package domain

import (
	"math"
)

type ItemID string
const (
	PacifierID ItemID = "pacifier"

	SwordID ItemID = "sword"
	SunglassesID ItemID = "sunglasses"
	GrayFishID ItemID = "grayFish"
	MartiniID ItemID = "martini"
	SunriseCocktailID ItemID = "sunriseCocktail"
	CosmapolitanGlitchID ItemID = "cosmapolitanGlitch"
	LightBulbID ItemID = "lightBulb"
	IceCreamID ItemID = "iceCream"
	StrawberryID ItemID = "strawberry"
	SkullID ItemID = "skull"
	ShipID ItemID = "ship"
	FireExtinguisherID ItemID = "fireExtinguisher"
	InvisibleID ItemID = "invisible"

	// For coins
	GlassOfWineID ItemID = "glassOfWine"
	CloverID ItemID        = "clover"
	HourglassID ItemID     = "hourglass"
	PiggyBankID ItemID = "piggyBank"
	ButterflyID ItemID = "butterfly"
	MonsterID ItemID = "monster"
	GhostID ItemID = "ghost"
	BeeID ItemID = "bee"
	BurgerID ItemID = "burger"

	CrownID ItemID = "crown"
	// Survival
	TorturerID ItemID = "torturer"
	HellAmuletID ItemID = "hellAmulet"
	SmirkingDemonID ItemID = "smirkingDemon"
)

var Invisible = NewItemNotForSale(InvisibleID, "Invisible Item", Side)

var Pacifier = NewSideItemForChips(PacifierID, "Pacifier", 50)
var LightBulb = NewSideItemForChips(LightBulbID, "Light Bulb", 100)
var Slippers = NewSideItemForChips("slippers", "Slippers", 150)
var IceCream = NewSideItemForChips(IceCreamID, "Ice Cream", 200)
var GrayFish = NewSideItemForChips(GrayFishID, "Grey Fish", 250)
var Ship = NewSideItemForChips(ShipID, "Ship", 250)
var FireExtinguisher = NewSideItemForChips(FireExtinguisherID, "Fire extinguisher", 300)
var Sunglasses = NewSideItemForChips(SunglassesID, "Sunglasses", 300)
var Gamepad = NewSideItemForChips("gamepad", "Gamepad", 300)
var CupCoffee = NewSideItemForChips("cupCoffee", "Cup of coffee", 300)
var Strawberry = NewSideItemForChips(StrawberryID, "Strawberry", 350)
var CubicToyCat = NewSideItemForChips("cubicToyCat", "Cubic toy cat", 350)
var ScubaGlasses = NewSideItemForChips("scubaGlasses", "Scuba glasses", 350)
var Sword = NewSideItemForChips(SwordID, "Sword", 400)
var Voodoo = NewSideItemForChips("voodoo", "Voodoo", 450)
var Mummy = NewSideItemForChips("mummy", "Mummy", 450)
var Martini = NewSideItemForChips(MartiniID, "Martini", 500)
var Compass = NewSideItemForChips("compass", "Compass", 550)
var SunriseCocktail = NewSideItemForChips(SunriseCocktailID, "Sunrise cocktail", 600)
var CosmapolitanGlitch = NewSideItemForChips(CosmapolitanGlitchID, "Cosmapolitan glitch", 700)
var ToyRobot = NewSideItemForChips("toyRobot", "Toy robot", 750)
var DiverHelmet = NewSideItemForChips("diverHelmet", "Diver's helmet", 750)
var Dinosaur = NewSideItemForChips("dinosaur", "Dinosaur", 800)
var BacklitPumpkin = NewSideItemForChips("backlitPumpkin", "Pumpkin", 850)
var Skull = NewSideItemForChips(SkullID, "Skull", 900)
var Vendetta = NewSideItemForChips("vendetta", "Vendetta Mask", 900)
var Rocket = NewSideItemForChips("rocket", "Rocket", 950)
var Shuriken = NewSideItemForChips("shuriken", "Shuriken", 1100)
var Grenade = NewSideItemForChips("grenade", "Grenade", 1200)
var Rip = NewSideItemForChips("rip", "R.I.P", 1300)
var VintageCar = NewSideItemForChips("vintageCar", "Vintage Car", 1500)
var Tutanchamun = NewSideItemForChips("tutanchamun", "Tutanchamun", 3000)

var Burger = NewSideItemForChips(BurgerID, "Burger", 5000)
var Bee = NewSideItemForChips(BeeID, "Bee", 5000)
var PiggyBank = NewSideItemForChips(PiggyBankID, "Piggy Bank", 5000)
var Clover =  NewSideItemForChips(CloverID, "Clover", 5000)
var GlassOfWine = NewSideItemForChips(GlassOfWineID, "Glass of Wine", 5000)
var Monster = NewSideItemForChips(MonsterID, "Monster", 10000)
var Ghost = NewSideItemForChips(GhostID, "Ghost", 10000)
var Butterfly = NewSideItemForChips(ButterflyID, "Butterfly", 10000)
var GoldenGlasses = NewSideItemForChips("goldenGlasses", "Golden Glasses", 10000)
var Hourglass = NewSideItemForChips(HourglassID, "Hourglass", 150000)

var Crown = NewItemNotForSale(CrownID, "Crown", Top)

var Torturer = NewItemNotForSale(TorturerID, "Torturer", Side)
var HellAmulet = NewItemNotForSale(HellAmuletID, "Hell Amulet", Side)
var SmirkingDemon = NewItemNotForSale(SmirkingDemonID, "Smirking Demon", Side)

var AnimatedItems = []*Item{
	Burger, Bee, PiggyBank, Clover, GlassOfWine, Monster, Ghost, Butterfly, GoldenGlasses, Hourglass,
}

// Don't mutate it
var ItemsOnSale = append(
	[]*Item{Pacifier, LightBulb, Slippers, IceCream, GrayFish, Ship, Sunglasses, FireExtinguisher, Gamepad, CupCoffee, Strawberry, CubicToyCat, ScubaGlasses, Sword, Voodoo, Mummy,
		Martini, Compass, SunriseCocktail, CosmapolitanGlitch, ToyRobot, DiverHelmet,
		Dinosaur, BacklitPumpkin, Skull, Vendetta, Rocket, Shuriken, Grenade, Rip, VintageCar, Tutanchamun},
	AnimatedItems...
	)

var AnimatedItemsMap = buildMapForAnimatedItems()
func buildMapForAnimatedItems() map[ItemID]*Item  {
	var result = make(map[ItemID]*Item)
	for _, item := range AnimatedItems {
		result[item.ID] = item
	}
	return result
}

func ItemCoinsDayPrice(itemID ItemID) int64 {
	item := AnimatedItemsMap[itemID]
	if item == nil {
		return 0
	}
	return item.PriceList.Day
}

var ItemsOnSaleMap = buildMapOnSaleItems()
func buildMapOnSaleItems() map[ItemID]*Item {
	var result = make(map[ItemID]*Item)
	for _, item := range ItemsOnSale {
		result[item.ID] = item
	}
	return result
}

var ItemsNotForSale = []*Item{Invisible, Torturer, HellAmulet, SmirkingDemon, Crown}
var ItemsNotForSaleMap = buildMapNotForSaleItems()
func buildMapNotForSaleItems() map[ItemID]*Item {
	var result = make(map[ItemID]*Item)
	for _, item := range ItemsNotForSale {
		result[item.ID] = item
	}
	return result
}
var ItemsNotForSaleOrder = buildMapNotForSaleOrder()
func buildMapNotForSaleOrder() map[ItemID]int {
	var result = make(map[ItemID]int)
	for i, item := range ItemsNotForSale {
		result[item.ID] = i
	}
	return result
}

type Item struct {
	ID        ItemID
	SaleType
	PositionType
	Name      string
	PriceList PriceList
}

func NewSideItemForChips(id ItemID, name string, dayPrice int64) *Item {
	return NewItemForChips(id, name, dayPrice, Side)
}

func NewItemForChips(id ItemID, name string, dayPrice int64, pt PositionType,) *Item {
	return &Item{ID: id, Name: name, PositionType: pt, PriceList: buildPriceList(dayPrice), SaleType: ForChips}
}

func NewItemForCoins(id ItemID, name string, dayPrice int64, pt PositionType) *Item {
	return &Item{ID: id, Name: name, PositionType: pt, PriceList: buildPriceList(dayPrice), SaleType: ForCoins}
}

func NewItemNotForSale(id ItemID, name string, pt PositionType) *Item {
	return &Item{ID: id, Name: name, SaleType: NotForSale, PositionType: pt}
}

func (i *Item) GetPrice(tf TimeFrame) (int64, error) {
	switch tf {
	case Day:
		return i.PriceList.Day, nil
	case Week:
		return i.PriceList.Week, nil
	case Month:
		return i.PriceList.Month, nil
	default:
		return 0, errUnknownTimeFrame(tf)
	}
}

type PriceList struct {
	Day int64
	Week int64
	Month int64
}

func buildPriceList(dayPrice int64) PriceList {
	return PriceList{Day: dayPrice, Week: int64(math.Round(float64(dayPrice*7) * 0.85)), Month: int64(math.Round(float64(dayPrice*30) * 0.7))}
}

type SaleType string
const (
	ForCoins   SaleType = "coins"
	ForChips   SaleType = "chips"
	NotForSale SaleType = "notForSale"
)

type PositionType string
const (
	Side     PositionType = "side"
	Top PositionType = "top"
)

func errNoSuchItem(ID ItemID) error {
	return E("item with id %s doesn't exist", ID)
}
