package request

import "github.com/spf13/cast"

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`   //关键字
}

// GetById Find by id structure
type GetById struct {
	ID string `json:"id" form:"id"` // 主键ID
}

func (r *GetById) Uint() uint {
	return cast.ToUint(r.ID)
}

func (r *GetById) Uint64() uint64 {
	return cast.ToUint64(r.ID)
}

type IdsReq struct {
	Ids []string `json:"ids" form:"ids"`
}

// GetAuthorityId Get role by id structure
type GetAuthorityId struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
}

type Empty struct{}
