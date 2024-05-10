package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiRes "github.com/flipped-aurora/gin-vue-admin/server/model/ai/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type CssFormatApi struct{}

// CreateCssFormat
// @Tags      CssFormat
// @Summary   创建排版
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.CssFormat            true  "排版用户名, 排版手机号码"
// @Success   200   {object}  response.Response{msg=string}  "创建排版"
// @Router    /cssFormat/cssFormat [post]
func (e *CssFormatApi) CreateCssFormat(c *gin.Context) {
	var cssFormat ai.CssFormat
	err := c.ShouldBindJSON(&cssFormat)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(cssFormat, utils.CssFormatVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = cssFormatService.CreateCssFormat(cssFormat)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteCssFormat
// @Tags      CssFormat
// @Summary   删除排版
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.CssFormat            true  "排版ID"
// @Success   200   {object}  response.Response{msg=string}  "删除排版"
// @Router    /cssFormat/cssFormat [delete]
func (e *CssFormatApi) DeleteCssFormat(c *gin.Context) {
	var cssFormat ai.CssFormat
	err := c.ShouldBindJSON(&cssFormat)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(cssFormat.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = cssFormatService.DeleteCssFormat(cssFormat)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateCssFormat
// @Tags      CssFormat
// @Summary   更新排版信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.CssFormat            true  "排版ID, 排版信息"
// @Success   200   {object}  response.Response{msg=string}  "更新排版信息"
// @Router    /cssFormat/cssFormat [put]
func (e *CssFormatApi) UpdateCssFormat(c *gin.Context) {
	var cssFormat ai.CssFormat
	err := c.ShouldBindJSON(&cssFormat)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(cssFormat.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(cssFormat, utils.CssFormatVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = cssFormatService.UpdateCssFormat(&cssFormat)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetCssFormat
// @Tags      CssFormat
// @Summary   获取单一排版信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     wechat.CssFormat                                                true  "排版ID"
// @Success   200   {object}  response.Response{data=exampleRes.CssFormatResponse,msg=string}  "获取单一排版信息,返回包括排版详情"
// @Router    /cssFormat/cssFormat [get]
func (e *CssFormatApi) GetCssFormat(c *gin.Context) {
	var cssFormat ai.CssFormat
	err := c.ShouldBindQuery(&cssFormat)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(cssFormat.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := cssFormatService.GetCssFormat(cast.ToUint64(cssFormat.ID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(aiRes.CssFormatResponse{CssFormat: data}, "获取成功", c)
}

// GetCssFormatList
// @Tags      CssFormat
// @Summary   分页获取权限排版列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限排版列表,返回包括列表,总数,页码,每页数量"
// @Router    /cssFormat/cssFormatList [get]
func (e *CssFormatApi) GetCssFormatList(c *gin.Context) {
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
	cssFormatList, total, err := cssFormatService.GetCssFormatList(utils.GetUserAuthorityId(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     cssFormatList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
