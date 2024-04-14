package utils

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"go.uber.org/zap"
	"log"
	"net/url"
	"regexp"
	"strings"
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

func CollectArticle(keyWord string) {
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	collector.SetRequestTimeout(time.Second * 60)

	subCollector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	// 请求发起时回调,一般用来设置请求头等
	collector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Cookie", "BIDUPSID=C155065D7D880DB16532F2B5ACAAA29F; PSTM=1647945341; BAIDUID=BD4B73003EE2F83D2376C3B1D0F34EC0:FG=1; BD_UPN=123253; H_WISE_SIDS_BFESS=39662_40207_40212_40217_40079_40364_40351_40298_40369_40338_40416_40466_40460_40482_40317; H_PS_PSSID=40298_40369_40416_40482_40512_40397_60037_60031_60047_40511; BDSFRCVID=K4_OJexroG3Kj2QtQY6ghbsZNzn3CubTDYrE8HDLnjRzg2_VY-fcEG0PttO_kXtb6jBvogKK3gOTH4PF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF=tbC8VCD5JCt3H48k-4QEbbQH-UnLqbOdHmOZ04n-ah05OnF40JOGWjK7MM5eXf6EaJ7mLRom3UTKsq76Wh35K5tTQP6rLqFL5D74KKJxbPbc_nvy5-K-MxDBhUJiB5JMBan7_UJIXKohJh7FM4tW3J0ZyxomtfQxtNRJ0DnjtpChbC8lejuajTQQeU5eetjK2CntsJOOaCvobtbOy4oWK441DPRBJqoJLGKq0bnt-McFJRuRLnoD3M04K4o9-hvT-54e2p3FBUQZKJQnQft20b03XnnWXn3aJDQGLJ7jWhk5ep72y5OmQlRX5q79atTMfNTJ-qcH0KQpsIJM5-DWbT8IjHCeJ6F8tRktoIvX24THD6rP-trf5DCShUFsB53dB2Q-XPoO3KOvjCoh5hQZj4-J-lbzWJTXt27t3fbgy4opbhneX-T80RIsDRnpXCrhK2TxoUJ2bp7rjUnm-6j8LRKebPRiJ-b9Qg-JbpQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0HPonHjKBj63y3H; sugstore=0; H_WISE_SIDS=40298_40369_40416_40482_40512_40397_60037_60031_60047_40511; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; BA_HECTOR=8k00808ha5ah0la1048580almpqnb21j1lb3k1t; delPer=0; BD_CK_SAM=1; PSINO=6; BAIDUID_BFESS=BD4B73003EE2F83D2376C3B1D0F34EC0:FG=1; BDSFRCVID_BFESS=K4_OJexroG3Kj2QtQY6ghbsZNzn3CubTDYrE8HDLnjRzg2_VY-fcEG0PttO_kXtb6jBvogKK3gOTH4PF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF_BFESS=tbC8VCD5JCt3H48k-4QEbbQH-UnLqbOdHmOZ04n-ah05OnF40JOGWjK7MM5eXf6EaJ7mLRom3UTKsq76Wh35K5tTQP6rLqFL5D74KKJxbPbc_nvy5-K-MxDBhUJiB5JMBan7_UJIXKohJh7FM4tW3J0ZyxomtfQxtNRJ0DnjtpChbC8lejuajTQQeU5eetjK2CntsJOOaCvobtbOy4oWK441DPRBJqoJLGKq0bnt-McFJRuRLnoD3M04K4o9-hvT-54e2p3FBUQZKJQnQft20b03XnnWXn3aJDQGLJ7jWhk5ep72y5OmQlRX5q79atTMfNTJ-qcH0KQpsIJM5-DWbT8IjHCeJ6F8tRktoIvX24THD6rP-trf5DCShUFsB53dB2Q-XPoO3KOvjCoh5hQZj4-J-lbzWJTXt27t3fbgy4opbhneX-T80RIsDRnpXCrhK2TxoUJ2bp7rjUnm-6j8LRKebPRiJ-b9Qg-JbpQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0HPonHjKBj63y3H; ZFY=v1swEOIXPmAIGlcNHtoru3xd:BSQvZ:AiBuyM5Zk9pBKc:C; H_PS_645EC=4e6dTxYQjGttjrjxL4knxn83YP2pEtkUCsBqdsfiYzIBcAyGl2M%2FQHLntxY; BDUSS=RoeGtSRmtVYUQzTGw0dnF-OUFYeDFPMzZHeVY5aHp1eGlILU9lb0NrZ0RPMEptSUFBQUFBJCQAAAAAAAAAAAEAAABFPI80sPyyy77NysfBpsG~AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOuGmYDrhpmWV; BDUSS_BFESS=RoeGtSRmtVYUQzTGw0dnF-OUFYeDFPMzZHeVY5aHp1eGlILU9lb0NrZ0RPMEptSUFBQUFBJCQAAAAAAAAAAAEAAABFPI80sPyyy77NysfBpsG~AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOuGmYDrhpmWV; ab_sr=1.0.1_MzFlZWFlNjNkODlmMTVjZjg1YjZkZjJhMjRjNGZkZjA5OWQ2NThiYTk1YzQ4MGY3MDNlNjE4ZjAzYjAzYThkOWU0NzljZmIyYTU5YzRhNGFhMjc3MGFhOWUyOTUwNjhjNTFkOWYwZGE3ZmIxYWE1MzY5YTBmYjE1M2NhMTQ1OGM0MDMzYTM0MzI2ODUzNzY0NDRhNDM2MzQxMWQyNWFjNA==; BDRCVFR[feWj1Vr5u3D]=I67x6TjHwwYf0; BDSVRTM=6; WWW_ST=1713024561191")
		fmt.Println("----> 开始请求了")
	})

	// 请求发起时回调,一般用来设置请求头等
	subCollector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Cookie", "BIDUPSID=C155065D7D880DB16532F2B5ACAAA29F; PSTM=1647945341; BAIDUID=BD4B73003EE2F83D2376C3B1D0F34EC0:FG=1; H_WISE_SIDS_BFESS=39662_40207_40212_40217_40079_40364_40351_40298_40369_40338_40416_40466_40460_40482_40317; theme=bjh; PHPSESSID=j5plog2r8sfs61oenbt9609oo1; RECENT_LOGIN=1; gray=1; canary=0; Hm_lvt_f7b8c775c6c8b6a716a75df506fb72df=1712501212; H_PS_PSSID=40298_40369_40416_40482_40512_40397_60037_60031_60047_40511; H_WISE_SIDS=40298_40369_40416_40482_40512_40397_60037_60031_60047_40511; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; BDUSS=RoeGtSRmtVYUQzTGw0dnF-OUFYeDFPMzZHeVY5aHp1eGlILU9lb0NrZ0RPMEptSUFBQUFBJCQAAAAAAAAAAAEAAABFPI80sPyyy77NysfBpsG~AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOuGmYDrhpmWV; BDUSS_BFESS=RoeGtSRmtVYUQzTGw0dnF-OUFYeDFPMzZHeVY5aHp1eGlILU9lb0NrZ0RPMEptSUFBQUFBJCQAAAAAAAAAAAEAAABFPI80sPyyy77NysfBpsG~AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOuGmYDrhpmWV; BAIDUID_BFESS=BD4B73003EE2F83D2376C3B1D0F34EC0:FG=1; delPer=0; BA_HECTOR=2ga58la5ag05a18h048g2ga4tn1cdm1j1lcod1t; ZFY=v1swEOIXPmAIGlcNHtoru3xd:BSQvZ:AiBuyM5Zk9pBKc:C; PSINO=6; BCLID=10827640111856734396; BCLID_BFESS=10827640111856734396; BDSFRCVID=x58OJexroG3hrIrtyZF3hbsZNNLKRoRTDYLEOwXPsp3LGJLVY5DNEG0PtH691jtb6j3eogKK3gOTH4PF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; BDSFRCVID_BFESS=x58OJexroG3hrIrtyZF3hbsZNNLKRoRTDYLEOwXPsp3LGJLVY5DNEG0PtH691jtb6j3eogKK3gOTH4PF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF=tbC8VCDKJDD3H48k-4QEbbQH-UnLqhTf02OZ04n-ah05j-oVKROBW65LMM5eXf6wBg3NWMom3UTdsq76Wh35K5tTQP6rLqFt3go4KKJxbpP-bDjF5-K-M6t8hUJiB5JMBan7_UJIXKohJh7FM4tW3J0ZyxomtfQxtNRJ0DnjtpChbC8lejuaj65LeU5eetjK2CntsJOOaCv4ohbOy4oWK441DU7HBRjJMGLqXJnt-McF8lRo04nC3M04X-o9-hvT-54e2p3FBUQZMqb8Qft20b0njGuDaJ3aJDLe0R7jWhk5ep72y5OmQlRX5q79atTMfNTJ-qcH0KQpsIJM5-DWbT8IjHCeJ6F8tRFqoCvt-5rDHJTg5DTjhPrMWbKLWMT-MTryKKOS-qTKf-TLQPOs-RtlX4TrBhjyLHnRhlRNB-3iV-OxDUvnyxAZWHoRWfQxtNRJaILXtRcMHCQFDPoobUPUDUJ9LUkJLgcdot5yBbc8eIna5hjkbfJBQttjQn3hfIkj2CKLtKDBbDDCj6L3-RJH-xQ0KnLXKKOLVh-bXh7ketn4hUt5MbDNjlO8Xf5K2IrehqrYLb5pjMj2QhrdQf4WWb3ebTJr32Qr-fIa3n3psIJM5bQWbRvX0MKqajJ-aKviaKOjBMb1MMJDBT5h2M4qMxtOLR3pWDTm_q5TtUJMeCnTDMFhe6JLeHuDtjKDfKresJoq2RbhKROvhjRWjjkgyxoObtRxtKbHa4o2MxcPhnnPKp6FXxAU0xnNLU3kBgT9LMnx--t58h3_XhjZjf0NQttjQn37JJKDBtbtJJF5eJ7TyU45bU47yaOT0q4Hb6b9BJcjfU5MSlcNLTjpQT8r5MDOK5OuJRQ2QJ8BJD_MMDoP; H_BDCLCKID_SF_BFESS=tbC8VCDKJDD3H48k-4QEbbQH-UnLqhTf02OZ04n-ah05j-oVKROBW65LMM5eXf6wBg3NWMom3UTdsq76Wh35K5tTQP6rLqFt3go4KKJxbpP-bDjF5-K-M6t8hUJiB5JMBan7_UJIXKohJh7FM4tW3J0ZyxomtfQxtNRJ0DnjtpChbC8lejuaj65LeU5eetjK2CntsJOOaCv4ohbOy4oWK441DU7HBRjJMGLqXJnt-McF8lRo04nC3M04X-o9-hvT-54e2p3FBUQZMqb8Qft20b0njGuDaJ3aJDLe0R7jWhk5ep72y5OmQlRX5q79atTMfNTJ-qcH0KQpsIJM5-DWbT8IjHCeJ6F8tRFqoCvt-5rDHJTg5DTjhPrMWbKLWMT-MTryKKOS-qTKf-TLQPOs-RtlX4TrBhjyLHnRhlRNB-3iV-OxDUvnyxAZWHoRWfQxtNRJaILXtRcMHCQFDPoobUPUDUJ9LUkJLgcdot5yBbc8eIna5hjkbfJBQttjQn3hfIkj2CKLtKDBbDDCj6L3-RJH-xQ0KnLXKKOLVh-bXh7ketn4hUt5MbDNjlO8Xf5K2IrehqrYLb5pjMj2QhrdQf4WWb3ebTJr32Qr-fIa3n3psIJM5bQWbRvX0MKqajJ-aKviaKOjBMb1MMJDBT5h2M4qMxtOLR3pWDTm_q5TtUJMeCnTDMFhe6JLeHuDtjKDfKresJoq2RbhKROvhjRWjjkgyxoObtRxtKbHa4o2MxcPhnnPKp6FXxAU0xnNLU3kBgT9LMnx--t58h3_XhjZjf0NQttjQn37JJKDBtbtJJF5eJ7TyU45bU47yaOT0q4Hb6b9BJcjfU5MSlcNLTjpQT8r5MDOK5OuJRQ2QJ8BJD_MMDoP; Hm_lpvt_f7b8c775c6c8b6a716a75df506fb72df=1713026619; devStoken=cf820b506aa32b87815664e87b22d06d6e5f9d2291d7d2724c5338fc1da7c46e; bjhStoken=6d7dd9ea0e054f5a1178f8baff05d9fe6e5f9d2291d7d2724c5338fc1da7c46e; __bid_n=18ea48c2c6987368819ffc; RT=\"z=1&dm=baidu.com&si=bec4a7db-e43d-432c-a462-6884534de5b1&ss=luybc3dc&sl=7&tt=ash&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&ld=il7m&nu=ycg2a0z&cl=i4rh\"; ab_sr=1.0.1_MDkzMzUzYjljYWM0Y2UzOTc1MjdkZTZjZmE3ZWZmMDBkMDIxM2M0YzUzNmI3YjEyZDZjODVmNTI3NDNkNGM4ODYzZGQ0ODU2MTIxYTcxZTdhOTkzMjBlMTU2ZDQxNjVjMjExZGRjNWYyMjUwZGYxZWE1YWVlZWNkZWRlZTg2MmZkOTFiMTQzM2NiZTRmYmE3Mjk3OThjMjVlYzk0ZTYxOQ==")
		fmt.Println("----> 开始请求了")
	})

	// 请求完成后回调
	collector.OnResponse(func(response *colly.Response) {
		fmt.Println("----> 开始返回了")
	})

	// 请求完成后回调
	subCollector.OnResponse(func(response *colly.Response) {
		fmt.Println("----> 开始返回了")
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

		item := ai.Article{
			Title:       title,
			Comment:     content,
			AuthorName:  author,
			PublishTime: publishTime,
		}

		//item.BASEMODEL = ai2.BaseModel()

		// 将文章添加到结果切片中
		result = append(result, item)
	})

	encodedParam := url.QueryEscape(keyWord)

	err := collector.Visit("http://www.baidu.com/s?tn=news&rtt=1&bsst=1&cl=2&wd=" + encodedParam)
	if err != nil {
		fmt.Println(err)
	}

}

