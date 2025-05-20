package artifact

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"

	"suggest-runtime/internal/suggester"
	"suggest-runtime/internal/util/gzippedReader"
)

type QueryInfo struct {
	Searches   int64  `json:"searches"`
	Contacts   int64  `json:"contacts"`
	Query      string `json:"query"`
	RightQuery string `json:"right_query"`
}

const queriesSliceCapacity = 3_500_000

func ReadQueriesFromJson(filename string) ([]*suggester.IndexItem, error) {
	jsonFile, err := gzippedReader.NewGzippedJsonReader(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	startToken, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	if startToken != json.Delim('[') {
		return nil, fmt.Errorf("invalid json file %s", filename)
	}

	queries := make([]*suggester.IndexItem, 0, queriesSliceCapacity)

	for decoder.More() {
		var info QueryInfo
		if err := decoder.Decode(&info); err != nil {
			return nil, err
		}
		queries = append(queries, &suggester.IndexItem{
			Query:           []rune(info.Query),
			NormalizedQuery: []rune(info.RightQuery),
			Score:           float64(suggester.Score(info.Searches, info.Contacts)),
		})
	}
	_, err = decoder.Token()
	if err != nil {
		return nil, err
	}

	slices.SortStableFunc(queries, func(a, b *suggester.IndexItem) int {
		return cmp.Compare(b.Score, a.Score)
	})

	return queries, nil
}
