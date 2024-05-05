import service from '@/utils/request'


// @Tags TemplateApi
// @Summary 删除模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Template true "生成模板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /template/generate [delete]
export const generateTemplate = (data) => {
    return service({
        url: '/template/generate',
        method: 'post',
    })
}


// @Tags TemplateApi
// @Summary 删除模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Template true "删除模板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /template/template [delete]
export const deleteTemplate = (data) => {
    return service({
        url: '/template/template',
        method: 'delete',
        data
    })
}

// @Tags TemplateApi
// @Summary 克隆模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Template true "克隆"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /template/clone [post]
export const cloneTemplate = (data) => {
    return service({
        url: '/template/clone',
        method: 'post',
        data
    })
}


// @Tags TemplateApi
// @Summary 获取单一模板信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Template true "获取单一模板信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /template/template [get]
export const getTemplate = (params) => {
    return service({
        url: '/template/template',
        method: 'get',
        params
    })
}

// @Tags TemplateApi
// @Summary 获取权限模板列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限模板列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /template/templateList [get]
export const getTemplateList = (params) => {
    return service({
        url: '/template/templateList',
        method: 'get',
        params
    })
}


// @Tags TemplateApi
// @Summary 获取权限模板列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限模板列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /template/update [get]
export const updateTemplate = (data) => {
    return service({
        url: '/template/update',
        method: 'post',
        data
    })
}
