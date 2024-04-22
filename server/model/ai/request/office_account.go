package request

// Modify  user's auth structure
type SetCreateTypes struct {
	ID             string
	CreateTypeList []int `json:"createTypeList"` // 角色ID
}
