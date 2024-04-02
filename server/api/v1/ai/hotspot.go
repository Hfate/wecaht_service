package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HotspotApi struct{}

// DeleteHotspot
// @Tags      Hotspot
// @Summary   删除头条
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Hotspot            true  "头条ID"
// @Success   200   {object}  response.Response{msg=string}  "删除头条"
// @Router    /hotspot/hotspot [delete]
func (e *HotspotApi) DeleteHotspot(c *gin.Context) {
	var hotspot ai.Hotspot
	err := c.ShouldBindJSON(&hotspot)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(hotspot.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = hotspotService.DeleteHotspot(hotspot)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// CreateArticle
// @Tags      Hotspot
// @Summary   删除头条
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Hotspot            true  "头条ID"
// @Success   200   {object}  response.Response{msg=string}  "删除头条"
// @Router    /hotspot/create [post]
func (e *HotspotApi) CreateArticle(c *gin.Context) {
	var hotspot ai.Hotspot
	err := c.ShouldBindJSON(&hotspot)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(hotspot.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = hotspotService.CreateArticle(hotspot.ID)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// GetHotspotList
// @Tags      Hotspot
// @Summary   分页获取权限头条列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限头条列表,返回包括列表,总数,页码,每页数量"
// @Router    /hotspot/hotspotList [get]
func (e *HotspotApi) GetHotspotList(c *gin.Context) {
	var pageInfo aiReq.HotspotSearch
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
	hotspotList, total, err := hotspotService.GetHotspotList(utils.GetUserAuthorityId(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     hotspotList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
