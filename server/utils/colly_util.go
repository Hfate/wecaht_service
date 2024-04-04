package utils

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"go.uber.org/zap"
	"log"
	"net/url"
	"regexp"
	"time"
)

type Item struct {
	Id     string
	Width  int
	Height int
	Links  Links
	Urls   Urls
}

type Links struct {
	Download string
}

type Urls struct {
	Regular string
}

func CollectUnsplashImgUrl(keyWord string) []string {
	c := colly.NewCollector(
		colly.Async(true),
	)
	err := c.Limit(&colly.LimitRule{
		DomainRegexp: `unsplash\.com`,
		RandomDelay:  500 * time.Millisecond,
		Parallelism:  12,
	})
	if err != nil {
		global.GVA_LOG.Error("SearchUnsplash", zap.Error(err))
	}

	result := make([]string, 0)
	c.OnResponse(func(r *colly.Response) {
		var items []*Item
		json.Unmarshal(r.Body, &items)
		for _, item := range items {
			result = append(result, item.Urls.Regular)
		}
	})

	encodedParam := url.QueryEscape(keyWord)

	c.Visit(fmt.Sprintf("https://unsplash.com/napi/photos?page=1&per_page=10&query=%s", encodedParam))

	c.Wait()

	return result
}

func CollectBaiduImgUrl(keyWord string) []string {
	collector := colly.NewCollector(
		func(collector *colly.Collector) {
			// 设置随机ua
			extensions.RandomUserAgent(collector)
		},
		func(collector *colly.Collector) {
			collector.OnRequest(func(request *colly.Request) {
				log.Println(request.URL, ", User-Agent:", request.Headers.Get("User-Agent"))
			})
		},
	)

	collector.SetRequestTimeout(time.Second * 60)

	// 请求发起时回调,一般用来设置请求头等
	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("----> 开始请求了")
	})

	result := make([]string, 0)
	// 请求完成后回调
	collector.OnResponse(func(response *colly.Response) {
		fmt.Println("----> 开始返回了")
		resp := string(response.Body)
		// 正则表达式匹配objURL
		re := regexp.MustCompile(`"objURL":"([^"]+)"`)
		matches := re.FindAllStringSubmatch(resp, -1)
		//打印所有找到的objURL
		for _, match := range matches {
			result = append(result, match[1])
		}

	})

	//请求发生错误回调
	collector.OnError(func(response *colly.Response, err error) {
		fmt.Printf("发生错误了:%v", err)
	})

	err := collector.Visit("http://image.baidu.com/search/flip?tn=baiduimage&ie=utf-8&ct=201326592&v=flip&word=" + keyWord)
	if err != nil {
		fmt.Println(err)
	}

	return result
}
