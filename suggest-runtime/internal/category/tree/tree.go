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

func (t *categoryTree) Children(id string) []*NodeInfo {
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

func (t *categoryTree) Parents(id string) []string {
	var parents []string
	node, ok := t.nodesById[id]
	if !ok {
		return nil
	}

	parent := node.ParentId
	for len(parent) > 0 {
		parents = append(parents, parent)
		parent = t.nodesById[parent].ParentId
	}
	return parents
}

func (t *categoryTree) Title(id string) string {
	node, ok := t.nodesById[id]
	if !ok {
		return ""
	}
	return node.Title
}

func (t *categoryTree) Parent(id string) *string {
	node, ok := t.nodesById[id]
	if !ok {
		return nil
	}

	if len(node.ParentId) > 0 {
		return &node.ParentId
	}
	return nil
}

func (t *categoryTree) Depth(id string) int {
	return len(t.Parents(id))
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
