package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

type PortalService struct {
}

//@function: CreatePortal
//@description: 创建门户
//@param: e model.Portal
//@return: err error

func (exa *PortalService) CreatePortal(e ai.Portal) (err error) {
	e.BASEMODEL = BaseModel()
	err = global.GVA_DB.Create(&e).Error
	return err
}

//@function: DeleteFileChunk
//@description: 删除门户
//@param: e model.Portal
//@return: err error

func (exa *PortalService) DeletePortal(e ai.Portal) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@function: UpdatePortal
//@description: 更新门户
//@param: e *model.Portal
//@return: err error

func (exa *PortalService) UpdatePortal(e *ai.Portal) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

//@function: GetPortal
//@description: 获取门户信息
//@param: id uint
//@return: customer model.Portal, err error

func (exa *PortalService) GetPortal(id uint64) (portal ai.Portal, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&portal).Error
	return
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetPortalList
// @description: 分页获取门户列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *PortalService) GetPortalList(sysUserAuthorityID uint, info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	var a system.SysAuthority
	a.AuthorityId = sysUserAuthorityID
	auth, err := systemService.AuthorityServiceApp.GetAuthorityInfo(a)
	if err != nil {
		return
	}
	var dataId []uint
	for _, v := range auth.DataAuthorityId {
		dataId = append(dataId, v.AuthorityId)
	}
	var portalList []ai.Portal

	db := global.GVA_DB.Model(&ai.Portal{})
	err = db.Count(&total).Error
	if err != nil {
		return portalList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&portalList).Error
	}
	return portalList, total, err
}
