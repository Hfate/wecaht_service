package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/wechat/request"
	aiRes "github.com/flipped-aurora/gin-vue-admin/server/model/wechat/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ArticleApi struct{}

// DeleteArticle
// @Tags      Article
// @Summary   删除文章
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Article            true  "文章ID"
// @Success   200   {object}  response.Response{msg=string}  "删除文章"
// @Router    /article/article [delete]
func (e *ArticleApi) DeleteArticle(c *gin.Context) {
	var portal wechat.Article
	err := c.ShouldBindJSON(&portal)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(portal.GVA_MODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = articleService.DeleteArticle(portal)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetArticle
// @Tags      Article
// @Summary   获取单一文章信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     wechat.Article                                                true  "文章ID"
// @Success   200   {object}  response.Response{data=exampleRes.ArticleResponse,msg=string}  "获取单一文章信息,返回包括文章详情"
// @Router    /article/article [get]
func (e *ArticleApi) GetArticle(c *gin.Context) {
	var portal wechat.Article
	err := c.ShouldBindQuery(&portal)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(portal.GVA_MODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := articleService.GetArticle(portal.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(aiRes.ArticleResponse{Article: data}, "获取成功", c)
}

// GetArticleList
// @Tags      Article
// @Summary   分页获取权限文章列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限文章列表,返回包括列表,总数,页码,每页数量"
// @Router    /article/articleList [get]
func (e *ArticleApi) GetArticleList(c *gin.Context) {
	var pageInfo aiReq.ArticleSearch
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
	portalList, total, err := articleService.GetArticleList(utils.GetUserAuthorityId(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     portalList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
