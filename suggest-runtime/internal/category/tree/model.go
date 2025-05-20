package tree

type NodeInfo struct {
	Id          string `json:"id"`
	ParentId    string `json:"parentId"`
	Title       string `json:"title"`
	HasChildren bool   `json:"-"`
}
