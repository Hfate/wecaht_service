import service from '@/utils/request'
// @Tags MediaApi
// @Summary 创建素材
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Media true "创建素材"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /media/media [post]
export const createMedia = (data) => {
    return service({
        url: '/media/media',
        method: 'post',
        data
    })
}

// @Tags MediaApi
// @Summary 更新素材信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Media true "更新素材信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /media/media [put]
export const updateMedia = (data) => {
    return service({
        url: '/media/media',
        method: 'put',
        data
    })
}

// @Tags MediaApi
// @Summary 删除素材
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Media true "删除素材"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /media/media [delete]
export const deleteMedia = (data) => {
    return service({
        url: '/media/media',
        method: 'delete',
        data
    })
}


// @Tags MediaApi
// @Summary 获取素材集合
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取素材集合"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /media/mediaList [get]
export const getMediaList = (params) => {
    return service({
        url: '/media/mediaList',
        method: 'get',
        params
    })
}


// @Tags MediaApi
// @Summary 获取素材分页
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取素材分页"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /media/mediaList [get]
export const getMediaPage = (params) => {
    return service({
        url: '/media/mediaPage',
        method: 'get',
        params
    })
}
