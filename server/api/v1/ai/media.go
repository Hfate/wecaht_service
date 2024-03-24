package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MediaApi struct{}

// CreateMedia
// @Tags      Media
// @Summary   创建素材
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     data  body      wechat.Media            true  "素材用户名, 素材手机号码"
// @Success   200   {object}  response.Response{msg=string}  "创建素材"
// @Router    /media/media [post]
func (e *MediaApi) CreateMedia(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	targetAccountId := c.Request.PostFormValue("targetAccountId")

	if targetAccountId == "" {
		response.FailWithMessage("必须选定公众号", c)
	}

	err = mediaService.CreateMedia(targetAccountId, header)

	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteMedia
// @Tags      Media
// @Summary   删除素材
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Media            true  "素材ID"
// @Success   200   {object}  response.Response{msg=string}  "删除素材"
// @Router    /media/media [delete]
func (e *MediaApi) DeleteMedia(c *gin.Context) {
	var media ai.Media
	err := c.ShouldBindJSON(&media)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(media.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = mediaService.DeleteMedia(media)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetMediaPage
// @Tags      Media
// @Summary   分页获取权限素材列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限素材列表,返回包括列表,总数,页码,每页数量"
// @Router    /media/mediaPage [get]
func (e *MediaApi) GetMediaPage(c *gin.Context) {
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
	mediaList, total, err := mediaService.GetMediaPage(utils.GetUserAuthorityId(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     mediaList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
