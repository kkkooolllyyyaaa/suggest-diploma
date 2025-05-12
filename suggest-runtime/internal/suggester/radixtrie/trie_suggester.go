package radixtrie

import (
	"slices"
	"strings"

	"suggest-runtime/internal/suggester"
)

type Suggester struct {
	suggests []*suggester.IndexItem
	trie     *Trie
}

func NewTrieSuggester() suggester.Suggester {
	return &Suggester{
		trie: NewTrie(),
	}
}

func (s *Suggester) Build(collection []*suggester.IndexItem) {
	for _, item := range collection {
		s.trie.Put(item)
	}

	s.suggests = collection
}

func (s *Suggester) Suggest(request suggester.SearchRequest) []*suggester.IndexItem {
	if len(request.Query) == 0 {
		return nil
	}

	tokens := strings.Fields(request.Query)
	lookupsResults := make([]*Node, 0, len(tokens))

	for _, token := range tokens {
		result := s.trie.Get([]rune(token))
		if result != nil {
			lookupsResults = append(lookupsResults, result)
		}
	}

	intersectedIndexes := s.intersectIndexes(lookupsResults)
	actualResult := make([]*suggester.IndexItem, 0, len(intersectedIndexes))

	for _, index := range intersectedIndexes {
		actualResult = append(actualResult, s.suggests[index])
	}

	slices.SortFunc(actualResult, func(a, b *suggester.IndexItem) int {
		if (b.Score - a.Score) >= 0.0 {
			return 1
		}
		return -1
	})

	if len(actualResult) > suggester.SuggestLimit {
		return actualResult[:suggester.SuggestLimit]
	}
	return actualResult
}

func (s *Suggester) intersectIndexes(results []*Node) []int {
	if len(results) == 0 {
		return nil
	}

	maps := make([]map[int]struct{}, 0, len(results))
	for i := 0; i < len(results); i++ {
		maps = append(maps, make(map[int]struct{}, len(results[i].Index)))
	}

	for i, node := range results {
		for _, id := range node.Index {
			maps[i][id] = struct{}{}
		}
	}

	intersected := make([]int, 0, len(maps[0])/2+1)
	for toCheck := range maps[0] {
		containsAll := true
		for i := 1; i < len(maps); i++ {
			if _, ok := maps[i][toCheck]; !ok {
				containsAll = false
				break
			}
		}

		if containsAll {
			intersected = append(intersected, toCheck)
		}
	}
	return intersected
}
