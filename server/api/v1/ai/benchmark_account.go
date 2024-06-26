package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	aiRes "github.com/flipped-aurora/gin-vue-admin/server/model/ai/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type BenchmarkAccountApi struct{}

// CreateBenchmarkAccount
// @Tags      Portal
// @Summary   创建对标账号
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.BenchmarkAccount            true  "门户用户名, 门户手机号码"
// @Success   200   {object}  response.Response{msg=string}  "创建对标账号"
// @Router    /benchmark/benchmark [post]
func (e *BenchmarkAccountApi) CreateBenchmarkAccount(c *gin.Context) {
	var benchmarkAccount ai.BenchmarkAccount
	err := c.ShouldBindJSON(&benchmarkAccount)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(benchmarkAccount, utils.BenchmarkAccountVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = benchmarkAccountService.CreateBenchmarkAccount(benchmarkAccount)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateBenchmarkAccount
// @Tags      Portal
// @Summary   更新对标账号
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.BenchmarkAccount            true  ""
// @Success   200   {object}  response.Response{msg=string}  "更新对标账号"
// @Router    /benchmark/benchmark [put]
func (e *BenchmarkAccountApi) UpdateBenchmarkAccount(c *gin.Context) {
	var benchmarkAccount ai.BenchmarkAccount
	err := c.ShouldBindJSON(&benchmarkAccount)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(benchmarkAccount, utils.BenchmarkAccountVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = benchmarkAccountService.UpdateBenchmarkAccount(benchmarkAccount)
	if err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败", c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

// UpdateWxToken
// @Tags      Portal
// @Summary   创建对标账号
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.BenchmarkAccount            true  "门户用户名, 门户手机号码"
// @Success   200   {object}  response.Response{msg=string}  "更新wx token"
// @Router    /benchmark/benchmark [post]
func (e *BenchmarkAccountApi) UpdateWxToken(c *gin.Context) {
	var wxToken ai.WxToken
	err := c.ShouldBindJSON(&wxToken)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(wxToken, utils.WxTokenVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = wxTokenService.UpdateWxToken(wxToken)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteBenchmarkAccount
// @Tags      BenchmarkAccount
// @Summary   删除对标账号
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.BenchmarkAccount            true  "对标账号ID"
// @Success   200   {object}  response.Response{msg=string}  "删除对标账号"
// @Router    /benchmark/benchmark [delete]
func (e *BenchmarkAccountApi) DeleteBenchmarkAccount(c *gin.Context) {
	var benchmarkAccount ai.BenchmarkAccount
	err := c.ShouldBindJSON(&benchmarkAccount)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(benchmarkAccount.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = benchmarkAccountService.DeleteBenchmarkAccount(benchmarkAccount)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetBenchmarkAccount
// @Tags      BenchmarkAccount
// @Summary   获取单一对标账号信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     wechat.BenchmarkAccount                                                true  "对标账号ID"
// @Success   200   {object}  response.Response{data=exampleRes.BenchmarkAccountResponse,msg=string}  "获取单一对标账号信息,返回包括对标账号详情"
// @Router    /benchmarkAccount/benchmarkAccount [get]
func (e *BenchmarkAccountApi) GetBenchmarkAccount(c *gin.Context) {
	var benchmarkAccount ai.BenchmarkAccount
	err := c.ShouldBindQuery(&benchmarkAccount)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(benchmarkAccount.BASEMODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := benchmarkAccountService.GetBenchmarkAccount(cast.ToUint64(benchmarkAccount.ID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(aiRes.BenchmarkAccountResponse{BenchmarkAccount: data}, "获取成功", c)
}

// GetBenchmarkAccountList
// @Tags      BenchmarkAccount
// @Summary   分页获取权限对标账号列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限对标账号列表,返回包括列表,总数,页码,每页数量"
// @Router    /benchmarkAccount/benchmarkAccountList [get]
func (e *BenchmarkAccountApi) GetBenchmarkAccountList(c *gin.Context) {
	var pageInfo aiReq.BenchmarkAccountSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	benchmarkAccountList, total, err := benchmarkAccountService.GetBenchmarkAccountList(utils.GetUserAuthorityId(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     benchmarkAccountList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
