package task

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	ai2 "github.com/flipped-aurora/gin-vue-admin/server/service/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

func HotspotSpider(db *gorm.DB) error {
	// 微博热点爬取
	spiderWeiboHeadline(db)

	// 头条热点爬取
	spiderToutiaoHeadline(db)

	spiderBaiduHeadline(db)

	//db.Exec("UPDATE hotspot SET avg_speed = CASE  WHEN TIMESTAMPDIFF(MINUTE, created_at, updated_at) = 0 THEN 0\n                    ELSE trending / NULLIF(TIMESTAMPDIFF(MINUTE, created_at, updated_at), 0)\n    END\nWHERE trending IS NOT NULL AND updated_at IS NOT NULL AND created_at IS NOT NULL;")

	return nil
}

func spiderBaiduHeadline(db *gorm.DB) {
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

	hotspotList := make([]*ai.Hotspot, 0)
	linkUrlList := make([]string, 0)

	collector.OnHTML(".container-bg_lQ801", func(element *colly.HTMLElement) {
		element.ForEach(".category-wrap_iQLoo", func(i int, element *colly.HTMLElement) {
			aLink := element.DOM.ChildrenFiltered("a")
			jumpLink, _ := aLink.Attr("href")
			title := element.ChildText(".content_1YWBm .c-single-text-ellipsis")
			trending := element.ChildText(".trend_2RttY .hot-index_1Bl1a")
			hotspotList = append(hotspotList, &ai.Hotspot{
				PortalName: "百度",
				BASEMODEL:  ai2.BaseModel(),
				Link:       jumpLink,
				Headline:   title,
				Trending:   cast.ToInt(trending),
			})
			linkUrlList = append(linkUrlList, jumpLink)
		})
	})

	err := collector.Visit("https://top.baidu.com/board?tab=realtime")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	if len(linkUrlList) == 0 {
		return
	}

	err = ai2.HotspotServiceImp.CreateHotspot(hotspotList)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func spiderWeiboHeadline(db *gorm.DB) {
	resp, err := utils.GetStr("https://weibo.com/ajax/side/hotSearch")
	if err != nil {
		fmt.Println(err)
		return
	}
	var webResp WebResp
	err = utils.JsonStrToStruct(resp, &webResp)
	if err != nil {
		fmt.Println(err)
		return
	}

	hotspotList := make([]*ai.Hotspot, 0)
	urlList := make([]string, 0)
	for _, realTime := range webResp.Data.Realtime {
		hotspot := &ai.Hotspot{
			PortalName: "微博",
			BASEMODEL:  ai2.BaseModel(),
			Link:       "https://s.weibo.com/weibo?q=" + realTime.Word,
			Headline:   realTime.Word,
			Trending:   realTime.Num,
			Topic:      realTime.Category,
		}
		urlList = append(urlList, hotspot.Link)
		hotspotList = append(hotspotList, hotspot)
	}

	err = ai2.HotspotServiceImp.CreateHotspot(hotspotList)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func spiderToutiaoHeadline(db *gorm.DB) {
	resp, err := utils.GetStr("https://www.toutiao.com/hot-event/hot-board/?origin=toutiao_pc")
	if err != nil {
		fmt.Println(err)
		return
	}
	var toutiaoResp ToutiaoResp
	err = utils.JsonStrToStruct(resp, &toutiaoResp)
	if err != nil {
		fmt.Println(err)
		return
	}

	hotspotList := make([]*ai.Hotspot, 0)
	urlList := make([]string, 0)
	headlineList := make([]string, 0)
	for _, item := range toutiaoResp.Data {
		topic := ""
		interestCategoryList := item.InterestCategory
		if len(interestCategoryList) > 0 {
			for _, i := range interestCategoryList {
				topic += TopicTranslateMap[i] + ","
			}
		} else {
			topic = item.LabelDesc
		}

		topic = strings.TrimRight(topic, ",")

		hotspot := &ai.Hotspot{
			PortalName: "头条",
			BASEMODEL:  ai2.BaseModel(),
			Link:       item.Url,
			Headline:   item.Title,
			Trending:   cast.ToInt(item.HotValue),
			Topic:      topic,
		}
		urlList = append(urlList, hotspot.Link)
		headlineList = append(headlineList, hotspot.Headline)
		hotspotList = append(hotspotList, hotspot)
	}

	err = ai2.HotspotServiceImp.CreateHotspot(hotspotList)
	if err != nil {
		fmt.Println(err)
	}
	return
}

var TopicTranslateMap = map[string]string{
	"sports":        "体育",
	"entertainment": "娱乐",
	"business":      "商业",
	"technology":    "科技",
	"finance":       "财经",
	"health":        "健康",
	"travel":        "旅游",
	"music":         "音乐",
	"lifestyle":     "生活",
	"food":          "美食",
	"car":           "汽车",
	"science":       "科学",
	"education":     "教育",
	"game":          "游戏",
	"anime":         "动漫",
	"comic":         "漫画",
	"news":          "新闻",
	"military":      "军事",
	"international": "时事",
	"other":         "其他",
	"stock":         "财经",
	"estate":        "财经",
	"taiwan":        "时事",
	"tourism":       "旅游",
	"culture":       "文化",
}

type ToutiaoResp struct {
	Data []struct {
		ClusterId int64  `json:"ClusterId"`
		Title     string `json:"Title"`
		LabelUrl  string `json:"LabelUrl"`
		Label     string `json:"Label"`
		Url       string `json:"Url"`
		HotValue  string `json:"HotValue"`
		Schema    string `json:"Schema"`
		LabelUri  struct {
			Uri     string `json:"uri"`
			Url     string `json:"url"`
			Width   int    `json:"width"`
			Height  int    `json:"height"`
			UrlList []struct {
				Url string `json:"url"`
			} `json:"url_list"`
			ImageType int `json:"image_type"`
		} `json:"LabelUri"`
		ClusterIdStr string `json:"ClusterIdStr"`
		ClusterType  int    `json:"ClusterType"`
		QueryWord    string `json:"QueryWord"`
		Image        struct {
			Uri     string `json:"uri"`
			Url     string `json:"url"`
			Width   int    `json:"width"`
			Height  int    `json:"height"`
			UrlList []struct {
				Url string `json:"url"`
			} `json:"url_list"`
			ImageType int `json:"image_type"`
		} `json:"Image"`
		LabelDesc        string   `json:"LabelDesc,omitempty"`
		InterestCategory []string `json:"InterestCategory,omitempty"`
	} `json:"data"`
	FixedTopData []struct {
		Id     int    `json:"Id"`
		Title  string `json:"Title"`
		Url    string `json:"Url"`
		Schema string `json:"Schema"`
	} `json:"fixed_top_data"`
	FixedTopStyle string `json:"fixed_top_style"`
	ImprId        string `json:"impr_id"`
	Status        string `json:"status"`
}

type WebResp struct {
	Ok   int `json:"ok"`
	Data struct {
		Realtime []struct {
			Realpos            int         `json:"realpos"`
			AdInfo             string      `json:"ad_info"`
			StarWord           int         `json:"star_word"`
			FunWord            int         `json:"fun_word"`
			Emoticon           string      `json:"emoticon"`
			Word               string      `json:"word"`
			WordScheme         string      `json:"word_scheme"`
			StarName           interface{} `json:"star_name"`
			TopicFlag          int         `json:"topic_flag"`
			Mid                string      `json:"mid"`
			Flag               int         `json:"flag"`
			LabelName          string      `json:"label_name"`
			IconDescColor      string      `json:"icon_desc_color,omitempty"`
			IconDesc           string      `json:"icon_desc,omitempty"`
			Note               string      `json:"note"`
			OnboardTime        int         `json:"onboard_time"`
			RawHot             int         `json:"raw_hot"`
			SmallIconDesc      string      `json:"small_icon_desc,omitempty"`
			SmallIconDescColor string      `json:"small_icon_desc_color,omitempty"`
			Category           string      `json:"category"`
			IsHot              int         `json:"is_hot,omitempty"`
			Num                int         `json:"num"`
			Expand             int         `json:"expand"`
			ChannelType        string      `json:"channel_type"`
			SubjectQuerys      string      `json:"subject_querys"`
			Extension          int         `json:"extension"`
			SubjectLabel       string      `json:"subject_label"`
			Rank               int         `json:"rank"`
			IsWarm             int         `json:"is_warm,omitempty"`
			FlagDesc           string      `json:"flag_desc,omitempty"`
			IsNew              int         `json:"is_new,omitempty"`
		} `json:"realtime"`
		Hotgovs []struct {
			Pos                int    `json:"pos"`
			Note               string `json:"note"`
			IconDescColor      string `json:"icon_desc_color"`
			IsGov              int    `json:"is_gov"`
			Word               string `json:"word"`
			SmallIconDescColor string `json:"small_icon_desc_color"`
			SmallIconDesc      string `json:"small_icon_desc"`
			IsHot              int    `json:"is_hot"`
			IconDesc           string `json:"icon_desc"`
			Url                string `json:"url"`
			Name               string `json:"name"`
			TopicFlag          int    `json:"topic_flag"`
			Mid                string `json:"mid"`
			Flag               int    `json:"flag"`
		} `json:"hotgovs"`
		Hotgov struct {
			Pos                int    `json:"pos"`
			Note               string `json:"note"`
			IconDescColor      string `json:"icon_desc_color"`
			IsGov              int    `json:"is_gov"`
			Word               string `json:"word"`
			SmallIconDescColor string `json:"small_icon_desc_color"`
			SmallIconDesc      string `json:"small_icon_desc"`
			IsHot              int    `json:"is_hot"`
			IconDesc           string `json:"icon_desc"`
			Url                string `json:"url"`
			Name               string `json:"name"`
			TopicFlag          int    `json:"topic_flag"`
			Mid                string `json:"mid"`
			Flag               int    `json:"flag"`
		} `json:"hotgov"`
	} `json:"data"`
}
