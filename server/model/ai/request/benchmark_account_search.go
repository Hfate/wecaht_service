package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type BenchmarkAccountSearch struct {
	ai.BenchmarkAccount
	request.PageInfo
}
