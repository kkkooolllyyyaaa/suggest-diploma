package artifact

import (
	"encoding/json"
	"fmt"

	"suggest-runtime/internal/category/cat_stats"
	"suggest-runtime/internal/util/gzippedReader"
)

const dictSize = 5_000_000

func ReadQueryNodeDictRaw(filename string) (cat_stats.QueriesCatDict, error) {
	jsonFile, err := gzippedReader.NewGzippedJsonReader("data/queries_categories_propagated.json.gz")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	startToken, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	if startToken != json.Delim('{') { // must be dict
		return nil, fmt.Errorf("invalid json file %s", filename)
	}

	sourceDict := make(cat_stats.QueriesCatDict, dictSize)

	for decoder.More() {
		queryToken, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		queryString, ok := queryToken.(string)
		if !ok {
			return nil, fmt.Errorf("invalid json file %s", filename)
		}
		var stats []cat_stats.CatStats
		if err := decoder.Decode(&stats); err != nil {
			return nil, err
		}
		sourceDict[queryString] = stats
	}
	_, err = decoder.Token()
	if err != nil {
		return nil, err
	}

	return sourceDict, nil
}
