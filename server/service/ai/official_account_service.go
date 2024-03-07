package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"time"
)

type OfficialAccountService struct {
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateOfficialAccount
//@description: 创建门户
//@param: e model.OfficialAccount
//@return: err error

func (exa *OfficialAccountService) CreateOfficialAccount(e wechat.OfficialAccount) (err error) {
	e.GVA_MODEL = global.GVA_MODEL{ID: utils.GenID(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err = global.GVA_DB.Create(&e).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteFileChunk
//@description: 删除门户
//@param: e model.OfficialAccount
//@return: err error

func (exa *OfficialAccountService) DeleteOfficialAccount(e wechat.OfficialAccount) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateOfficialAccount
//@description: 更新门户
//@param: e *model.OfficialAccount
//@return: err error

func (exa *OfficialAccountService) UpdateOfficialAccount(e *wechat.OfficialAccount) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetOfficialAccount
//@description: 获取门户信息
//@param: id uint
//@return: customer model.OfficialAccount, err error

func (exa *OfficialAccountService) GetOfficialAccount(id uint64) (officialAccount wechat.OfficialAccount, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&officialAccount).Error
	return
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetOfficialAccountList
// @description: 分页获取门户列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *OfficialAccountService) GetOfficialAccountList(sysUserAuthorityID uint, info request.PageInfo) (list interface{}, total int64, err error) {
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
	var officialAccountList []wechat.OfficialAccount

	db := global.GVA_DB.Model(&wechat.OfficialAccount{})
	err = db.Count(&total).Error
	if err != nil {
		return officialAccountList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Find(&officialAccountList).Error
	}
	return officialAccountList, total, err
}
