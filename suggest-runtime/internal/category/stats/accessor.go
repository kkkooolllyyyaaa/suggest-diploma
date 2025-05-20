package stats

type CatStatsAccessor interface {
	QueryFreq(stats CatStats) int64
	CategoryRate(stats CatStats) float64
}

// contact Accessor
type categoryContactsAccessor struct{}

func NewCategoryContactsAccessor() CatStatsAccessor {
	return &categoryContactsAccessor{}
}

func (n categoryContactsAccessor) QueryFreq(stats CatStats) int64 {
	return stats.Contacts
}

func (n categoryContactsAccessor) CategoryRate(stats CatStats) float64 {
	return stats.CategoryContactRate
}

// score Accessor
type categoryScoreAccessor struct{}

func NewCategoryScoreAccessor() CatStatsAccessor {
	return &categoryScoreAccessor{}
}

func (n categoryScoreAccessor) QueryFreq(stats CatStats) int64 {
	return stats.Score
}

func (n categoryScoreAccessor) CategoryRate(stats CatStats) float64 {
	return stats.CategoryScoreRate
}

// search Accessor
type categorySearchAccessor struct{}

func NewCategorySearchAccessor() CatStatsAccessor {
	return &categorySearchAccessor{}
}

func (n categorySearchAccessor) QueryFreq(stats CatStats) int64 {
	return stats.Searches
}

func (n categorySearchAccessor) CategoryRate(stats CatStats) float64 {
	return stats.CategorySearchRate
}
