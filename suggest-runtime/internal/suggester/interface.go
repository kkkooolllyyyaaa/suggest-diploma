package suggester

const SuggestLimit = 8

type SuggestId int32

type Suggester interface {
	Build(collection []*IndexItem)
	Suggest(request SearchRequest) []*IndexItem
}
