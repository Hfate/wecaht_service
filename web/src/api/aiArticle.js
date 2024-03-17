import service from '@/utils/request'


// @Tags AIArticleApi
// @Summary 删除文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.AIArticle true "删除文章"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /aiArticle/aiArticle [delete]
export const deleteAIArticle = (data) => {
    return service({
        url: '/aiArticle/aiArticle',
        method: 'delete',
        data
    })
}


// @Tags AIArticleApi
// @Summary 改写文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.AIArticle true "改写文章"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /aiArticle/recreation [delete]
export const recreationAIArticle = (data) => {
    return service({
        url: '/aiArticle/recreation',
        method: 'post',
        data
    })
}

// @Tags AIArticleApi
// @Summary 删除文章集合
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.AIArticle true "删除文章集合"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /aiArticle/deleteAIArticlesByIds [delete]
export const deleteAIArticlesByIds = (data) => {
    return service({
        url: '/aiArticle/deleteAIArticlesByIds',
        method: 'delete',
        data
    })
}

// @Tags AIArticleApi
// @Summary 获取单一文章信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.AIArticle true "获取单一文章信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /aiArticle/aiArticle [get]
export const getAIArticle = (params) => {
    return service({
        url: '/aiArticle/aiArticle',
        method: 'get',
        params
    })
}

// @Tags AIArticleApi
// @Summary 获取权限文章列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限文章列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /aiArticle/aiArticleList [get]
export const getAIArticleList = (params) => {
    return service({
        url: '/aiArticle/aiArticleList',
        method: 'get',
        params
    })
}


// @Tags AIArticleApi
// @Summary 获取权限文章列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限文章列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /aiArticle/aiArticleList [get]
export const publishArticle = (data) => {
    return service({
        url: '/aiArticle/publish',
        method: 'post',
        data
    })
}
