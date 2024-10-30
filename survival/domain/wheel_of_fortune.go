package domain

import (
	"math/rand"
)

type Slot struct {
	Chips int64
	Item  *Item
}

type Item struct {
	ItemID string
	Days int64
}

type WheelOfFortune struct {
	Slots []Slot
	WonSlotIdx int
}

func (wof *WheelOfFortune) WonSlot() Slot {
	return wof.Slots[wof.WonSlotIdx]
}

func (s *Survival) CreateWheelOfFortune() {
	if s.Level == 1 {
		return
	}
	s.WoF = createWoF(int64(s.Level-1))
}

func (s *Survival) GetWheelOfFortune() *WheelOfFortune {
	return s.WoF
}

func createWoF(passedLevel int64) *WheelOfFortune {
	var chips int64
	switch {
	case passedLevel < 6:
		chips = passedLevel * 1000
	case passedLevel < 9:
		chips = 5000 + (passedLevel-5)*2000
	default:
		chips = 11000 + (passedLevel-8)*3000
	}

	a := chipSlots(chips)
	itemSlot := getItemSlot(passedLevel)
	var wonSlot Slot
	if itemSlot != nil {
		a = append(a, Slot{Item: itemSlot})
		wonSlot = a[spin(1.0)]
	} else {
		wonSlot = a[spin(0.9)]
		// add repeats
		a = append(a, Slot{Chips: chips})
	}


	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	var wonSlotIdx int
	for i, slot := range a {
		if wonSlot.Chips == slot.Chips {
			wonSlotIdx = i
		}
	}

	return &WheelOfFortune{Slots: a, WonSlotIdx: wonSlotIdx}
}

func getItemSlot(passedLevel int64) *Item {
	if passedLevel < 4 {
		return nil
	}
	switch {
	case passedLevel <= 5:
		return &Item{ItemID: "torturer", Days: passedLevel - 3}
	case passedLevel <= 8:
		return &Item{ItemID: "hellAmulet", Days: passedLevel - 5}
	default:
		return &Item{ItemID: "smirkingDemon", Days: passedLevel - 8}
	}
}

func chipSlots(chips int64) []Slot {
	return []Slot{
		{Chips: chips/4}, {Chips: chips/2}, {Chips: chips/4*3}, {Chips: chips/10*8}, {Chips: chips/10*9},
		{Chips:chips},
		{Chips: chips/2*3}, {Chips: chips*2}, {Chips: chips*3}, {Chips: 5*chips}, {Chips: 10*chips},
	}
}

var slotChances = []float64{0.01, 0.09, 0.1, 0.12, 0.18, 0.3, 0.05, 0.025, 0.015, 0.009, 0.001, 0.1}
var slotSummedChances []float64
func init() {
	slotSummedChances = make([]float64, 0, len(slotChances))
	var sum float64
	for _, chance := range slotChances {
		sum += chance
		slotSummedChances = append(slotSummedChances, sum)
	}
}

func spin(coef float64) int {
	r := rand.Float64()*coef
	for i, chance := range slotSummedChances {
		if r < chance {
			return i
		}
	}

	return 5
}