package utils

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/cast"
	"testing"
	"time"
)

func TestCollectArticle(t *testing.T) {
	curTime := time.Now()
	time.Sleep(2 * time.Second)

	fmt.Println(curTime.Sub(time.Now()).Milliseconds())

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	// 请求发起时回调,一般用来设置请求头等
	collector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Cookie", "__ac_nonce=0665c7753007e1cb7a608; _ga_QEHZPBE5HH=GS1.1.1717331872.5.1.1717335891.0.0.0; __ac_signature=_02B4Z6wo00f01RGuJyQAAIDDrCEMzB4yQDkRiiOAACI6ygmk--UIXy5rYsqsj77ZsKNZ1BUBXNjXmFRHWTmse5fovWgCM2ZCcr7uLGHU11lH73Ic9X7qRlSym7LE36oEZb12cjvPbtdyGBKFd1; __ac_referer=__ac_blank; msToken=CmnfAvSofZoACxNZrC_KTPsiwKryd4S0Yq-ProVEXgDhKVI4Bc31IuLl34lxYGu_VO_XkrNAYQtKNmaMzfrVGvDQc2kYbk42cobWHsMN; ttwid=1%7CjJhV-ZQBbp-rRPiTpQX02ojxAnJzbXA3xow5e5Q7iTA%7C1717335892%7Cb3e97916cc5efc25a7706234d6b4810ae649450704a42f890f87091b690f756d")
	})

	collector.OnResponse(func(response *colly.Response) {
		fmt.Println("----> 开始返回了" + cast.ToString(response.StatusCode))
	})

	result := make([]string, 0)
	collector.OnHTML("div.feed-card-article-l", func(element *colly.HTMLElement) {
		aLink := element.DOM.ChildrenFiltered("a")
		jumpLink, _ := aLink.Attr("href")
		fmt.Println(jumpLink)

		result = append(result, jumpLink)

	})

	link := "https://www.toutiao.com/trending/7379493903504654390/?category_name=topic_innerflow&event_type=hot_board&log_pb=%7B%22category_name%22%3A%22topic_innerflow%22%2C%22cluster_type%22%3A%226%22%2C%22enter_from%22%3A%22click_category%22%2C%22entrance_hotspot%22%3A%22outside%22%2C%22event_type%22%3A%22hot_board%22%2C%22hot_board_cluster_id%22%3A%227379493903504654390%22%2C%22hot_board_impr_id%22%3A%22202406121940219BCB28E60A584F7F9F96%22%2C%22jump_page%22%3A%22hot_board_page%22%2C%22location%22%3A%22news_hot_card%22%2C%22page_location%22%3A%22hot_board_page%22%2C%22rank%22%3A%2227%22%2C%22source%22%3A%22trending_tab%22%2C%22style_id%22%3A%2240132%22%2C%22title%22%3A%22%E6%96%B0%E5%8A%A0%E5%9D%A1%E9%98%9F%E9%97%A8%E5%B0%86%E6%A1%91%E5%B0%BC%E7%9A%84%E9%A4%90%E5%8E%85%E7%88%86%E7%81%AB%22%7D&rank=27&style_id=40132&topic_id=7379493903504654390"
	err := collector.Visit(link)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

func TestCollectWechatArticle(t *testing.T) {
	CollectWechatArticle("https://mp.weixin.qq.com/s?__biz=MzUyNzE4OTE1Mw==&mid=2247756365&idx=1&sn=abd56014b9839b7d434f4af12dd4df9b&chksm=fb0e06df07aab431bc5afd2f17a1d027efd09fa770cab2f515bac138d00a59181398a4909bbd&scene=0&xtrack=1#rd")
}
