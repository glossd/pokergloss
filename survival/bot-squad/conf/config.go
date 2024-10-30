package conf

import (
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	UserID  string `mapstructure:"user_id"`
	TableID string `mapstructure:"table_id"`
	Squad
	TableEvents `mapstructure:"table_events"`
	Table
	UserPosition int `mapstructure:"user_position"`
	Token        string
	Protocol     string
}

type Squad struct {
	Positions []string
	// Deprecated.
	WeakPositions []string
	Looseness     []string
	Aggression    []string
}

func (s Squad) GetPositions() []int {
	return stringsToInts(s.Positions, "PB_SQUAD_POSITIONS")
}

func (s Squad) GetWeakPositions() []int {
	return stringsToInts(s.WeakPositions, "PB_SQUAD_WEAK_POSITIONS")
}

func (s Squad) GetLooseness() []float64 {
	looseness := s.Looseness
	return stringsToFloats(looseness, "PB_SQUAD_LOOSENESS")
}

func (s Squad) GetAggression() []float64 {
	return stringsToFloats(s.Aggression, "PB_SQUAD_AGGRESSION")
}

func (s Squad) GetWeakPositionsSet() map[int]struct{} {
	res := make(map[int]struct{})
	for _, pos := range s.GetWeakPositions() {
		res[pos] = struct{}{}
	}
	return res
}

func stringsToInts(strSlice []string, description string) []int {
	var positions []int
	for _, position := range strSlice {
		trimed := strings.TrimSpace(position)
		for _, s := range strings.Split(trimed, " ") {
			pos, err := strconv.Atoi(s)
			if err != nil {
				panic(fmt.Sprintf("Failed to parse %s %v", description, strSlice))
			}
			positions = append(positions, pos)
		}
	}
	return positions
}

func stringsToFloats(strSlice []string, description string) []float64 {
	var positions []float64
	for _, position := range strSlice {
		trimed := strings.TrimSpace(position)
		for _, s := range strings.Split(trimed, " ") {
			pos, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic(fmt.Sprintf("Failed to parse %s %v", description, strSlice))
			}
			positions = append(positions, pos)
		}
	}
	return positions
}

type Table struct {
	Scheme string
	Host   string
	Port   int
}

type TableEvents struct {
	Host string
	Port int
}
