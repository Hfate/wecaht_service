import service from '@/utils/request'


// @Tags PortalApi
// @Summary 删除热点
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Portal true "删除热点"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /hotspot/hotspot [delete]
export const deleteHotspot = (data) => {
    return service({
        url: '/hotspot/hotspot',
        method: 'delete',
        data
    })
}


// @Tags PortalApi
// @Summary 删除热点
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Portal true "删除热点"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /hotspot/create [delete]
export const hotspotCreate = (data) => {
    return service({
        url: '/hotspot/create',
        method: 'post',
        data
    })
}


// @Tags PortalApi
// @Summary 获取热点列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取热点列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/portalList [get]
export const getHotspotList = (params) => {
    return service({
        url: '/hotspot/hotspotList',
        method: 'get',
        params
    })
}
