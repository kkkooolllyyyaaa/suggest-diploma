package tree

const RootNodeId = "1"

type CategoryTree interface {
	GetChildren(nodeId string) []*NodeInfo
}
