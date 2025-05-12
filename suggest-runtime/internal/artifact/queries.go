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
	Freq        int64  `json:"freq"`
	ContactFreq int64  `json:"contact_freq"`
	BaseQuery   string `json:"query"`
	RightQuery  string `json:"right_query"`
}

func ReadQueriesFromJson() ([]*suggester.IndexItem, error) {
	jsonFile, err := gzippedReader.NewGzippedJsonReader("data/queries.json.gz")
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
		return nil, fmt.Errorf("invalid json file")
	}

	queries := make([]*suggester.IndexItem, 0, 5_500_000)

	for decoder.More() {
		var info QueryInfo
		if err := decoder.Decode(&info); err != nil {
			return nil, err
		}
		queries = append(queries, &suggester.IndexItem{
			Query:           []rune(info.BaseQuery),
			NormalizedQuery: []rune(info.RightQuery),
			Score:           float64(info.Freq),
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
