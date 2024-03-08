package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/wechat"

type BenchmarkAccountResponse struct {
	BenchmarkAccount ai.BenchmarkAccount `json:"benchmarkAccount"`
}
