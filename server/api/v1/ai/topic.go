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

type TopicApi struct{}

// CreateTopic
// @Tags      Topic
// @Summary   创建主题
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Topic            true  "主题用户名, 主题手机号码"
// @Success   200   {object}  response.Response{msg=string}  "创建主题"
// @Router    /topic/topic [post]
func (e *TopicApi) CreateTopic(c *gin.Context) {
	var topic ai.Topic
	err := c.ShouldBindJSON(&topic)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(topic, utils.TopicVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = topicService.CreateTopic(topic)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteTopic
// @Tags      Topic
// @Summary   删除主题
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Topic            true  "主题ID"
// @Success   200   {object}  response.Response{msg=string}  "删除主题"
// @Router    /topic/topic [delete]
func (e *TopicApi) DeleteTopic(c *gin.Context) {
	var topic ai.Topic
	err := c.ShouldBindJSON(&topic)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(topic.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = topicService.DeleteTopic(topic)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateTopic
// @Tags      Topic
// @Summary   更新主题信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Topic            true  "主题ID, 主题信息"
// @Success   200   {object}  response.Response{msg=string}  "更新主题信息"
// @Router    /topic/topic [put]
func (e *TopicApi) UpdateTopic(c *gin.Context) {
	var topic ai.Topic
	err := c.ShouldBindJSON(&topic)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(topic.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(topic, utils.TopicVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = topicService.UpdateTopic(&topic)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetTopicPage
// @Tags      Topic
// @Summary   分页获取权限主题列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限主题列表,返回包括列表,总数,页码,每页数量"
// @Router    /topic/topicPage [get]
func (e *TopicApi) GetTopicPage(c *gin.Context) {
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
	topicList, total, err := topicService.GetTopicPage(utils.GetUserAuthorityId(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     topicList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetTopicList
// @Tags      Topic
// @Summary   分页获取权限主题列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限主题列表,返回包括列表,总数,页码,每页数量"
// @Router    /topic/topicList [get]
func (e *TopicApi) GetTopicList(c *gin.Context) {
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
	topicList, err := topicService.GetTopicList(utils.GetUserAuthorityId(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     topicList,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
