package task

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	ai2 "github.com/flipped-aurora/gin-vue-admin/server/service/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timeutil"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"time"
)

// SpiderHotArticle

func SpiderHotArticle() {
	hotspotList := make([]ai.Hotspot, 0)
	twoHourAgo := timeutil.AddHours(timeutil.GetCurTime(), -10)

	err := global.GVA_DB.Where("created_at > ?", twoHourAgo).Where("spider_num=0").Find(&hotspotList).Error

	if err != nil {
		global.GVA_LOG.Error("SpiderHotArticle", zap.Error(err))
		return
	}

	for _, item := range hotspotList {
		err = spiderHotspot(item)
		global.GVA_LOG.Info("spiderHotspot", zap.String("hotspot", utils.Parse2Json(item)), zap.Error(err))
		time.Sleep(3 * time.Second)
	}

}

func spiderHotspot(hotspot ai.Hotspot) error {

	articleList := collectArticle(hotspot)
	if len(articleList) > 0 {
		err := global.GVA_DB.Create(articleList).Error
		hotspot.SpiderNum = len(articleList)
		global.GVA_DB.Save(hotspot)

		global.GVA_LOG.Info("spiderHotspot", zap.String("hotspot", utils.Parse2Json(hotspot)))
		return err
	}

	return nil
}

func collectArticle(hotspot ai.Hotspot) []ai.Article {
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	collector.SetRequestTimeout(time.Second * 60)

	subCollector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	// 请求发起时回调,一般用来设置请求头等
	collector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Cookie", global.GVA_CONFIG.QianFan.Cookie)
		fmt.Println("----> 开始请求了")
	})

	// 请求发起时回调,一般用来设置请求头等
	subCollector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Cookie", global.GVA_CONFIG.QianFan.Cookie)
		fmt.Println(request.URL.Path + "----> 开始请求了")
	})

	// 请求完成后回调
	collector.OnResponse(func(response *colly.Response) {
		fmt.Println("----> 开始返回了")
	})

	// 请求完成后回调
	subCollector.OnResponse(func(response *colly.Response) {
		fmt.Println(response.Request.URL.Path + "----> 开始返回了")
	})

	result := make([]ai.Article, 0)

	collectNum := 0

	// 定义一个回调函数，处理页面响应
	collector.OnHTML("h3", func(e *colly.HTMLElement) {
		articleUrl := e.ChildAttr("a", "href")

		if strings.Contains(articleUrl, "baijiahao") && collectNum <= 3 {
			// 解析URL
			parsedURL, err := url.Parse(articleUrl)
			if err != nil {
				fmt.Println("Error parsing URL:", err)
				return
			}

			// 更改URL的协议为http
			parsedURL.Scheme = "http"

			// 解析查询参数
			queryParams := parsedURL.Query()

			// 删除特定的查询参数wfr
			queryParams.Del("wfr")

			// 更新URL的查询参数
			parsedURL.RawQuery = queryParams.Encode()

			// 访问文章URL
			subCollector.Visit(parsedURL.String())

			time.Sleep(3 * time.Second)

			collectNum++
		}

	})

	//请求发生错误回调
	collector.OnError(func(response *colly.Response, err error) {
		fmt.Printf("发生错误了:%v", err)
	})

	// 提取标题
	subCollector.OnHTML("div.EaCvy", func(element *colly.HTMLElement) {

		title := element.ChildText(".sKHSJ")
		author := element.ChildText("._2gGWi")
		publishTime := element.ChildText("._2sjh9")
		content := element.ChildText("._18p7x")
		//
		title = strings.TrimSpace(title) // 移除多余的空格
		author = strings.TrimSpace(author)
		publishTime = strings.TrimSpace(publishTime)
		content = strings.TrimSpace(content)

		topic := hotspot.Topic
		if topic == "" {
			topic = "时事"
		}

		item := ai.Article{
			Title:       title,
			Comment:     content,
			AuthorName:  author,
			PublishTime: publishTime,
			HotspotId:   cast.ToUint64(hotspot.ID),
			Topic:       topic,
			PortalName:  "百家号",
		}

		publishTimeInt, _ := timeutil.StrToTimeStamp(publishTime, "2006-01-02 15:04:05")
		// 发布时间需大于今年
		if publishTimeInt < timeutil.GetYearStartTime(int64(time.Now().Year())) {
			return
		}

		item.BASEMODEL = ai2.BaseModel()

		// 将文章添加到结果切片中
		result = append(result, item)
	})

	encodedParam := url.QueryEscape(hotspot.Headline)

	err := collector.Visit("http://www.baidu.com/s?tn=news&rtt=1&bsst=1&cl=2&wd=" + encodedParam)
	if err != nil {
		global.GVA_LOG.Error("collectArticle", zap.Error(err))
	}

	return result
}
