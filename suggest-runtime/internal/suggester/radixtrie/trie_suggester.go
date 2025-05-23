package radixtrie

import (
	"cmp"
	"slices"
	"strings"

	"suggest-runtime/internal/suggester"
)

const maxIntersected = 100

type trieSuggester struct {
	suggests []*suggester.IndexItem
	trie     *Trie
}

func NewTrieSuggester() suggester.Suggester {
	return &trieSuggester{
		trie: NewTrie(),
	}
}

func (s *trieSuggester) Build(collection []*suggester.IndexItem) {
	for _, item := range collection {
		s.trie.Put(item)
	}

	s.suggests = collection
}

func (s *trieSuggester) Suggest(request suggester.SearchRequest) []*suggester.IndexItem {
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

	uniqueQueriesIdx := map[string]int{}
	for _, index := range intersectedIndexes {
		indexItem := s.suggests[index]
		normQuery := string(indexItem.NormalizedQuery)
		if idx, ok := uniqueQueriesIdx[normQuery]; ok {
			if actualResult[idx].Score < indexItem.Score {
				actualResult[idx] = indexItem
			}
			continue
		}

		uniqueQueriesIdx[normQuery] = len(actualResult)
		actualResult = append(actualResult, s.suggests[index])
	}

	slices.SortFunc(actualResult, func(a, b *suggester.IndexItem) int {
		return cmp.Compare(b.Score, a.Score)
	})

	if len(actualResult) > suggester.SuggestLimit {
		return actualResult[:suggester.SuggestLimit]
	}
	return actualResult
}

func (s *trieSuggester) intersectIndexes(results []*Node) []int {
	if len(results) == 0 {
		return nil
	} else if len(results) == 1 {
		toReturn := results[0].Index
		if len(toReturn) > maxIntersected {
			toReturn = toReturn[:maxIntersected]
		}
		return toReturn
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
			if len(intersected) == maxIntersected {
				break
			}
		}
	}
	return intersected
}
