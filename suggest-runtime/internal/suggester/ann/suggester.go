package ann

import (
	"cmp"
	"slices"

	"suggest-runtime/internal/suggester"
	"suggest-runtime/internal/vector"
)

const maxItems = 4

type annSuggester struct {
	annIndex      vector.AnnIndex
	queriesScores map[string]float64
}

func NewAnnSuggester(
	annIndex vector.AnnIndex,
) suggester.Suggester {
	return &annSuggester{
		annIndex: annIndex,
	}
}

func (s *annSuggester) Build(collection []*suggester.IndexItem) {
	queriesScores := make(map[string]float64, len(collection))
	for _, q := range collection {
		queriesScores[string(q.Query)] = q.Score
	}
	s.queriesScores = queriesScores
}

func (s *annSuggester) Suggest(request suggester.SearchRequest) []*suggester.IndexItem {
	userQueryNearestQueries := s.annIndex.NearestQueries(request.Query)
	slices.SortFunc(userQueryNearestQueries, func(a, b string) int {
		return cmp.Compare(s.queriesScores[a], s.queriesScores[b])
	})

	indexItems := make([]*suggester.IndexItem, 0, maxItems)
	for _, q := range userQueryNearestQueries {
		runes := []rune(q)
		indexItems = append(indexItems, &suggester.IndexItem{
			Query:           runes,
			NormalizedQuery: runes,
			Score:           s.queriesScores[q],
		})
		if len(indexItems) >= maxItems {
			break
		}
	}
	return indexItems
}
