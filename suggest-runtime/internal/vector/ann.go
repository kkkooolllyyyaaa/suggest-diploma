package vector

import (
	"fmt"
	"strings"

	"suggest-runtime/internal/config"

	"github.com/mariotoffia/goannoy/builder"
	"github.com/mariotoffia/goannoy/interfaces"
)

type AnnIndex struct {
	annoyIndex    interfaces.AnnoyIndex[float32, uint32]
	tokensVectors TokensVectors
	indexQueries  map[int]string
	dim           int
	count         int
	minScore      float32
}

func NewIndex(
	cfg *config.Config,
	queriesVectors QueriesVectors,
	tokensVectors TokensVectors,
) AnnIndex {
	idx := builder.Index[float32, uint32]().
		AngularDistance(cfg.Vector.Dimension).
		UseMultiWorkerPolicy().
		MmapIndexAllocator().
		IndexNumHint(100_000).
		Build()

	idx.Load(cfg.Artifact.AnnoyIndex)

	indexQueries := make(map[int]string, len(queriesVectors))
	for q, v := range queriesVectors {
		indexQueries[v.Index] = q
	}

	return AnnIndex{
		annoyIndex:    idx,
		tokensVectors: tokensVectors,
		indexQueries:  indexQueries,
		dim:           cfg.Vector.Dimension,
		count:         cfg.Vector.Count,
		minScore:      cfg.Vector.MinDist,
	}
}

func (idx AnnIndex) meanVector(vectors []TokenVector) TokenVector {
	if len(vectors) == 0 {
		return nil
	} else if len(vectors) == 1 {
		return vectors[0]
	}

	vector := make(TokenVector, 0, idx.dim)
	for i := 0; i < idx.dim; i++ {
		vector = append(vector, 0.0)
	}

	for _, v := range vectors {
		for i, f := range v {
			vector[i] += f
		}
	}

	for i := 0; i < idx.dim; i++ {
		vector[i] /= float32(len(vectors))
	}
	return vector
}

func (idx AnnIndex) NearestQueries(query string) []string {
	tokens := strings.Fields(query)
	tokensVectors := make([]TokenVector, 0, len(tokens))
	for _, token := range tokens {
		if vec, ok := idx.tokensVectors[token]; !ok {
			return nil
		} else {
			tokensVectors = append(tokensVectors, vec)
		}
	}

	if len(tokensVectors) == 0 {
		return nil
	}

	meanVector := idx.meanVector(tokensVectors)
	ctx := idx.annoyIndex.CreateContext()

	indexes, scores := idx.annoyIndex.GetNnsByVector(meanVector, idx.count, -1, ctx)
	fmt.Println(len(indexes))
	result := make([]string, 0, min(len(indexes), idx.count))
	for i, index := range indexes {
		if scores[i] < idx.minScore {
			break
		}
		fmt.Println("Has Score!")

		if q, ok := idx.indexQueries[int(index)]; ok {
			result = append(result, q)
		}
	}
	return result
}
