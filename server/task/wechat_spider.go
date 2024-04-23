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
	err := db.Model(ai.BenchmarkAccount{}).Where("init_num>0").Order("id desc").Find(&benchmarkList).Error
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

				if count > 0 {
					continue
				}

				resp, _ := utils.GetStr(a.Link)
				var content *WechatContent
				response := graphquery.ParseFromString(resp, GRAPH_QUERY)
				response.Decode(&content)

				a.Topic = item.Topic
				a.Content = content.Content

				db.Model(&ai.Article{}).Create(&a)
			}

		}
	}

	// 更新对标账号已爬取数目
	db.Exec("update wechat_benchmark_account a inner join\n    (select portal_name, count(*) totalNum from wechat_article group by portal_name) b on a.account_name = b.portal_name\nset a.finish_num = b.totalNum\nwhere 1=1;")

	return nil
}

type WechatContent struct {
	Content string `json:"content"`
}
