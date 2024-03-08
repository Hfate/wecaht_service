package task

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/spf13/cast"
	"github.com/storyicon/graphquery"
	"gorm.io/gorm"
	"net/url"
	"strings"
	"sync"
	"time"
)

func PortalSpider(db *gorm.DB) error {
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
			spiderPortal(db, p)
		}()
	}

	wg.Wait()

	return nil
}

func spiderPortal(db *gorm.DB, portal *ai.Portal) {

	urlList := collectAllUrl(portal.Link, portal.PortalKey, portal.TargetNum)
	notInDbUrls := findNotInDb(db, urlList)
	size := cast.ToString(len(notInDbUrls))
	fmt.Println(portal.PortalName + " collect url size " + size)

	urlSet := make(map[string]bool)
	collectSize := 0
	for index, articleUrl := range notInDbUrls {
		if !strings.Contains(articleUrl, portal.ArticleKey) {
			continue
		}

		if _, ok := urlSet[articleUrl]; ok {
			continue
		}
		urlSet[articleUrl] = true

		fmt.Println(articleUrl)
		fmt.Println(portal.PortalName + "[" + cast.ToString(index) + "/" + size + "]")

		time.Sleep(10 * time.Millisecond)

		articleResp, er := spiderArticle(portal.GraphQuery, articleUrl)
		if er != nil {
			continue
		}
		if articleResp.Content == "" {
			continue
		}

		articleResp.ReadNum = strings.ReplaceAll(articleResp.ReadNum, "阅读", "")

		article := &ai.Article{
			BASEMODEL:   global.BASEMODEL{ID: utils.GenID(), CreatedAt: time.Now(), UpdatedAt: time.Now()},
			PortalName:  portal.PortalName,
			AuthorName:  articleResp.AuthorName,
			Topic:       articleResp.Topic,
			Title:       articleResp.Title,
			Link:        articleUrl,
			Content:     articleResp.Content,
			PublishTime: articleResp.PublishTime,
			ReadNum:     cast.ToInt(articleResp.ReadNum),
			CommentNum:  cast.ToInt(articleResp.CommentNum),
			LikeNum:     cast.ToInt(articleResp.LikeNum),
		}

		err := db.Model(&ai.Article{}).Create(article).Error
		if err != nil {
			fmt.Println(err)
		}

		err = db.Model(&ai.Url{}).Create(&ai.Url{
			Url:       articleUrl,
			BASEMODEL: global.BASEMODEL{ID: utils.GenID(), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}).Error

		fmt.Println(portal.PortalName + "[collectSize" + ":" + cast.ToString(collectSize) + "]")

		collectSize++
	}
}

func findNotInDb(db *gorm.DB, urlList []string) []string {
	dbUrls := make([]string, 0)
	err := db.Model(&ai.Url{}).Where("url in ?", urlList).Select("url").Find(&dbUrls).Error
	if err != nil {
		return urlList
	}

	dbUrlSet := make(map[string]bool)
	for _, item := range dbUrls {
		dbUrlSet[item] = true
	}

	result := make([]string, 0)

	for _, item := range urlList {
		if _, ok := dbUrlSet[item]; !ok {
			result = append(result, item)
		}
	}
	return result
}

func collectAllUrl(portalUrl string, portalKey string, targetNum int) []string {
	urlList := []string{portalUrl}
	result := make([]string, 0)
	urlSet := make(map[string]bool)

	// 爬到目标数量网页 或者 遍历了100遍 责退出循环
	loopNum := 0
	for len(result) < targetNum && loopNum < 50 {
		subList := collectUrlList(portalKey, urlList)

		time.Sleep(10 * time.Millisecond)

		newSubList := make([]string, 0)
		for _, item := range subList {
			if !strings.Contains(item, portalKey) {
				item = portalUrl + item
			}

			if !isValidUrl(item) {
				continue
			}

			if _, ok := urlSet[item]; !ok {
				result = append(result, item)
				urlSet[item] = true
				newSubList = append(newSubList, item)
			}

		}
		loopNum++
		urlList = newSubList

		fmt.Println("loopNum->" + cast.ToString(loopNum) + ",resultNum->" + cast.ToString(len(result)))
	}
	return result
}

func collectUrlList(portalKey string, urlList []string) []string {
	result := make([]string, 0)
	urlSet := make(map[string]bool)
	for _, u := range urlList {
		if !strings.Contains(u, portalKey) {
			continue
		}
		subList := collectUrl(u)
		for _, item := range subList {
			if _, ok := urlSet[item]; !ok {
				result = append(result, item)
				urlSet[item] = true
			}
		}
	}
	return result
}

func collectUrl(url string) []string {
	htmlResult, err := utils.GetStr(url)
	if err != nil {
		return []string{}
	}
	query := "a `css(\"a\")` [{  url  `attr(\"href\")` }]"

	response := graphquery.ParseFromString(htmlResult, query)

	urls := make([]*URL, 0)
	response.Decode(&urls)

	result := make([]string, 0)
	for _, item := range urls {
		result = append(result, item.URL)
	}

	return result
}

func spiderArticle(expression string, url string) (*ArticleResp, error) {
	htmlResult, err := utils.GetStr(url)
	if err != nil {
		return nil, err
	}
	//readNum = https://www.yicai.com/api/ajax/getnewsdetail?id=102009612
	response := graphquery.ParseFromString(htmlResult, expression)
	article := &ArticleResp{}
	err = response.Decode(&article)
	return article, err
}

func isValidUrl(str string) bool {
	count := strings.Count(str, "http")
	if count > 1 {
		return false
	}
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Hostname() == "" {
		return false
	}
	return true
}

type URL struct {
	URL string `json:"url"`
}

type ArticleResp struct {
	PortalName  string `json:"portalName"`
	Title       string `json:"title"`       //
	Topic       string `json:"topic"`       //
	AuthorName  string `json:"authorName" ` //
	Link        string `json:"link"`        //
	PublishTime string `json:"publishTime"`
	LikeNum     string `json:"likeNum"`
	ReadNum     string `json:"readNum"`
	CommentNum  string `json:"commentNum"`
	Content     string `json:"content"`
}
