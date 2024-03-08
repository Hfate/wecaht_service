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

type OfficialAccountApi struct{}

// CreateOfficialAccount
// @Tags      OfficialAccount
// @Summary   创建门户
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.OfficialAccount            true  "门户用户名, 门户手机号码"
// @Success   200   {object}  response.Response{msg=string}  "创建门户"
// @Router    /officialAccount/officialAccount [post]
func (e *OfficialAccountApi) CreateOfficialAccount(c *gin.Context) {
	var officialAccount ai.OfficialAccount
	err := c.ShouldBindJSON(&officialAccount)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(officialAccount, utils.OfficialAccountVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = officialAccountService.CreateOfficialAccount(officialAccount)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteOfficialAccount
// @Tags      OfficialAccount
// @Summary   删除门户
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.OfficialAccount            true  "门户ID"
// @Success   200   {object}  response.Response{msg=string}  "删除门户"
// @Router    /officialAccount/officialAccount [delete]
func (e *OfficialAccountApi) DeleteOfficialAccount(c *gin.Context) {
	var officialAccount ai.OfficialAccount
	err := c.ShouldBindJSON(&officialAccount)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(officialAccount.GVA_MODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = officialAccountService.DeleteOfficialAccount(officialAccount)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateOfficialAccount
// @Tags      OfficialAccount
// @Summary   更新门户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.OfficialAccount            true  "门户ID, 门户信息"
// @Success   200   {object}  response.Response{msg=string}  "更新门户信息"
// @Router    /officialAccount/officialAccount [put]
func (e *OfficialAccountApi) UpdateOfficialAccount(c *gin.Context) {
	var officialAccount ai.OfficialAccount
	err := c.ShouldBindJSON(&officialAccount)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(officialAccount.GVA_MODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(officialAccount, utils.OfficialAccountVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = officialAccountService.UpdateOfficialAccount(&officialAccount)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetOfficialAccount
// @Tags      OfficialAccount
// @Summary   获取单一门户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     wechat.OfficialAccount                                                true  "门户ID"
// @Success   200   {object}  response.Response{data=exampleRes.OfficialAccountResponse,msg=string}  "获取单一门户信息,返回包括门户详情"
// @Router    /officialAccount/officialAccount [get]
func (e *OfficialAccountApi) GetOfficialAccount(c *gin.Context) {
	var officialAccount ai.OfficialAccount
	err := c.ShouldBindQuery(&officialAccount)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(officialAccount.GVA_MODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := officialAccountService.GetOfficialAccount(officialAccount.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(aiRes.OfficialAccountResponse{OfficialAccount: data}, "获取成功", c)
}

// GetOfficialAccountList
// @Tags      OfficialAccount
// @Summary   分页获取权限门户列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限门户列表,返回包括列表,总数,页码,每页数量"
// @Router    /officialAccount/officialAccountList [get]
func (e *OfficialAccountApi) GetOfficialAccountList(c *gin.Context) {
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
	officialAccountList, total, err := officialAccountService.GetOfficialAccountList(utils.GetUserAuthorityId(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     officialAccountList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
