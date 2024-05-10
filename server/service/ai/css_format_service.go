package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiRes "github.com/flipped-aurora/gin-vue-admin/server/model/ai/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type CssFormatService struct {
}

var CssFormatServiceApp = new(CssFormatService)

func (exa *CssFormatService) FindTopicList() []string {
	result := make([]string, 0)

	err := global.GVA_DB.Model(&ai.CssFormat{}).Distinct("topic").Find(&result).Error
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		return result
	}

	return result
}

//@function: CreateCssFormat
//@description: 创建公众号
//@param: e model.CssFormat
//@return: err error

func (exa *CssFormatService) CreateCssFormat(e ai.CssFormat) (err error) {
	e.BASEMODEL = BaseModel()
	err = global.GVA_DB.Create(&e).Error
	return err
}

//@function: DeleteFileChunk
//@description: 删除公众号
//@param: e model.CssFormat
//@return: err error

func (exa *CssFormatService) DeleteCssFormat(e ai.CssFormat) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@function: UpdateCssFormat
//@description: 更新公众号
//@param: e *model.CssFormat
//@return: err error

func (exa *CssFormatService) UpdateCssFormat(e *ai.CssFormat) (err error) {
	oldCss, err := exa.GetCssFormat(cast.ToUint64(e.ID))
	if err != nil {
		return err
	}

	// 更新名字
	global.GVA_DB.Model(&ai.OfficialAccount{}).Where("css_format=?", oldCss.FormatName).Update("css_format", e.FormatName)

	err = global.GVA_DB.Save(e).Error
	return err
}

//@function: GetCssFormat
//@description: 获取公众号信息
//@param: id uint
//@return: customer model.CssFormat, err error

func (exa *CssFormatService) GetCssFormat(id uint64) (cssFormat *ai.CssFormat, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&cssFormat).Error
	return
}

//@function: GetCssFormat
//@description: 获取公众号信息
//@param: id uint
//@return: customer model.CssFormat, err error

func (exa *CssFormatService) FindByFormatName(cssFormatName string) (cssFormat *ai.CssFormat, err error) {
	err = global.GVA_DB.Where("css_format_name = ?", cssFormatName).Last(&cssFormat).Error
	return
}

//@function: GetCssFormat
//@description: 获取公众号信息
//@param: id uint
//@return: customer model.CssFormat, err error

func (exa *CssFormatService) List() (list []*ai.CssFormat, err error) {
	err = global.GVA_DB.Where("1=1").Find(&list).Error
	return
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetCssFormatList
// @description: 分页获取公众号列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *CssFormatService) GetCssFormatList(sysUserAuthorityID uint, info request.PageInfo) (list interface{}, total int64, err error) {
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
	var cssFormatList []*ai.CssFormat

	db := global.GVA_DB.Model(&ai.CssFormat{})
	err = db.Count(&total).Error
	if err != nil {
		return cssFormatList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&cssFormatList).Error
	}

	result := make([]*aiRes.CssFormatResponse, 0)

	for _, v := range cssFormatList {
		result = append(result, &aiRes.CssFormatResponse{
			CssFormat: v,
		})
	}

	return result, total, err
}
