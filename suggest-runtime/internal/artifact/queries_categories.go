package artifact

import (
	"encoding/json"
	"fmt"

	"suggest-runtime/internal/category/stats"
	"suggest-runtime/internal/util/gzippedReader"
)

const dictCapacity = 4_000_000

func ReadQueriesCategories(filename string) (stats.QueriesCategoriesDict, error) {
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
	if startToken != json.Delim('{') {
		return nil, fmt.Errorf("invalid json file %s", filename)
	}

	queriesCategories := make(stats.QueriesCategoriesDict, dictCapacity)

	for decoder.More() {
		queryToken, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		queryString, ok := queryToken.(string)
		if !ok {
			return nil, fmt.Errorf("invalid json file %s", filename)
		}
		var stats []stats.CatStats
		if err := decoder.Decode(&stats); err != nil {
			return nil, err
		}
		queriesCategories[queryString] = stats
	}
	_, err = decoder.Token()
	if err != nil {
		return nil, err
	}

	return queriesCategories, nil
}
