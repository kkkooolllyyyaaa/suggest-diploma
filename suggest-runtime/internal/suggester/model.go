package suggester

type IndexItem struct {
	Query           []rune
	NormalizedQuery []rune
	Score           float64
}

type SearchRequest struct {
	Query  string
	UserId string
}

func Score(searches, contacts int64) int64 {
	return searches + 10*contacts
}
