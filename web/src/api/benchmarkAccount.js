import service from '@/utils/request'


// @Tags BenchmarkAccountApi
// @Summary 删除对标账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.BenchmarkAccount true "删除对标账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /benchmarkAccount/benchmarkAccount [delete]
export const deleteBenchmark = (data) => {
    return service({
        url: '/benchmark/benchmark',
        method: 'delete',
        data
    })
}


// @Tags BenchmarkAccountApi
// @Summary 更新对标账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.BenchmarkAccount true "更新对标账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /benchmarkAccount/benchmarkAccount [delete]
export const updateBenchmark = (data) => {
    return service({
        url: '/benchmark/benchmark',
        method: 'put',
        data
    })
}

// @Tags BenchmarkAccountApi
// @Summary 创建对标账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Portal true "创建对标账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/portal [post]
export const createBenchmark = (data) => {
    return service({
        url: '/benchmark/benchmark',
        method: 'post',
        data
    })
}


// @Tags BenchmarkAccountApi
// @Summary 更新WX Token
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Portal true "更新WX Token"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新WX Token"}"
// @Router /wxToken/wxToken [post]
export const updateWxToken = (data) => {
    return service({
        url: '/wxToken/wxToken',
        method: 'put',
        data
    })
}

// @Tags BenchmarkAccountApi
// @Summary 获取权限对标账号列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限对标账号列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /benchmarkAccount/benchmarkAccountList [get]
export const getBenchmarkAccountList = (params) => {
    return service({
        url: '/benchmark/benchmarkList',
        method: 'get',
        params
    })
}


// @Tags BenchmarkAccountApi
// @Summary 获取权限对标账号列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取权限对标账号列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /benchmarkAccount/benchmarkAccountList [get]
export const getBenchmarkAccount = (params) => {
    return service({
        url: '/benchmark/benchmark',
        method: 'get',
        params
    })
}