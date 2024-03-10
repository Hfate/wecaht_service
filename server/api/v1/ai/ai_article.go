package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	aiRes "github.com/flipped-aurora/gin-vue-admin/server/model/ai/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AIArticleApi struct{}

// DeleteAIArticle
// @Tags      AIArticle
// @Summary   删除文章
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.AIArticle            true  "文章ID"
// @Success   200   {object}  response.Response{msg=string}  "删除文章"
// @Router    /aiArticle/aiArticle [delete]
func (e *AIArticleApi) DeleteAIArticle(c *gin.Context) {
	var aiArticle ai.AIArticle
	err := c.ShouldBindJSON(&aiArticle)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(aiArticle.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = aiArticleService.DeleteAIArticle(aiArticle)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteAIArticlesByIds
// @Tags      AIArticle
// @Summary   删除选中文章
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.IdsReq                 true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "删除选中文章"
// @Router    /aiArticle/deleteAIArticlesByIds [delete]
func (e *AIArticleApi) DeleteAIArticlesByIds(c *gin.Context) {
	var ids request.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = aiArticleService.DeleteAIArticlesByIds(ids)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetAIArticle
// @Tags      AIArticle
// @Summary   获取单一文章信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     wechat.AIArticle                                                true  "文章ID"
// @Success   200   {object}  response.Response{data=exampleRes.AIArticleResponse,msg=string}  "获取单一文章信息,返回包括文章详情"
// @Router    /aiArticle/aiArticle [get]
func (e *AIArticleApi) GetAIArticle(c *gin.Context) {
	var aiArticle ai.AIArticle
	err := c.ShouldBindQuery(&aiArticle)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(aiArticle.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := aiArticleService.GetAIArticle(aiArticle.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(aiRes.AIArticleResponse{AIArticle: data}, "获取成功", c)
}

// Recreation
// @Tags      AIArticle
// @Summary   改写文章
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.AIArticle            true  "文章ID"
// @Success   200   {object}  response.Response{msg=string}  "改写文章"
// @Router    /aiArticle/recreation [post]
func (e *AIArticleApi) Recreation(c *gin.Context) {
	var aiArticle ai.AIArticle
	err := c.ShouldBindJSON(&aiArticle)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(aiArticle.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = aiArticleService.Recreation(aiArticle.ID)
	if err != nil {
		global.GVA_LOG.Error("改写失败!", zap.Error(err))
		response.FailWithMessage("改写失败", c)
		return
	}
	response.OkWithMessage("改写成功", c)
}

// GetAIArticleList
// @Tags      AIArticle
// @Summary   分页获取权限文章列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限文章列表,返回包括列表,总数,页码,每页数量"
// @Router    /aiArticle/aiArticleList [get]
func (e *AIArticleApi) GetAIArticleList(c *gin.Context) {
	var pageInfo aiReq.AIArticleSearch
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
	aiArticleList, total, err := aiArticleService.GetAIArticleList(utils.GetUserAuthorityId(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     aiArticleList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
