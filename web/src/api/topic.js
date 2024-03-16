import service from '@/utils/request'
// @Tags TopicApi
// @Summary 创建主题
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Topic true "创建主题"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topic/topic [post]
export const createTopic = (data) => {
    return service({
        url: '/topic/topic',
        method: 'post',
        data
    })
}

// @Tags TopicApi
// @Summary 更新主题信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Topic true "更新主题信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topic/topic [put]
export const updateTopic = (data) => {
    return service({
        url: '/topic/topic',
        method: 'put',
        data
    })
}

// @Tags TopicApi
// @Summary 删除主题
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Topic true "删除主题"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topic/topic [delete]
export const deleteTopic = (data) => {
    return service({
        url: '/topic/topic',
        method: 'delete',
        data
    })
}


// @Tags TopicApi
// @Summary 获取主题集合
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取主题集合"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topic/topicList [get]
export const getTopicList = (params) => {
    return service({
        url: '/topic/topicList',
        method: 'get',
        params
    })
}


// @Tags TopicApi
// @Summary 获取主题分页
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取主题分页"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topic/topicList [get]
export const getTopicPage = (params) => {
    return service({
        url: '/topic/topicPage',
        method: 'get',
        params
    })
}
