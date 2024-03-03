import service from '@/utils/request'
// @Tags PortalApi
// @Summary 创建门户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Portal true "创建门户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/portal [post]
export const createPortal = (data) => {
    return service({
        url: '/portal/portal',
        method: 'post',
        data
    })
}

// @Tags PortalApi
// @Summary 更新门户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Portal true "更新门户信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/portal [put]
export const updatePortal = (data) => {
    return service({
        url: '/portal/portal',
        method: 'put',
        data
    })
}

// @Tags PortalApi
// @Summary 删除门户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Portal true "删除门户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/portal [delete]
export const deletePortal = (data) => {
    return service({
        url: '/portal/portal',
        method: 'delete',
        data
    })
}

// @Tags PortalApi
// @Summary 获取单一门户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Portal true "获取单一门户信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/portal [get]
export const getPortal = (params) => {
    return service({
        url: '/portal/portal',
        method: 'get',
        params
    })
}

// @Tags PortalApi
// @Summary 获取权限门户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限门户列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/portalList [get]
export const getPortalList = (params) => {
    return service({
        url: '/portal/portalList',
        method: 'get',
        params
    })
}
