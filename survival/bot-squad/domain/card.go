package domain

const MaxFaceRank = 12
const MinFaceRank = 0

var faces = []Face{Deuce, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
var facesDesc = []Face{Ace, King, Queen, Jack, Ten, Nine, Eight, Seven, Six, Five, Four, Three, Deuce}

var faceRanks = make(map[Face]int, len(faces))
var rankFaces = make(map[int]Face, len(faces))

type Card string

func CardFrom(f Face, s Suit) Card {
	return Card(string(f) + string(s))
}

func FaceFromRank(rank int) Face {
	if rank > MaxFaceRank {
		rank = MaxFaceRank
	}
	if rank < MinFaceRank {
		rank = MinFaceRank
	}
	return rankFaces[rank]
}

func (c Card) Face() Face {
	return Face(c[0])
}

// Deuce = 0, Three = 1, ..., Ace = 12
func (c Card) FaceRank() int {
	return faceRanks[c.Face()]
}

func (c Card) Suit() Suit {
	return Suit(c[1])
}

func (c Card) IncFace() Card {
	return CardFrom(c.Face().inc(), c.Suit())
}

func (c Card) String() string {
	return string(c)
}

type Suit rune

const (
	Diamonds Suit = 'd'
	Spades Suit = 's'
	Hearts Suit = 'h'
	Clubs Suit = 'c'
)

type Face rune
const (
	UnknownFace Face = '-'
	Ace         Face = 'A'
	King        Face = 'K'
	Queen       Face = 'Q'
	Jack        Face = 'J'
	Ten         Face = 'T'
	Nine        Face = '9'
	Eight       Face = '8'
	Seven       Face = '7'
	Six         Face = '6'
	Five        Face = '5'
	Four  Face = '4'
	Three Face = '3'
	Deuce Face = '2'
)

func (f Face) Rank() int {
	return faceRanks[f]
}

func (f Face) Incremented() (Face, bool) {
	rank := f.Rank()
	if rank == MaxFaceRank {
		return f, false
	}
	return FaceFromRank(rank +1), true
}

func (f Face) inc() Face {
	res, _ := f.Incremented()
	return res
}

func (f Face) dec() Face {
	res, _ := f.Decremented()
	return res
}

func (f Face) decStraight() Face {
	res, _ := f.DecrementedStraight()
	return res
}

func (f Face) Decremented() (Face, bool) {
	rank := f.Rank()
	if rank == MinFaceRank {
		return f, false
	}
	return FaceFromRank(rank -1), true
}

func (f Face) DecrementedStraight() (Face, bool) {
	rank := f.Rank()
	if rank == MinFaceRank {
		return Ace, true
	}
	return FaceFromRank(rank -1), true
}

