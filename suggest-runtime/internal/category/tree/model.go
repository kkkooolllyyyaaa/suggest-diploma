package tree

type NodeInfo struct {
	Id       string `json:"id"`
	ParentId string `json:"parentId"`
	Title    string `json:"title"`
	Icon     struct {
		Uri  string `json:"uri"`
		Size struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"size"`
	} `json:"server_icon"`
	HasChildren bool `json:"-"`
}
