package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
)

type BenchmarkAccountSearch struct {
	wechat.BenchmarkAccount
	request.PageInfo
}
