package vector

type QueryVector struct {
	Index  int       `json:"index"`
	Vector []float32 `json:"-"`
}
type QueriesVectors map[string]QueryVector

type TokenVector []float32
type TokensVectors map[string]TokenVector
