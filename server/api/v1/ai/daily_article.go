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
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type DailyArticleApi struct{}

// DeleteDailyArticle
// @Tags      DailyArticle
// @Summary   删除文章
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.DailyArticle            true  "文章ID"
// @Success   200   {object}  response.Response{msg=string}  "删除文章"
// @Router    /dailyArticle/dailyArticle [delete]
func (e *DailyArticleApi) DeleteDailyArticle(c *gin.Context) {
	var article ai.DailyArticle
	err := c.ShouldBindJSON(&article)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(article.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dailyArticleService.DeleteDailyArticle(article)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteDailyArticlesByIds
// @Tags      DailyArticle
// @Summary   删除选中文章
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.IdsReq                 true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "删除选中文章"
// @Router    /dailyArticle/deleteDailyArticlesByIds [delete]
func (e *DailyArticleApi) DeleteDailyArticlesByIds(c *gin.Context) {
	var ids request.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dailyArticleService.DeleteDailyArticlesByIds(ids)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetDailyArticle
// @Tags      DailyArticle
// @Summary   获取单一文章信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     wechat.DailyArticle                                                true  "文章ID"
// @Success   200   {object}  response.Response{data=exampleRes.DailyArticleResponse,msg=string}  "获取单一文章信息,返回包括文章详情"
// @Router    /dailyArticle/dailyArticle [get]
func (e *DailyArticleApi) GetDailyArticle(c *gin.Context) {
	var article ai.DailyArticle
	err := c.ShouldBindQuery(&article)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(article.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := dailyArticleService.GetDailyArticle(cast.ToUint64(article.ID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(aiRes.DailyArticleResponse{DailyArticle: data}, "获取成功", c)
}

// Recreation
// @Tags      DailyArticle
// @Summary   改写文章
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.DailyArticle            true  "文章ID"
// @Success   200   {object}  response.Response{msg=string}  "改写文章"
// @Router    /dailyArticle/recreation [post]
func (e *DailyArticleApi) Recreation(c *gin.Context) {
	var article ai.DailyArticle
	err := c.ShouldBindJSON(&article)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(article.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dailyArticleService.Recreation(cast.ToUint64(article.ID))
	if err != nil {
		global.GVA_LOG.Error("改写失败!", zap.Error(err))
		response.FailWithMessage("改写失败", c)
		return
	}
	response.OkWithMessage("改写成功", c)
}

// GetDailyArticleList
// @Tags      DailyArticle
// @Summary   分页获取权限文章列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限文章列表,返回包括列表,总数,页码,每页数量"
// @Router    /dailyArticle/dailyArticleList [get]
func (e *DailyArticleApi) GetDailyArticleList(c *gin.Context) {
	var pageInfo aiReq.DailyArticleSearch
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
	articleList, total, err := dailyArticleService.GetDailyArticleList(utils.GetUserAuthorityId(c))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     articleList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
