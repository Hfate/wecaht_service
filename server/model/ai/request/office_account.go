package request

// Modify  user's auth structure
type SetCreateTypes struct {
	ID             uint64
	CreateTypeList []int `json:"createTypeList"` // 角色ID
}
