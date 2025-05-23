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
