package helpers

type SelectOption struct {
	Id    int64  `json:"id"`
	Label string `json:"label"`
}

type TreeItem struct {
	Id   string
	Pid  string
	Name string
}

type TreeNode struct {
	Id       string     `json:"id"`
	Label    string     `json:"label"`
	Children []TreeNode `json:"children"`
}

func getTreeChildren(items []TreeItem, pid string) []TreeNode {
	result := make([]TreeNode, 0)
	for _, item := range items {
		if item.Pid == pid {
			result = append(result, TreeNode{
				Id:       item.Id,
				Label:    item.Name,
				Children: getTreeChildren(items, item.Id),
			})
		}
	}
	return result
}

func ToTree(items []TreeItem) []TreeNode {
	result := make([]TreeNode, 0)
	for _, item := range items {
		if item.Pid == "0" {
			result = append(result, TreeNode{
				Id:       item.Id,
				Label:    item.Name,
				Children: getTreeChildren(items, item.Id),
			})
		}
	}
	return result
}
