package stats

type CatEngine interface {
	Suggest(query string) []string
}

type QueriesCategoriesDict map[string][]CatStats

type CatStats struct {
	Category string `json:"node_id"`
	Contacts int64  `json:"total_contacts"`
	Searches int64  `json:"total_searches"`
	Score    int64  `json:"total_score"`

	CategoryContacts int64 `json:"node_contacts"`
	CategorySearches int64 `json:"node_searches"`
	CategoryScore    int64 `json:"node_score"`

	CategoryContactRate float64 `json:"node_contact_rate"`
	CategorySearchRate  float64 `json:"node_search_rate"`
	CategoryScoreRate   float64 `json:"node_score_rate"`
}