func CollectToutiaoArticle(link string) {
	collector := colly.NewCollector(
		colly.Async(true),
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

	subCollector := colly.NewCollector(
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

	// 请求发起时回调,一般用来设置请求头等
	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("----> 开始请求了")
	})

	// 请求发起时回调,一般用来设置请求头等
	subCollector.OnRequest(func(request *colly.Request) {
		fmt.Println("----> 开始请求了")
	})

	// 请求完成后回调
	collector.OnResponse(func(response *colly.Response) {
		fmt.Println("----> 开始返回了")
	})

	// 请求完成后回调
	subCollector.OnResponse(func(response *colly.Response) {
		fmt.Println("----> 开始返回了")
	})

	result := make([]ai.Article, 0)

	// 定义一个回调函数，处理页面响应
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if strings.Contains(url, "article") {
			time.Sleep(5 * time.Second)
			// 访问文章URL
			subCollector.Visit(url)
		}

	})

	//请求发生错误回调
	collector.OnError(func(response *colly.Response, err error) {
		fmt.Printf("发生错误了:%v", err)
	})

	// 提取标题
	subCollector.OnHTML("div.EaCvy", func(element *colly.HTMLElement) {
		fmt.Println(1)
		//title := element.ChildText("._3tNyU ._2oTsX .sKHSJ")
		//author := element.ChildText("._3tNyU ._2oTsX .bH7m7 ._2bKNC ._2gGWi")
		//publishTime := element.ChildText("._3tNyU ._2oTsX .bH7m7 ._2bKNC ._2sjh9")
		//content := element.ChildText("._18p7x")
		//
		//title = strings.TrimSpace(title) // 移除多余的空格
		//author = strings.TrimSpace(author)
		//publishTime = strings.TrimSpace(publishTime)
		//content = strings.TrimSpace(content)
		// 将文章添加到结果切片中
		//result = append(result, item)
	})

	err := collector.Visit(link)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(10 * time.Second)

	fmt.Println(Parse2Json(result))
}
