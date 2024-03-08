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
