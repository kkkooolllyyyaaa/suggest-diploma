package artifact

import (
	"encoding/json"
	"fmt"

	"suggest-runtime/internal/category/tree"
	"suggest-runtime/internal/util/gzippedReader"
)

const nodesSliceCapacity = 3500

func ReadNodesFromJson(filename string) ([]*tree.NodeInfo, error) {
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

	nodes := make([]*tree.NodeInfo, 0, nodesSliceCapacity)

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

	fmt.Println(len(nodes))

	return nodes, nil
}
