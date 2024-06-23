package utils

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/cast"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"os"
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

	link := "https://www.toutiao.com/trending/7379493903504654390/?category_name=topic_innerflow&=hot_board&log_pb=%7B%22category_name%22%3A%22topic_innerflow%22%2C%22cluster_type%22%3A%226%22%2C%22enter_from%22%3A%22click_category%22%2C%22entrance_hotspot%22%3A%22outside%22%2C%22event_type%22%3A%22hot_board%22%2C%22hot_board_cluster_id%22%3A%227379493903504654390%22%2C%22hot_board_impr_id%22%3A%22202406121940219BCB28E60A584F7F9F96%22%2C%22jump_page%22%3A%22hot_board_page%22%2C%22location%22%3A%22news_hot_card%22%2C%22page_location%22%3A%22hot_board_page%22%2C%22rank%22%3A%2227%22%2C%22source%22%3A%22trending_tab%22%2C%22style_id%22%3A%2240132%22%2C%22title%22%3A%22%E6%96%B0%E5%8A%A0%E5%9D%A1%E9%98%9F%E9%97%A8%E5%B0%86%E6%A1%91%E5%B0%BC%E7%9A%84%E9%A4%90%E5%8E%85%E7%88%86%E7%81%AB%22%7D&rank=27&style_id=40132&topic_id=7379493903504654390"
	err := collector.Visit(link)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

func TestCollectWechatArticle(t *testing.T) {
	// 设置WebDriver路径和浏览器选项
	const (
		// These paths will be different on your system.
		seleniumPath    = "vendor/selenium-server-standalone-3.4.jar"
		geckoDriverPath = "/usr/local/bin/chromedriver"
		port            = 8080
	)
	opts := []selenium.ServiceOption{
		//selenium.StartFrameBuffer(),            // Start an X frame buffer for the browser to run in.
		//selenium.ChromeDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr), // Output debug information to STDERR.
	}
	selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(geckoDriverPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}

	//禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	// chrome设置
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		//Path:  "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		Args: []string{
			//静默执行请求
			//"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			//"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36", // 模拟user-agent，防反爬,
			//"--disable-blink-features=AutomationControlled", // 从 Chrome 88 开始，它的 V8 引擎升级了，加了这个参数，window.navigator.webdriver=false
			//"--proxy-server=socks5://127.0.0.1:1080",
		},
		ExcludeSwitches: []string{
			"enable-automation", // 禁用左上角的控制显示
		},
	}
	caps.AddChrome(chromeCaps)

	// 设置Cookie
	cookie1 := &selenium.Cookie{
		Name:   "__ac_nonce",            // 替换为实际的cookie名称
		Value:  "0666ec66c0011c75926f1", // 替换为实际的cookie值
		Path:   "/",                     // 路径，通常为"/"
		Domain: "www.toutiao.com",       // 域名
	}

	cookie2 := &selenium.Cookie{
		Name:   "__ac_signature",                                                                                                                                      // 替换为实际的cookie名称
		Value:  "_02B4Z6wo00f01PZtaJwAAIDCS-JDd52bMPj2SWwAAFv7I7XYOlYSUGHkhGcNZuRCqZXuSEVWidh.UOF8ukFMy.COulbSInC.hCRBbMXBqR9cvEUCm7-iyZGFuKu5XUEncIP7d6LTui5DvYYhcb", // 替换为实际的cookie值
		Path:   "/",                                                                                                                                                   // 路径，通常为"/"
		Domain: "www.toutiao.com",                                                                                                                                     // 域名
	}

	cookie3 := &selenium.Cookie{
		Name:   "msToken",                                                                                                                          // 替换为实际的cookie名称
		Value:  "DJvJaotce0mVwfUUqX9zv5312fm7_8cQRnfNeaI-vbdLi2Myb-Y99F_RixakvBgYGhhBnHcR-HlPgznbGrbJutEwG1oABS5POZeO4qpJMQKG8g63n9hDonuYRjctNlw=", // 替换为实际的cookie值
		Path:   "/",                                                                                                                                // 路径，通常为"/"
		Domain: ".bytedance.com",                                                                                                                   // 域名
	}

	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer driver.Quit()

	driver.AddCookie(cookie1)
	driver.AddCookie(cookie2)
	driver.AddCookie(cookie3)
	// 导航到网页
	err = driver.Get("https://www.baidu.com/s?wd=%E5%A7%9C%E8%90%8D&rsv_spt=1&rsv_iqid=0xd7e81c74009f9ea9&issp=1&f=3&rsv_bp=1&rsv_idx=2&ie=utf-8&tn=baiduhome_pg&rsv_dl=ts_0&rsv_enter=1&rsv_sug3=4&rsv_sug1=3&rsv_sug7=100&rsv_sug2=0&rsv_btype=i&prefixsug=jia&rsp=0&inputT=1680&rsv_sug4=2408")
	if err != nil {
		log.Fatalf("Failed to navigate to page: %v", err)
	}

	time.Sleep(5 * time.Second)

	// 获取网页的HTML内容
	htmlSource, err := driver.PageSource()
	if err != nil {
		log.Fatalf("Failed to get page source: %v", err)
	}

	// 输出网页HTML内容
	fmt.Println(htmlSource)
}
