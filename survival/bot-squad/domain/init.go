package domain

func init() {
	for i, face := range faces {
		faceRanks[face] = i
	}
	for i, face := range faces {
		rankFaces[i]=face
	}
	for i, s := range patternRanks {
		patternRankMap[s] = i
	}
}
