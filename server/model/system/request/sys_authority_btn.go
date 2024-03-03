package request

type SysAuthorityBtnReq struct {
	MenuID      uint64 `json:"menuID"`
	AuthorityId uint   `json:"authorityId"`
	Selected    []uint `json:"selected"`
}
