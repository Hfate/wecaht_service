import service from '@/utils/request'
// @Tags OfficialAccountApi
// @Summary 创建门户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.OfficialAccount true "创建门户"
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
// @Summary 更新门户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.OfficialAccount true "更新门户信息"
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
// @Summary 删除门户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.OfficialAccount true "删除门户"
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
// @Summary 获取单一门户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.OfficialAccount true "获取单一门户信息"
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
// @Summary 获取权限门户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限门户列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /officialAccount/officialAccountList [get]
export const getOfficialAccountList = (params) => {
    return service({
        url: '/officialAccount/officialAccountList',
        method: 'get',
        params
    })
}
