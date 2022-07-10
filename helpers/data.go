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

func toParentLevels(items []TreeOption, itemId int64) []TreeOption {
	result := make([]TreeOption, 0)
	for _, item := range items {
		if item.Id == itemId {
			result = append(result, item)
			parents := toParentLevels(items, item.Pid)
			if len(parents) > 0 {
				result = append(result, parents...)
			}
			break
		}
	}
	return result
}

func ToParentLevels(items []TreeOption, itemId int64) []TreeOption {
	result := make([]TreeOption, 0)
	for _, item := range items {
		if item.Id == itemId {
			result = append(result, item)
			parents := toParentLevels(items, item.Pid)
			if len(parents) > 0 {
				result = append(result, parents...)
			}
			break
		}
	}
	Reverse[TreeOption](result)
	return result
}

func findChildren(items []TreeOption, itemId int64) []TreeOption {
	result := make([]TreeOption, 0)
	for _, item := range items {
		if item.Pid == itemId {
			result = append(result, item)
			children := findChildren(items, item.Id)
			if len(children) > 0 {
				children = append(result, children...)
			}
		}
	}
	return result
}

func FindChildren(items []TreeOption, itemId int64) []TreeOption {
	result := make([]TreeOption, 0)
	for _, item := range items {
		if item.Id == itemId {
			result = append(result, item)
			children := findChildren(items, item.Id)
			if len(children) > 0 {
				result = append(result, children...)
			}
			break
		}
	}
	return result
}

func findChildrenIdList(items []TreeOption, itemId int64) []int64 {
	result := make([]int64, 0)
	for _, item := range items {
		if item.Pid == itemId {
			result = append(result, item.Id)
			children := findChildrenIdList(items, item.Id)
			if len(children) > 0 {
				children = append(result, children...)
			}
		}
	}
	return result
}

func FindChildrenIdList(items []TreeOption, itemId int64) []int64 {
	result := make([]int64, 0)
	for _, item := range items {
		if item.Id == itemId {
			children := findChildrenIdList(items, item.Id)
			if len(children) > 0 {
				result = append(result, children...)
			}
			break
		}
	}
	return result
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

type ITreeNode interface {
	GetId() int64
	GetPid() int64
	GetLabel() string
	GetChildren() []ITreeNode
	Get(field string) interface{}
	Append(child ITreeNode)
}

type TreeNodeV2 struct {
	Id       int64       `json:"id"`
	Pid      int64       `json:"pid"`
	Label    string      `json:"label"`
	Children []ITreeNode `json:"children"`
}

func (node *TreeNodeV2) GetId() int64 {
	return node.Id
}
func (node *TreeNodeV2) GetPid() int64 {
	return node.Pid
}
func (node *TreeNodeV2) GetLabel() string {
	return node.Label
}
func (node *TreeNodeV2) Get(field string) interface{} {
	switch field {
	case "id":
		return node.Id
	case "pid":
		return node.Pid
	case "label":
		return node.Label
	}
	return nil
}
func (node *TreeNodeV2) GetChildren() []ITreeNode {
	return node.Children
}
func (node *TreeNodeV2) Append(child ITreeNode) {
	if node.Children == nil {
		node.Children = make([]ITreeNode, 0)
	}
	node.Children = append(node.Children, child)
}

func getTreeV2Children(items []ITreeNode, pid int64) []ITreeNode {
	result := make([]ITreeNode, 0)
	for _, item := range items {
		if item.GetPid() == pid {
			children := getTreeV2Children(items, item.GetId())
			for _, child := range children {
				item.Append(child)
			}
			result = append(result, item)
		}
	}
	return result
}

func ToTreeV2(items []ITreeNode) []ITreeNode {
	result := make([]ITreeNode, 0)
	for _, item := range items {
		if item.GetPid() == EmptyId {
			children := getTreeV2Children(items, item.GetId())
			for _, child := range children {
				item.Append(child)
			}
			result = append(result, item)
		}
	}
	return result
}
