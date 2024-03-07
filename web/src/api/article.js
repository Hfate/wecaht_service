import service from '@/utils/request'


// @Tags ArticleApi
// @Summary 删除文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Article true "删除文章"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /article/article [delete]
export const deleteArticle = (data) => {
    return service({
        url: '/article/article',
        method: 'delete',
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
// @Router /article/deleteArticlesByIds [delete]
export const deleteArticlesByIds = (data) => {
    return service({
        url: '/article/deleteArticlesByIds',
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
// @Router /article/article [get]
export const getArticle = (params) => {
    return service({
        url: '/article/article',
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
// @Router /article/articleList [get]
export const getArticleList = (params) => {
    return service({
        url: '/article/articleList',
        method: 'get',
        params
    })
}


// @Tags ArticleApi
// @Summary 获取文章导入模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取文章导入模板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /article/template [get]
export const getTemplate = (params) => {
    return service({
        url: '/article/template',
        method: 'get',
        params
    })
}