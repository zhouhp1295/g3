package helpers

const EmptyId = 0

type SelectOption struct {
	Id    int64  `json:"id"`
	Label string `json:"label"`
}

type TreeOption struct {
	Id    int64  `json:"id"`
	Pid   int64  `json:"pid"`
	Label string `json:"label"`
}

type TreeNode struct {
	TreeOption
	Children []TreeNode `json:"children"`
}

func getTreeChildren(items []TreeOption, pid int64) []TreeNode {
	result := make([]TreeNode, 0)
	for _, item := range items {
		if item.Pid == pid {
			result = append(result, TreeNode{
				TreeOption: TreeOption{
					Pid:   pid,
					Id:    item.Id,
					Label: item.Label,
				},
				Children: getTreeChildren(items, item.Id),
			})
		}
	}
	return result
}

func ToTree(items []TreeOption) []TreeNode {
	result := make([]TreeNode, 0)
	for _, item := range items {
		if item.Pid == EmptyId {
			result = append(result, TreeNode{
				TreeOption: TreeOption{
					Pid:   EmptyId,
					Id:    item.Id,
					Label: item.Label,
				},
				Children: getTreeChildren(items, item.Id),
			})
		}
	}
	return result
}
