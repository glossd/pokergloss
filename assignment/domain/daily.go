package domain

import (
	"math/rand"
	"time"
)

const layoutISO = "2006-01-02"
type DailyID string
func NewDailyID(t time.Time) DailyID {
	return DailyID(t.Format(layoutISO))
}

type Daily struct {
	Day DailyID `bson:"_id"`
	Assignments []*Assignment
	// seconds
	CreatedAt int64
}

func NewDaily(t time.Time) *Daily {
	first := NewAssignment(getRand(handUnits))

	var second *Assignment
	secondRand := rand.Float64()
	if secondRand < 0.2 {
		second = NewAssignment(getRand(pairOfUnits))
	} else if secondRand < 0.4  {
		second = NewAssignment(getRand(sharePotUnits))
	} else if secondRand < 0.6 {
		second = NewAssignment(getRand(faceUnits))
	} else {
		second = NewAssignment(getRand(winUnits))
	}

	var third *Assignment
	thirdRand := rand.Float64()
	if thirdRand < 0.33 {
		third = NewAssignment(getRand(bustUnits))
	} else if thirdRand < 0.66 {
		third = NewAssignment(getRand(defeatUnits))
	} else {
		third = NewAssignment(getRand(scareAwayUnits))
	}

	return &Daily{
		Day: NewDailyID(t),
		CreatedAt: time.Now().Unix(),
		Assignments: []*Assignment{first, second, third},
	}
}

func (d *Daily) ContainsLosing() bool {
	for _, assignment := range d.Assignments {
		if assignment.IsToLose() {
			return true
		}
	}
	return false
}

var handUnits []*AssignmentUnit
var pairOfUnits []*AssignmentUnit
var faceUnits []*AssignmentUnit

var bustUnits []*AssignmentUnit
var defeatUnits []*AssignmentUnit
var scareAwayUnits []*AssignmentUnit
var sharePotUnits []*AssignmentUnit

var winUnits []*AssignmentUnit

func init() {
	handUnits = append(joinWithHands(WinWithHand), joinWithHands(LoseWithHand)...)
	pairOfUnits = append(joinFaces(LoseWithPairOf), joinFaces(WinWithPairOf)...)
	faceUnits = append(faceUnits,
		NewWithFace(WinWithFace, 'A'),
		NewWithFace(WinWithFace, 'K'),
		NewWithFace(LoseWithFace, 'A'),
		NewWithFace(LoseWithFace, 'K'),
	)

	bustUnits = append(bustUnits,
		NewWithNumber(BustPlayers, 1),
		NewWithNumber(BustPlayers, 2),
		NewWithNumber(BustPlayers, 3),
		)

	defeatUnits = append(defeatUnits,
		NewWithNumber(DefeatPlayers, 2),
		NewWithNumber(DefeatPlayers, 3),
		NewWithNumber(DefeatPlayers, 4),
		)

	scareAwayUnits = append(scareAwayUnits, joinWithRound(ScareAway, 1)...)
	scareAwayUnits = append(scareAwayUnits, joinWithRound(ScareAway, 2)...)
	scareAwayUnits = append(scareAwayUnits, joinWithRound(ScareAway, 3)...)
	scareAwayUnits = append(scareAwayUnits, joinWithRound(ScareAway, 4)...)

	sharePotUnits = append(sharePotUnits, NewWithNumber(SharePot, 1), NewWithNumber(SharePot, 2))

	winUnits = append(winUnits, New(Win), New(WinLive), New(WinSitNGo), New(WinMultiSitNGo))
}

func joinFaces(t AssignmentType) (a []*AssignmentUnit) {
	for _, f := range []rune{'A', 'K', 'Q', 'J'} {
		a = append(a, NewWithFace(t, f))
	}
	return
}

func getRand(a []*AssignmentUnit) AssignmentUnit {
	return *a[rand.Intn(len(a))]
}
