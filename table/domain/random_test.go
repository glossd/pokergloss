package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkRealAlgo_RandomThreeCards(b *testing.B) {
	table, err := NewTable(defaultTableParams)
	assert.Nil(b, err)
	for i := 0; i < b.N; i++ {
		f, s, th := Algo.RandomAvailableThreeCards(table)
		assert.NotEqualValues(b, f, s)
		assert.NotEqualValues(b, f, th)
		assert.NotEqualValues(b, s, th)
	}
}

func BenchmarkRealAlgo_RandomTwoCards(b *testing.B) {
	table, err := NewTable(defaultTableParams)
	assert.Nil(b, err)
	for i := 0; i < b.N; i++ {
		f, s:= Algo.RandomAvailableTwoCards(table)
		assert.NotEqualValues(b, f, s)
	}
}

func Algo_2P_MockSecondPlayerLoses() *MockAlgo {
	return NewMockCards("As", "Ks", "2d", "7h", "Qs", "Js", "Ts", "8h", "4c")
}

func Algo_2P_MockFirstPlayerLoses() *MockAlgo {
	return NewMockCards("2d", "7h", "As", "Ks", "Qs", "Js", "Ts", "8h", "4c")
}


func Algo_2P_MockDraw() *MockAlgo {
	return NewMockCards("2d", "7h", "2h", "7d", "Qs", "Qd", "Qc", "Qh", "4c")
}
func AlgoMock_3P_Second_Third_First(t *testing.T) *MockAlgo {
	mock, err :=  NewMockAlgo(CardsStr(/*first game doesn't matter*/"2c", "4s", "6c", "4h",
		/*new game started*/ "2d", "7h", "As", "Ks", "9s", "8s", "Qs", "Js", "Ts", "8h", "4c"))
	assert.Nil(t, err)
	return mock
}


func AlgoMock_3POrdered_Second_Third_First(t *testing.T) *MockAlgo {
	mock, err :=  NewMockAlgoMultiGame(
		CardsStr("2c", "2s", "2h", "2d"),
		CardsStr("2c", "2s", "2h", "2d"),
		CardsStr(
			"2d", "7h",
			"As", "Ks",
			"9s", "8s",
			"Qs", "Js", "Ts", "8h", "4c"))
	assert.Nil(t, err)
	return mock
}

func AlgoMock_4P_FirstAndSecond_Third_Fourth(t *testing.T) *MockAlgo {
	mock, err :=  NewMockAlgoMultiGame(
		CardsStr("2c", "2d", "2h", "2s"),
		CardsStr("2c", "2d", "2h", "2s"),
		CardsStr("2c", "2d", "2h", "2s"),
		CardsStr(
		"Ah", "Kc",
		"Ac", "Ks",
		"8s", "5h",
		"5c", "6h",
		"As", "Ad", "Kh", "8d", "4c"))
	assert.Nil(t, err)
	return mock
}

func AlgoMock_3P_ThirdAndFirst_Second(t *testing.T) *MockAlgo {
	mock, err :=  NewMockAlgo(CardsStr(/*first game doesn't matter*/"2c", "2d", "2h", "2s",
		/*new game started*/ "Ad", "7c", "4h", "3s", "As", "7d",
		/*community cards*/ "4s", "6h", "8c", "Ah", "Ac"))
	assert.Nil(t, err)
	return mock
}
