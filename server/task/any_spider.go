package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	ai2 "github.com/flipped-aurora/gin-vue-admin/server/service/ai"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
	"time"
)

func AnySpider() {
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	subCollector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)
	collector.SetRequestTimeout(time.Second * 60)
	subCollector.SetRequestTimeout(time.Second * 60)

	collector.OnHTML("div.cl", func(e *colly.HTMLElement) {
		articleUrl := e.ChildAttr("a", "href")
		// 访问文章URL
		subCollector.Visit("http://m.yiyuzheng.com.cn" + articleUrl)
	})

	subCollector.OnHTML("div.page_con", func(element *colly.HTMLElement) {
		title := element.ChildText(".titlem")
		author := element.ChildText("._2gGWi")
		publishTime := element.ChildText(".fr")
		content := element.ChildText(".descp")
		//
		title = strings.TrimSpace(title) // 移除多余的空格
		author = strings.TrimSpace(author)
		publishTime = strings.TrimSpace(publishTime)
		publishTime = strings.ReplaceAll(publishTime, "发布时间：", "")
		content = strings.TrimSpace(content)

		item := ai.Article{
			Topic:       "抑郁症",
			Title:       title,
			Comment:     content,
			AuthorName:  author,
			PublishTime: publishTime,
		}

		item.BASEMODEL = ai2.BaseModel()
		global.GVA_DB.Model(&ai.Article{}).Create(item)
	})

	err := collector.Visit("http://m.yiyuzheng.com.cn/")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
}
