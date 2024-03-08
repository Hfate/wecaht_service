package task

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	ai2 "github.com/flipped-aurora/gin-vue-admin/server/service/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/storyicon/graphquery"
	"gorm.io/gorm"
	"sync"
)

func HotspotSpider(db *gorm.DB) error {
	// 微博热点爬取
	spiderWeiboHeadline(db)

	portalList := make([]*ai.Portal, 0)
	err := db.Model(ai.Portal{}).Where("target_num>0").Find(&portalList).Error
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for _, portal := range portalList {
		wg.Add(1)
		p := portal
		go func() {
			defer wg.Done()
			spiderHeadlines(db, p)
		}()
	}

	wg.Wait()

	return nil
}

func spiderHeadlines(db *gorm.DB, portal *ai.Portal) {
	hotspotList := collectHeadlinesInfo(portal.Link, portal.GraphQuery)
	if len(hotspotList) == 0 {
		return
	}
	urlList := make([]string, 0)
	for _, hotspot := range hotspotList {
		urlList = append(urlList, hotspot.Link)
		hotspot.BASEMODEL = ai2.BaseModel()
		hotspot.PortalName = portal.PortalName
	}

	err := db.Where("link in ?", urlList).Unscoped().Delete(&ai.Hotspot{}).Error
	if err != nil {
		return
	}

	err = db.Create(&hotspotList).Error

	fmt.Println(err)
}

func collectHeadlinesInfo(url, gQuery string) []*ai.Hotspot {
	htmlResult, err := utils.GetStr(url)
	if err != nil {
		return []*ai.Hotspot{}
	}

	response := graphquery.ParseFromString(htmlResult, gQuery)

	hotspotList := make([]*ai.Hotspot, 0)
	response.Decode(&hotspotList)

	return hotspotList
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

	err = db.Where("link in ?", urlList).Unscoped().Delete(&ai.Hotspot{}).Error
	if err != nil {
		return
	}

	err = db.Create(&hotspotList).Error
	if err != nil {
		fmt.Println(err)
	}
	return
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
