package request

type SysAuthorityBtnReq struct {
	MenuID      string `json:"menuID"`
	AuthorityId uint   `json:"authorityId"`
	Selected    []uint `json:"selected"`
}
