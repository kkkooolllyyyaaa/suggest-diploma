package artifact

import (
	"encoding/json"
	"fmt"

	"suggest-runtime/internal/category/tree"
	"suggest-runtime/internal/util/gzippedReader"
)

func ReadNodesFromJson() ([]*tree.NodeInfo, error) {
	jsonFile, err := gzippedReader.NewGzippedJsonReader("data/nodes.json.gz")
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

	nodes := make([]*tree.NodeInfo, 0, 3500)

	for decoder.More() {
		var info tree.NodeInfo
		if err := decoder.Decode(&info); err != nil {
			return nil, err
		}
		nodes = append(nodes, &info)
	}
	_, err = decoder.Token()
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
