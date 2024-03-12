package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiRes "github.com/flipped-aurora/gin-vue-admin/server/model/ai/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PromptApi struct{}

// CreatePrompt
// @Tags      Prompt
// @Summary   创建prompt
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Prompt            true  "prompt用户名, prompt手机号码"
// @Success   200   {object}  response.Response{msg=string}  "创建prompt"
// @Router    /prompt/prompt [post]
func (e *PromptApi) CreatePrompt(c *gin.Context) {
	var prompt ai.Prompt
	err := c.ShouldBindJSON(&prompt)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(prompt, utils.PromptVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = promptService.CreatePrompt(prompt)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeletePrompt
// @Tags      Prompt
// @Summary   删除prompt
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Prompt            true  "promptID"
// @Success   200   {object}  response.Response{msg=string}  "删除prompt"
// @Router    /prompt/prompt [delete]
func (e *PromptApi) DeletePrompt(c *gin.Context) {
	var prompt ai.Prompt
	err := c.ShouldBindJSON(&prompt)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(prompt.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = promptService.DeletePrompt(prompt)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdatePrompt
// @Tags      Prompt
// @Summary   更新prompt信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Prompt            true  "promptID, prompt信息"
// @Success   200   {object}  response.Response{msg=string}  "更新prompt信息"
// @Router    /prompt/prompt [put]
func (e *PromptApi) UpdatePrompt(c *gin.Context) {
	var prompt ai.Prompt
	err := c.ShouldBindJSON(&prompt)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(prompt.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(prompt, utils.PromptVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = promptService.UpdatePrompt(&prompt)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetPrompt
// @Tags      Prompt
// @Summary   获取单一prompt信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     wechat.Prompt                                                true  "promptID"
// @Success   200   {object}  response.Response{data=exampleRes.PromptResponse,msg=string}  "获取单一prompt信息,返回包括prompt详情"
// @Router    /prompt/prompt [get]
func (e *PromptApi) GetPrompt(c *gin.Context) {
	var prompt ai.Prompt
	err := c.ShouldBindQuery(&prompt)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(prompt.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := promptService.GetPrompt(prompt.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(aiRes.PromptResponse{Prompt: data}, "获取成功", c)
}

// GetPromptList
// @Tags      Prompt
// @Summary   分页获取权限prompt列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限prompt列表,返回包括列表,总数,页码,每页数量"
// @Router    /prompt/promptList [get]
func (e *PromptApi) GetPromptList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	promptList, total, err := promptService.GetPromptList(utils.GetUserAuthorityId(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     promptList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
