package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"time"
)

type HotspotService struct {
}

//@function: CreateHotspot
//@description: 创建热点
//@param: e model.Hotspot
//@return: err error

func (exa *HotspotService) CreateHotspot(e ai.Hotspot) (err error) {
	e.BASEMODEL = global.BASEMODEL{ID: utils.GenID(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err = global.GVA_DB.Create(&e).Error
	return err
}

//@function: DeleteFileChunk
//@description: 删除热点
//@param: e model.Hotspot
//@return: err error

func (exa *HotspotService) DeleteHotspot(e ai.Hotspot) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@function: UpdateHotspot
//@description: 更新热点
//@param: e *model.Hotspot
//@return: err error

func (exa *HotspotService) UpdateHotspot(e *ai.Hotspot) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

//@function: GetHotspot
//@description: 获取热点信息
//@param: id uint
//@return: customer model.Hotspot, err error

func (exa *HotspotService) GetHotspot(id uint64) (portal ai.Hotspot, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&portal).Error
	return
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetHotspotList
// @description: 分页获取热点列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *HotspotService) GetHotspotList(sysUserAuthorityID uint, info aiReq.HotspotSearch) (list interface{}, total int64, err error) {
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
	var portalList []ai.Hotspot

	db := global.GVA_DB.Model(&ai.Hotspot{})
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Headline != "" {
		db = db.Where("headline LIKE ?", "%"+info.Headline+"%")
	}
	if info.PortalName != "" {
		db = db.Where("portal_name LIKE ?", "%"+info.PortalName+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return portalList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("created_at desc,trending desc").Find(&portalList).Error
	}
	return portalList, total, err
}
