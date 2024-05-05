package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/spf13/cast"
)

type TemplateService struct {
}

var TemplateServiceApp = new(TemplateService)

func (exa *TemplateService) Create(template ai.Template) error {
	return global.GVA_DB.Model(&ai.Template{}).Create(template).Error
}

func (exa *TemplateService) FindByAccountId(appId string) *ai.Template {
	result := &ai.Template{}
	global.GVA_DB.Model(&ai.Template{}).Where("account_id=?", appId).Last(&result)
	return result
}

func (exa *TemplateService) Update(template ai.Template) error {
	account, _ := OfficialAccountServiceApp.FindByAppId(template.AccountId)
	template.AccountName = account.AccountName
	return global.GVA_DB.Save(template).Error
}

func (exa *TemplateService) Clone(template ai.Template) error {
	global.GVA_DB.Model(&ai.Template{}).Where("id=?", template.ID).Last(&template)
	template.ID = cast.ToString(utils.GenID())
	return global.GVA_DB.Model(&ai.Template{}).Create(template).Error
}

//@function: DeleteFileChunk
//@description: 删除文章
//@param: e model.Template
//@return: err error

func (exa *TemplateService) DeleteTemplate(e ai.Template) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@function: GetTemplate
//@description: 获取文章信息
//@param: id uint
//@return: customer model.Template, err error

func (exa *TemplateService) GetTemplate(id uint64) (template ai.Template, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&template).Error
	return
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetTemplateList
// @description: 分页获取文章列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *TemplateService) GetTemplateList(sysUserAuthorityID uint, info request.PageInfo) (list interface{}, total int64, err error) {
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
	var templateList []ai.Template

	db := global.GVA_DB.Model(&ai.Template{})

	err = db.Count(&total).Error
	if err != nil {
		return templateList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&templateList).Error
	}
	return templateList, total, err
}
