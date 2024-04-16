import service from '@/utils/request'
// @Tags OfficialAccountApi
// @Summary 创建公众号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.OfficialAccount true "创建公众号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /officialAccount/officialAccount [post]
export const createOfficialAccount = (data) => {
    return service({
        url: '/officialAccount/officialAccount',
        method: 'post',
        data
    })
}


// @Tags OfficialAccountApi
// @Summary 创建公众号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.OfficialAccount true "创建公众号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /officialAccount/officialAccount [post]
export const setCreateTypes = (data) => {
    return service({
        url: '/officialAccount/updateCreateTypes',
        method: 'put',
        data
    })
}

// @Tags OfficialAccountApi
// @Summary 更新公众号信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.OfficialAccount true "更新公众号信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /officialAccount/officialAccount [put]
export const updateOfficialAccount = (data) => {
    return service({
        url: '/officialAccount/officialAccount',
        method: 'put',
        data
    })
}

// @Tags OfficialAccountApi
// @Summary 删除公众号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.OfficialAccount true "删除公众号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /officialAccount/officialAccount [delete]
export const deleteOfficialAccount = (data) => {
    return service({
        url: '/officialAccount/officialAccount',
        method: 'delete',
        data
    })
}


// @Tags OfficialAccountApi
// @Summary 公众号创建文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.OfficialAccount true "公众号创建文章"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /officialAccount/create [post]
export const officialAccountCreate = (data) => {
    return service({
        url: '/officialAccount/create',
        method: 'post',
        data
    })
}

// @Tags OfficialAccountApi
// @Summary 获取单一公众号信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.OfficialAccount true "获取单一公众号信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /officialAccount/officialAccount [get]
export const getOfficialAccount = (params) => {
    return service({
        url: '/officialAccount/officialAccount',
        method: 'get',
        params
    })
}

// @Tags OfficialAccountApi
// @Summary 获取权限公众号列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限公众号列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /officialAccount/officialAccountList [get]
export const getOfficialAccountList = (params) => {
    return service({
        url: '/officialAccount/officialAccountList',
        method: 'get',
        params
    })
}
