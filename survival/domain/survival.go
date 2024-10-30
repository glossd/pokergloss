package domain

import (
	"github.com/glossd/pokergloss/auth/authid"
	"math/rand"
)

type Survival struct {
	UserID      string `bson:"_id"`
	Iden        authid.Identity
	IsAnonymous bool
	IsIdle      bool
	Level       int
	TableID     string
	WoF         *WheelOfFortune
}

type Params struct {
	Anonymous bool
	Idle      bool
}

func New(iden authid.Identity, params Params) *Survival {
	var id = iden
	if params.Anonymous {
		id.Username = "Anonymous"
	}
	return &Survival{
		UserID:      iden.UserId,
		Iden:        id,
		Level:       1,
		IsAnonymous: params.Anonymous,
		IsIdle:      params.Idle,
	}
}

type TableParams struct {
	Name                 string
	BigBlind             int64
	DecisionTimeoutSec   int64
	ThemeID              string
	User                 authid.Identity
	UserStack            int64
	Bots                 []Bot
	LevelIncreaseTimeSec int64
	SurvivalLevel        int64
}

type Bot struct {
	Name    string
	Picture string
	Stack   int64

	Looseness  float64
	Aggression float64
	IsWeak     bool
}

func fogSpirit() Bot {
	return Bot{
		Name:       "Fog Spirit",
		Picture:    "https://storage.googleapis.com/pokerblow/story/spirit-world/fog-spirit.jpg",
		Stack:      30,
		Looseness:  0.8 + randomness(),
		Aggression: 0.2 + randomness(),
		IsWeak:     true,
	}
}

func randomness() float64 {
	return (rand.Float64() - 0.5) * 0.1
}

func elder() Bot {
	return Bot{
		Name:       "Elder",
		Picture:    "https://storage.googleapis.com/pokerblow/story/hell/torturer.jpg",
		Stack:      60,
		Looseness:  0.6 + randomness(),
		Aggression: 0.4 + randomness(),
	}
}

func ancient() Bot {
	return Bot{
		Name:       "Ancient",
		Picture:    "https://storage.googleapis.com/pokerblow/story/hell/sr-torturer.jpg",
		Stack:      80,
		Looseness:  0.5 + randomness(),
		Aggression: 0.5 + randomness(),
	}
}

func demon(loose float64, stack int64) Bot {
	return Bot{
		Name:       "Demon",
		Picture:    "https://storage.googleapis.com/pokerblow/story/hell/demon-v2.jpg",
		Stack:      stack,
		Looseness:  loose + randomness(),
		Aggression: 0.7 + randomness(),
	}
}

func archDemon(loose float64, stack int64) Bot {
	return Bot{
		Name:       "Archdemon",
		Picture:    "https://storage.googleapis.com/pokerblow/story/hell/archdemon.jpg",
		Stack:      stack,
		Looseness:  loose + randomness(),
		Aggression: 0.85 + 2*randomness(),
	}
}

func (s *Survival) NewLevel() {
	s.Level++
}

func (s *Survival) GetTableParams() TableParams {
	hellGamblingParams := TableParams{
		Name:                 "Hell, Forross, Gambling Room",
		BigBlind:             2,
		DecisionTimeoutSec:   10,
		ThemeID:              "hell",
		User:                 s.Iden,
		UserStack:            100,
		LevelIncreaseTimeSec: 3 * 60,
		SurvivalLevel:        int64(s.Level),
	}

	defaultParams := TableParams{
		BigBlind:             2,
		DecisionTimeoutSec:   20,
		User:                 s.Iden,
		LevelIncreaseTimeSec: 3 * 60,
		SurvivalLevel:        int64(s.Level),
	}

	var bots []Bot

	switch s.Level {
	case 1:
		defaultParams.Name = "Paradise, angel Raziel"
		defaultParams.ThemeID = "paradise"
		defaultParams.UserStack = 40
		defaultParams.Bots = []Bot{{
			Name:       "Raziel",
			Picture:    "https://storage.googleapis.com/pokerblow/story/paradise/raziel.png",
			Stack:      20,
			Looseness:  0.9,
			Aggression: 0.1,
			IsWeak:     true,
		}}
		return defaultParams
	case 2:
		defaultParams.Name = "Spirit World, Dosimour town, First Gambling House"
		defaultParams.ThemeID = "spirit-world"
		defaultParams.UserStack = 50
		defaultParams.Bots = []Bot{fogSpirit()}
		return defaultParams
	case 3:
		defaultParams.Name = "Spirit World, Dosimour town, First Gambling House"
		defaultParams.ThemeID = "spirit-world"
		defaultParams.UserStack = 60
		defaultParams.Bots = []Bot{fogSpirit(), {
			Name:       "Hangman Joe",
			Picture:    "https://storage.googleapis.com/pokerblow/story/spirit-world/hangman.jpg",
			Stack:      40,
			Looseness:  0.7,
			Aggression: 0.3,
		}}
		return defaultParams
	case 4:
		defaultParams.Name = "Hell, Forross, Hidden Cave"
		defaultParams.ThemeID = "hell"
		defaultParams.UserStack = 70
		defaultParams.DecisionTimeoutSec = 13
		defaultParams.Bots = []Bot{elder()}
		return defaultParams
	case 5:
		defaultParams.Name = "Hell, Forross, Hidden Cave"
		defaultParams.ThemeID = "hell"
		defaultParams.UserStack = 100
		defaultParams.DecisionTimeoutSec = 13
		defaultParams.Bots = []Bot{elder(), ancient()}
		return defaultParams
	case 6:
		bots = demons(1, 100)
	case 7:
		bots = demons(2, 80)
	case 8:
		bots = demons(3, 60)
	default:
		bots = s.archDemons()
	}

	shuffle(bots)
	hellGamblingParams.Bots = bots
	return hellGamblingParams
}

func (s *Survival) archDemons() []Bot {
	num, stack := s.archLevelData()
	return archDemons(num, stack)
}

func (s *Survival) archLevelData() (int, int64) {
	sum := 0
	cycle := 0
	archLevel := s.Level - 8
	for i := 4; sum < archLevel; i++ {
		if archLevel <= sum+i {
			num := archLevel - sum
			return num, 120 + int64(cycle)*20 - int64(20*(num-1))
		}
		sum += i
		cycle++
	}
	return 1, 100
}

func demons(n int, stack int64) []Bot {
	looses := numberToLooses(n)
	var bots []Bot
	for _, loose := range looses {
		bots = append(bots, demon(loose, stack))
	}
	return bots
}

func archDemons(n int, stack int64) []Bot {
	looses := numberToLooses(n)
	var bots []Bot
	for _, loose := range looses {
		bots = append(bots, archDemon(loose, stack))
	}
	return bots
}

func numberToLooses(n int) []float64 {
	var looses []float64
	switch n {
	case 1:
		looses = []float64{0.5}
	case 2:
		looses = []float64{0.4, 0.6}
	case 3:
		looses = []float64{0.4, 0.5, 0.6}
	case 4:
		looses = []float64{0.4, 0.5, 0.6, 0.7}
	case 5:
		looses = []float64{0.35, 0.45, 0.55, 0.6, 0.65}
	case 6:
		looses = []float64{0.2, 0.3, 0.4, 0.5, 0.6, 0.7}
	case 7:
		looses = []float64{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}
	case 8:
		looses = []float64{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}
	default:
		looses = []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}
	}
	return looses
}

func shuffle(bots []Bot) {
	rand.Shuffle(len(bots), func(i, j int) { bots[i], bots[j] = bots[j], bots[i] })
}
