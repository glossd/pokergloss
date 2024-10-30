package domain

import (
	"github.com/glossd/pokergloss/gomq/mqtable"
	log "github.com/sirupsen/logrus"
	"math"
)

type Hand string

const (
	SF  Hand = "Straight Flush"
	FoK Hand = "Four of a Kind"
	FH  Hand = "Full House"
	F   Hand = "Flush"
	S   Hand = "Straight"
	ToK Hand = "Three of a Kind"
	TP  Hand = "Two Pair"
	P   Hand = "Pair"
	Hi  Hand = "High Card"
)

var handIndex = map[Hand]int{
	TP:  0,
	ToK: 1,
	S:   2,
	F:   3,
	FH:  4,
	FoK: 5,
	SF:  6,
}

type HandsCounter struct {
	Hands map[Hand]*CounterPerHand
	prize Prize
}

type CounterPerHand struct {
	Count int64
	Level int
}

type CounterWithHand struct {
	Hand
	*CounterPerHand
}

func NewHandsCounter() *HandsCounter {
	return &HandsCounter{
		Hands: map[Hand]*CounterPerHand{
			TP:  newCounter(),
			ToK: newCounter(),
			S:   newCounter(),
			F:   newCounter(),
			FH:  newCounter(),
			FoK: newCounter(),
			SF:  newCounter(),
		},
	}
}

func newCounter() *CounterPerHand {
	return &CounterPerHand{Count: 0, Level: 0}
}

func (hc *HandsCounter) Update(winner *mqtable.Winner) error {
	hc.prize = Prize{}

	if winner.Hand == "" {
		log.Tracef("Winner didn't get to flop, userId=%s", winner.UserId)
		return nil
	}
	hand, err := MapToHand(winner.Hand)
	if err != nil {
		log.Errorf("HandsCounter.Update hand mapping: %s", err)
		return err
	}
	if !isHandCounted(hand) {
		return nil
	}
	counter, ok := hc.Hands[hand]
	if !ok {
		log.Errorf("No counter found in hands by hand=%s", hand)
		return err
	}

	counter.Count++
	oldLevel := counter.Level
	counter.Level = computeHandAchievementLevel(counter.Count)
	if oldLevel != counter.Level {
		hc.prize = Prize{
			Chips: ComputeHandLevelPrize(counter.Level, hand),
			Name:  string(hand),
		}
	}
	return nil
}

func (hc *HandsCounter) OrderedCounters() []CounterWithHand {
	var counters = make([]CounterWithHand, len(handIndex))
	for hand, counter := range hc.Hands {
		counters[handIndex[hand]] = CounterWithHand{Hand: hand, CounterPerHand: counter}
	}
	return counters
}

func (hc *HandsCounter) GetPrize() Prize {
	return hc.prize
}

func (hc *HandsCounter) GetPrizeHandCount() int64 {
	return hc.Hands[Hand(hc.GetPrize().Name)].Count
}

func ComputeHandLevelPrize(level int, hand Hand) int64 {
	return getWinningUnit(hand) * GetLevelCount(level)
}

func computeHandAchievementLevel(count int64) int {
	if count == 0 {
		return 0
	}
	level := math.Log10(float64(count))/math.Log10(5) + 1
	return int(level)
}

func GetLevelCount(level int) int64 {
	if level == 0 {
		return 0
	}
	return int64(math.Pow(5, float64(level-1)))
}

func isHandCounted(h Hand) bool {
	if h == Hi || h == P {
		return false
	}
	return true
}

func MapToHand(h string) (Hand, error) {
	switch h {
	case string(Hi), string(P), string(TP), string(ToK), string(S), string(F), string(FH), string(FoK), string(SF):
		return Hand(h), nil
	}
	return "", newEf("no such hand %s", h)
}

func (c *CounterPerHand) LevelCount() int64 {
	return GetLevelCount(c.Level)
}

func (c *CounterPerHand) NextLevelCount() int64 {
	return GetLevelCount(c.Level + 1)
}

func (c *CounterWithHand) NextLevelPrize() int64 {
	return ComputeHandLevelPrize(c.Level+1, c.Hand)
}

func getWinningUnit(h Hand) int64 {
	switch h {
	case TP:
		return 100
	case ToK:
		return 500
	case S:
		return 500
	case F:
		return 750
	case FH:
		return 1000
	case FoK:
		return 15000
	case SF:
		return 75000
	default:
		return 0
	}
}
