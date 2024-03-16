package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

type PromptService struct {
}

var PromptServiceApp = new(PromptService)

//@function: CreatePrompt
//@description: 创建prompt
//@param: e model.Prompt
//@return: err error

func (exa *PromptService) CreatePrompt(e ai.Prompt) (err error) {
	e.BASEMODEL = BaseModel()
	err = global.GVA_DB.Create(&e).Error
	return err
}

//@function: DeleteFileChunk
//@description: 删除prompt
//@param: e model.Prompt
//@return: err error

func (exa *PromptService) DeletePrompt(e ai.Prompt) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@function: UpdatePrompt
//@description: 更新prompt
//@param: e *model.Prompt
//@return: err error

func (exa *PromptService) UpdatePrompt(e *ai.Prompt) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

//@function: GetPrompt
//@description: 获取prompt信息
//@param: id uint
//@return: customer model.Prompt, err error

func (exa *PromptService) GetPrompt(id uint64) (prompt ai.Prompt, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&prompt).Error
	return
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetPromptList
// @description: 分页获取prompt列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *PromptService) GetPromptList(sysUserAuthorityID uint, info request.PageInfo) (list interface{}, total int64, err error) {
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
	var promptList []ai.Prompt

	db := global.GVA_DB.Model(&ai.Prompt{})
	err = db.Count(&total).Error
	if err != nil {
		return promptList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&promptList).Error
	}
	return promptList, total, err
}

func (exa *PromptService) FindPromptByTopicAndType(topic string, promptType int) (prompt ai.Prompt, err error) {
	err = global.GVA_DB.Where("topic = ?", topic).Where("prompt_type=?", promptType).Last(&prompt).Error
	return

}
