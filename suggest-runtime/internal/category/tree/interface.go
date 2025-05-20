package tree

const RootCategoryId = "1"

type CategoryTree interface {
	Children(nodeId string) []*NodeInfo
	Parents(nodeId string) []string
	Parent(id string) *string
	Depth(nodeId string) int
}
