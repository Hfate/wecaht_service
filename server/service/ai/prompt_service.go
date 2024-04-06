package ai

import (
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiRes "github.com/flipped-aurora/gin-vue-admin/server/model/ai/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/spf13/cast"
	"go.uber.org/zap"
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

func (exa *PromptService) GetPrompt(id uint64) (promptResp aiRes.PromptResp, err error) {
	prompt := &ai.Prompt{}
	err = global.GVA_DB.Where("id = ?", id).First(&prompt).Error
	if err != nil {
		return aiRes.PromptResp{}, err
	}

	// 定义一个空的字符串切片
	promptList := make([]string, 0)

	// 使用json.Unmarshal将JSON字符串解析到字符串切片
	err = json.Unmarshal([]byte(utils.EscapeSpecialCharacters(prompt.Prompt)), &promptList)

	promptResp = aiRes.PromptResp{
		BASEMODEL:  prompt.BASEMODEL,
		Topic:      prompt.Topic,
		PromptType: prompt.PromptType,
		PromptList: promptList,
		Language:   prompt.Language,
	}

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

func (exa *PromptService) FindPromptByTopicAndType(topic string, promptType int) (prompt string, err error) {

	promptModel := &ai.Prompt{}

	err = global.GVA_DB.Where("topic = ?", topic).Where("prompt_type=?", promptType).Last(&promptModel).Error

	// 没找到 则使用默认的
	if err != nil {
		err = global.GVA_DB.Where("topic = ?", "default").Where("prompt_type=?", promptType).Last(&promptModel).Error
		if err != nil {
			global.GVA_LOG.Info("无法找到topic相关的prompt", zap.Error(err), zap.String("topic", topic))
			return "", errors.New("无法找到topic相关的prompt,promptType=" + cast.ToString(promptType))
		}
	}

	prompt = promptModel.Prompt
	return

}
