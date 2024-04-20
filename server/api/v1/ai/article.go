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

type ArticleApi struct{}

// UploadArticle
// @Tags      Article
// @Summary   上传文章
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Article            true  "文章ID"
// @Success   200   {object}  response.Response{msg=string}  "上传文章"
// @Router    /article/upload [delete]
func (e *ArticleApi) UploadArticle(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	err = articleService.UploadArticle(header)
	if err != nil {
		global.GVA_LOG.Error("上传失败!", zap.Error(err))
		response.FailWithMessage("上传失败", c)
		return
	}
	response.OkWithMessage("上传成功", c)

}

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
	var article ai.Article
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
	err = articleService.DeleteArticle(article)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteArticlesByIds
// @Tags      Article
// @Summary   删除选中文章
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.IdsReq                 true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "删除选中文章"
// @Router    /article/deleteArticlesByIds [delete]
func (e *ArticleApi) DeleteArticlesByIds(c *gin.Context) {
	var ids request.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = articleService.DeleteArticlesByIds(ids)
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
	var article ai.Article
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
	data, err := articleService.GetArticle(cast.ToUint64(article.ID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(aiRes.ArticleResponse{Article: data}, "获取成功", c)
}

// Recreation
// @Tags      Article
// @Summary   改写文章
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Article            true  "文章ID"
// @Success   200   {object}  response.Response{msg=string}  "改写文章"
// @Router    /article/recreation [post]
func (e *ArticleApi) Recreation(c *gin.Context) {
	var article ai.Article
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
	err = articleService.Recreation(cast.ToUint64(article.ID))
	if err != nil {
		global.GVA_LOG.Error("改写失败!", zap.Error(err))
		response.FailWithMessage("改写失败", c)
		return
	}
	response.OkWithMessage("改写成功", c)
}

// Template
// @Tags      Article
// @Summary   获取文章上传模板
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=exampleRes.ArticleResponse,msg=string}  "获取文章上传模板"
// @Router    /article/template [get]
func (e *ArticleApi) Template(c *gin.Context) {
	var filename = "article_template.xlsx"
	var filePath = "./template/" + filename
	//返回文件流
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filePath)
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
	articleList, total, err := articleService.GetArticleList(utils.GetUserAuthorityId(c), pageInfo)
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

// Download
// @Tags      Article
// @Summary   分页获取权限文章列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限文章列表,返回包括列表,总数,页码,每页数量"
// @Router    /article/download [get]
func (e *ArticleApi) Download(c *gin.Context) {
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

	pageInfo.PageSize = 200
	articleService.Download(c, utils.GetUserAuthorityId(c), pageInfo)
}
