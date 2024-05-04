package ai

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"strings"
)

type SiTongService struct {
}

var SiTongServiceApp = new(SiTongService)

func (*SiTongService) Similarity(req *SimilarityReq) (*SimilarityResp, error) {
	header := make(map[string]string)
	header["secret-id"] = global.GVA_CONFIG.Sitong.AccessKey
	header["secret-key"] = global.GVA_CONFIG.Sitong.SecretKey

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(req.Compare[0]))
	if err == nil {
		req.Compare[0] = dom.Text()
	}

	_, resp, err := utils.PostWithHeaders(global.GVA_CONFIG.Sitong.ApiUrl, utils.Parse2Json(req), header)
	if err != nil {
		return nil, err
	}

	var similarityResp SimilarityResp
	err = utils.JsonStrToStruct(string(resp), &similarityResp)
	if err != nil {
		global.GVA_LOG.Error("Similarity", zap.Error(err))
		return nil, err
	}

	return &similarityResp, nil
}

type SimilarityReq struct {
	Text    string   `json:"text"`
	Compare []string `json:"compare"`
}

type SimilarityResp struct {
	Msg     string `json:"msg"`
	Results []struct {
		Similarity float64 `json:"similarity"`
	} `json:"results"`
	Status string `json:"status"`
}
