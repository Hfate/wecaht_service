import service from '@/utils/request'


// @Tags ArticleApi
// @Summary 删除文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Article true "删除文章"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /dailyArticle/dailyArticle [delete]
export const deleteArticle = (data) => {
    return service({
        url: '/dailyArticle/dailyArticle',
        method: 'delete',
        data
    })
}


// @Tags ArticleApi
// @Summary 改写文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Article true "改写文章"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /dailyArticle/recreation [delete]
export const recreationArticle = (data) => {
    return service({
        url: '/dailyArticle/recreation',
        method: 'post',
        data
    })
}

// @Tags ArticleApi
// @Summary 删除文章集合
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Article true "删除文章集合"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /dailyArticle/deleteArticlesByIds [delete]
export const deleteArticlesByIds = (data) => {
    return service({
        url: '/dailyArticle/deleteArticlesByIds',
        method: 'delete',
        data
    })
}

// @Tags ArticleApi
// @Summary 获取单一文章信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Article true "获取单一文章信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /dailyArticle/dailyArticle [get]
export const getArticle = (params) => {
    return service({
        url: '/dailyArticle/dailyArticle',
        method: 'get',
        params
    })
}

// @Tags ArticleApi
// @Summary 获取权限文章列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限文章列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /dailyArticle/articleList [get]
export const getArticleList = (params) => {
    return service({
        url: '/dailyArticle/dailyArticleList',
        method: 'get',
        params
    })
}



