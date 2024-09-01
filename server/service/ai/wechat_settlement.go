package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type WechatSettlementService struct {
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetSettlementList
// @description: 分页获取门户列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *WechatSettlementService) GetSettlementList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	var settlementList []ai.WechatSettlement

	db := global.GVA_DB.Model(&ai.WechatSettlement{})
	err = db.Count(&total).Error
	if err != nil {
		return settlementList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("date desc,account_id desc").Find(&settlementList).Error
	}
	return settlementList, total, err
}
