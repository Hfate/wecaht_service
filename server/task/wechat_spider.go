package task

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/storyicon/graphquery"
	"gorm.io/gorm"
)

var GRAPH_QUERY = "{content `xpath(\"//*[@id='js_content']\")`}"

func WechatSpider(db *gorm.DB) error {
	benchmarkList := make([]ai.BenchmarkAccount, 0)
	err := db.Model(ai.BenchmarkAccount{}).Where("init_num>0").Find(&benchmarkList).Error
	if err != nil {
		return err
	}

	wxToken := &ai.WxToken{}
	er := global.GVA_DB.Model(&ai.WxToken{}).Where("1=1").Last(&wxToken).Error
	if er != nil {
		fmt.Println(er)
		return err
	}

	for _, item := range benchmarkList {
		articleList := service.ServiceGroupApp.AIServiceGroup.BenchmarkAccountService.SpiderOfficialAccount(wxToken, item)

		if len(articleList) > 0 {
			for _, a := range articleList {
				var count int64
				db.Model(&ai.Url{}).Where("url = ?", a.Link).Count(&count)

				//if count > 0 {
				//	continue
				//}

				resp, _ := utils.GetStr(a.Link)
				var content *WechatContent
				response := graphquery.ParseFromString(resp, GRAPH_QUERY)
				response.Decode(&content)

				a.Content = content.Content

				db.Model(&ai.Article{}).Create(&a)
			}
		}
	}

	return nil
}

type WechatContent struct {
	Content string `json:"content"`
}
