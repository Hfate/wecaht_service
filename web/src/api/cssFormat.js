import service from '@/utils/request'
// @Tags CssFormatApi
// @Summary 创建公众号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.CssFormat true "创建公众号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cssFormat/cssFormat [post]
export const createCssFormat = (data) => {
    return service({
        url: '/cssFormat/cssFormat',
        method: 'post',
        data
    })
}


// @Tags CssFormatApi
// @Summary 更新公众号信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.CssFormat true "更新公众号信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cssFormat/cssFormat [put]
export const updateCssFormat = (data) => {
    return service({
        url: '/cssFormat/cssFormat',
        method: 'put',
        data
    })
}

// @Tags CssFormatApi
// @Summary 删除公众号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.CssFormat true "删除公众号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cssFormat/cssFormat [delete]
export const deleteCssFormat = (data) => {
    return service({
        url: '/cssFormat/cssFormat',
        method: 'delete',
        data
    })
}


// @Tags CssFormatApi
// @Summary 公众号创建文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.CssFormat true "公众号创建文章"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cssFormat/create [post]
export const cssFormatCreate = (data) => {
    return service({
        url: '/cssFormat/create',
        method: 'post',
        data
    })
}

// @Tags CssFormatApi
// @Summary 获取单一公众号信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.CssFormat true "获取单一公众号信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cssFormat/cssFormat [get]
export const getCssFormat = (params) => {
    return service({
        url: '/cssFormat/cssFormat',
        method: 'get',
        params
    })
}

// @Tags CssFormatApi
// @Summary 获取权限公众号列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限公众号列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cssFormat/cssFormatList [get]
export const getCssFormatList = (params) => {
    return service({
        url: '/cssFormat/cssFormatList',
        method: 'get',
        params
    })
}
