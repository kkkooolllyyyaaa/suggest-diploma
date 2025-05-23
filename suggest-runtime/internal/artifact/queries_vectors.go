package artifact

import (
	"encoding/json"
	"fmt"

	"suggest-runtime/internal/util/gzippedReader"
	"suggest-runtime/internal/vector"
)

const queriesVectorsCapacity = 110_000

func ReadQueriesVectors(filename string) (vector.QueriesVectors, error) {
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

	queriesCategories := make(vector.QueriesVectors, queriesVectorsCapacity)

	for decoder.More() {
		queryToken, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		queryString, ok := queryToken.(string)
		if !ok {
			return nil, fmt.Errorf("invalid json file %s", filename)
		}
		var queryVector vector.QueryVector
		if err := decoder.Decode(&queryVector); err != nil {
			return nil, err
		}
		queriesCategories[queryString] = queryVector
	}
	_, err = decoder.Token()
	if err != nil {
		return nil, err
	}

	return queriesCategories, nil
}

const tokensCapacity = 130_000

func ReadTokensVectors(filename string) (vector.TokensVectors, error) {
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

	tokensVectors := make(vector.TokensVectors, queriesVectorsCapacity)

	for decoder.More() {
		queryToken, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		queryString, ok := queryToken.(string)
		if !ok {
			return nil, fmt.Errorf("invalid json file %s", filename)
		}
		var tokenVector vector.TokenVector
		if err := decoder.Decode(&tokenVector); err != nil {
			return nil, err
		}
		tokensVectors[queryString] = tokenVector
	}
	_, err = decoder.Token()
	if err != nil {
		return nil, err
	}

	return tokensVectors, nil
}
