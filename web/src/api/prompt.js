import service from '@/utils/request'
// @Tags PromptApi
// @Summary 创建prompt
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Prompt true "创建prompt"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /prompt/prompt [post]
export const createPrompt = (data) => {
    return service({
        url: '/prompt/prompt',
        method: 'post',
        data
    })
}

// @Tags PromptApi
// @Summary 更新prompt信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Prompt true "更新prompt信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /prompt/prompt [put]
export const updatePrompt = (data) => {
    return service({
        url: '/prompt/prompt',
        method: 'put',
        data
    })
}

// @Tags PromptApi
// @Summary 删除prompt
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Prompt true "删除prompt"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /prompt/prompt [delete]
export const deletePrompt = (data) => {
    return service({
        url: '/prompt/prompt',
        method: 'delete',
        data
    })
}

// @Tags PromptApi
// @Summary 获取单一prompt信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Prompt true "获取单一prompt信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /prompt/prompt [get]
export const getPrompt = (params) => {
    return service({
        url: '/prompt/prompt',
        method: 'get',
        params
    })
}

// @Tags PromptApi
// @Summary 获取权限prompt列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限prompt列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /prompt/promptList [get]
export const getPromptList = (params) => {
    return service({
        url: '/prompt/promptList',
        method: 'get',
        params
    })
}
