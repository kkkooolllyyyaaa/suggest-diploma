package tree

type categoryTree struct {
	nodesById map[string]*NodeInfo
	children  map[string][]string
	depths    map[string]int
}

func NewCategoryTree(infos []*NodeInfo) CategoryTree {
	categoryTree := categoryTree{}
	categoryTree.initTree(infos)
	return &categoryTree
}

func (t *categoryTree) GetChildren(id string) []*NodeInfo {
	children, ok := t.children[id]
	if !ok {
		return nil
	}

	result := make([]*NodeInfo, 0, len(children))
	for _, nodeId := range children {
		got := t.nodesById[nodeId]
		_, got.HasChildren = t.children[nodeId]
		result = append(result, got)
	}
	return result
}

func (t *categoryTree) initTree(nodes []*NodeInfo) {
	children := make(map[string][]string, len(nodes))
	nodesById := make(map[string]*NodeInfo, len(nodes))

	for _, nodeInfo := range nodes {
		nodesById[nodeInfo.Id] = nodeInfo

		parentId := nodeInfo.ParentId
		if parentId == "" {
			continue
		}

		children[parentId] = append(children[parentId], nodeInfo.Id)
	}

	t.nodesById = nodesById
	t.children = children
}
