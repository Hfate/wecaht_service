import service from '@/utils/request'


// @Tags SettlementApi
// @Summary 获取权限公众号列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限公众号列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /Settlement/SettlementList [get]
export const getSettlementList = (params) => {
    return service({
        url: '/settlement/list',
        method: 'get',
        params
    })
}
