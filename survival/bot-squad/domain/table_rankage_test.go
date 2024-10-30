package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_straightRankage(t *testing.T) {
	assert.EqualValues(t, 1, straightRankage([]string{"As", "Qd"}, "Kc", "Jh", "Td", "7h", "2c"))
	assert.EqualValues(t, 0.45, straightRankage([]string{"As", "Qd"}, "Kc", "5h", "4d", "3h", "2c"))
	assert.EqualValues(t, 1, straightRankage([]string{"7s", "6c"},  "9c", "8c", "5s", "2d"))
}

func Test_twoPairRankage(t *testing.T) {
	assert.EqualValues(t, 1, twoPairRankage([]string{"Ks", "Qs"}, "Kd", "Qd", "6s", "5s"))
	assert.EqualValues(t, 5.0/6, twoPairRankage([]string{"Ks", "6c"}, "Kd", "Qd", "6s", "5s"))
	assert.EqualValues(t, 4.0/6, twoPairRankage([]string{"Ks", "5c"}, "Kd", "Qd", "6s", "5s"))
	assert.EqualValues(t, 3.0/6, twoPairRankage([]string{"Qc", "6c"}, "Kd", "Qd", "6s", "5s"))

	assert.EqualValues(t, 1.0/6, twoPairRankage([]string{"6c", "5c"}, "Kd", "Qd", "6s", "5s"))
}

func Test_pairRankage(t *testing.T) {
	assert.EqualValues(t, 1, pairRankage([]string{"9s", "2c"}, "9s", "8c", "3h"))
	assert.EqualValues(t, 2.0/3, pairRankage([]string{"8s", "2c"}, "9s", "8c", "3h"))

	assert.EqualValues(t, 1, pairRankage([]string{"Js", "Jc"}, "9s", "8c", "3h"))
	assert.EqualValues(t, 2.0/4, pairRankage([]string{"7s", "7c"}, "9s", "8c", "3h"))
	assert.EqualValues(t, 3.0/4, pairRankage([]string{"8s", "8c"}, "9s", "7c", "3h"))
	assert.EqualValues(t, 0.2, pairRankage([]string{"2s", "2c"}, "9s", "8c", "3h"))
}

func Test_highCardRankage(t *testing.T) {
	assert.EqualValues(t, 1, highCardRankage("As", "Qs", "Jc", "7h", "6d", "2h"))
	assert.EqualValues(t, 7.0/8, highCardRankage("Ks", "Qs", "Jc", "7h", "6d", "2h"))
	assert.EqualValues(t, 2.0/8, highCardRankage("4s", "Qs", "Jc", "8h", "7d", "3h"))
	assert.EqualValues(t, 1.0/8, highCardRankage("2s", "Qs", "Jc", "8h", "7d", "3h"))


	assert.EqualValues(t, 1.0, highCardRankage("As", "Qs", "Jc", "8h", "7d"))
	assert.EqualValues(t, 1.0/9, highCardRankage("2s", "Qs", "Jc", "8h", "7d"))
}

func TestTable_FindCardsOfStraight(t *testing.T) {
	assert.Len(t, findStraight("Ks", "Qd", "Jh"), 3)
	assert.Len(t, findStraight("Ks", "Qd", "Th"), 3)
	assert.Len(t, findStraight("Ks", "Qd", "9h"), 3)
	assert.Len(t, findStraight("Ks", "Td", "9h"), 3)


	assert.Len(t, findStraight("Ks", "Qd", "Jh", "Tc"), 4)
	assert.Len(t, findStraight("Ks", "Qd", "Th", "9c"), 4)

	straight3OutOf4 := findStraight("Ks", "Qd", "9h", "2c")
	assert.Len(t, straight3OutOf4, 3)
	assert.EqualValues(t, straight3OutOf4[2], "9h")

	assert.Len(t, findStraight("Ks", "Qd", "9h", "8c"), 4)
	assert.Len(t, findStraight("Ks", "Jd", "9h", "7c"), 4)


	assert.Len(t, findStraight("Ks", "Jd", "9h", "7c", "5s"), 5)
	assert.Len(t, findStraight("Ks", "9s", "8c", "5d", "4s"), 4)
}

func findStraight(cards ...string) Cards {
	return (&Table{CommCards: cardsStr(cards...)}).FindCardsOfStraight()
}

func straightRankage(pocket []string, comm ...string) float64 {
	return (&Table{CommCards: cardsStr(comm...)}).straightRankage(HoleCards(cardsStr(pocket...)))
}

func twoPairRankage(pocket []string, comm ...string,) float64 {
	return (&Table{CommCards: cardsStr(comm...)}).twoPairRankage(HoleCards(cardsStr(pocket...)))
}

func pairRankage(pocket []string, comm ...string,) float64 {
	return (&Table{CommCards: cardsStr(comm...)}).pairRankage(HoleCards(cardsStr(pocket...)))
}

func highCardRankage(pocket string, comm ...string, ) float64 {
	return (&Table{CommCards: cardsStr(comm...)}).highCardRankage(Card(pocket).Face())
}

func Test_factorial(t *testing.T) {
	assert.EqualValues(t, 1, F(0))
	assert.EqualValues(t, 1, F(1))
	assert.EqualValues(t, 2, F(2))
	assert.EqualValues(t, 6, F(3))
	assert.EqualValues(t, 24, F(4))
	assert.EqualValues(t, 120, F(5))
}

func cardsStr(cards ...string) (c Cards) {
	for _, card := range cards {
		c = append(c, Card(card))
	}
	return
}
