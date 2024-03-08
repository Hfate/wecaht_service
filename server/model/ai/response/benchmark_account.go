package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/ai"

type BenchmarkAccountResponse struct {
	BenchmarkAccount ai.BenchmarkAccount `json:"benchmarkAccount"`
}
