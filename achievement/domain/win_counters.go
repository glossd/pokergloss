package domain

type WinCounter struct {
	SimpleCounter `bson:",inline"`
}

func (w WinCounter) GetType() AchievementType {return SimpleCounterType}
func (w WinCounter) GetName() string {return "Winner"}

func NewWinCounter() *WinCounter {
	return &WinCounter{SimpleCounter: NewSimpleCounter(50, 10, false)}
}

type SitngoWinCounter struct {
	SimpleCounter `bson:",inline"`
}

func NewSitngoWinCounter() *SitngoWinCounter {
	return &SitngoWinCounter{SimpleCounter: NewSimpleCounter(100, 5, false)}
}

func (w SitngoWinCounter) GetType() AchievementType {return SimpleCounterType}
func (w SitngoWinCounter) GetName() string {return "SitNGo Winner"}

type MultiWinCounter struct {
	SimpleCounter `bson:",inline"`
}

func NewMultiWinCounter() *MultiWinCounter {
	return &MultiWinCounter{SimpleCounter: NewSimpleCounter(150, 3, false)}
}

func (w MultiWinCounter) GetType() AchievementType {return SimpleCounterType}
func (w MultiWinCounter) GetName() string {return "Multi SitNGo Winner"}
