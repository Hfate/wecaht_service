package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiRes "github.com/flipped-aurora/gin-vue-admin/server/model/ai/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timeutil"
	"github.com/spf13/cast"
	"sort"
	"strings"
)

type DailyArticleService struct {
}

var DailyArticleServiceApp = new(DailyArticleService)

//@function: DeleteFileChunk
//@description: 删除文章
//@param: e model.DailyArticle
//@return: err error

func (exa *DailyArticleService) DeleteDailyArticle(e ai.DailyArticle) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@function: DeleteDailyArticlesByIds
//@description: 删除选中文章
//@param: ids []wechat.DailyArticle
//@return: err error

func (exa *DailyArticleService) DeleteDailyArticlesByIds(ids request.IdsReq) (err error) {
	var DailyArticles []ai.DailyArticle
	err = global.GVA_DB.Find(&DailyArticles, "id in ?", ids.Ids).Delete(&DailyArticles).Error
	return err
}

//@function: GetDailyArticle
//@description: 获取文章信息
//@param: id uint
//@return: customer model.DailyArticle, err error

func (exa *DailyArticleService) GetDailyArticle(id uint64) (DailyArticle ai.DailyArticle, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&DailyArticle).Error
	return
}

func (exa *DailyArticleService) Recreation(id uint64) error {
	DailyArticle := ai.DailyArticle{}
	err := global.GVA_DB.Where("id = ?", id).First(&DailyArticle).Error
	if err != nil {
		return err
	}

	context := &ArticleContext{
		Topic: DailyArticle.Topic,
		Link:  DailyArticle.Link,
		Title: DailyArticle.Title,
	}

	chatGptResp, err := ChatModelServiceApp.Recreation(context)
	if err != nil {
		return err
	}

	aiDailyArticle := ai.AIArticle{
		OriginId:   cast.ToUint64(DailyArticle.ID),
		Title:      chatGptResp.Title,
		PortalName: DailyArticle.PortalName,
		Topic:      chatGptResp.Topic,
		AuthorName: DailyArticle.AuthorName,
		Tags:       strings.Join(chatGptResp.Tags, ","),
		Content:    chatGptResp.Content,
	}
	aiDailyArticle.BASEMODEL = BaseModel()

	err = AIArticleServiceApp.CreateAIArticle(aiDailyArticle)

	if err != nil {
		return err
	}
	DailyArticle.UseTimes++
	err = global.GVA_DB.Save(&DailyArticle).Error

	return err
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetDailyArticleList
// @description: 分页获取文章列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *DailyArticleService) GetDailyArticleList(sysUserAuthorityID uint) (list interface{}, total int64, err error) {
	var a system.SysAuthority
	a.AuthorityId = sysUserAuthorityID
	_, err = systemService.AuthorityServiceApp.GetAuthorityInfo(a)
	if err != nil {
		return
	}

	var articleList []ai.DailyArticle

	// 获取所有的公众号
	var officialAccountList []ai.OfficialAccount
	err = global.GVA_DB.Find(&officialAccountList).Error
	if err != nil {
		return
	}

	for _, account := range officialAccountList {
		batchId := timeutil.GetCurDate() + account.AppId

		var dailyArticleList []ai.DailyArticle
		err = global.GVA_DB.Where("batch_id = ?", batchId).Find(&dailyArticleList).Error
		if err != nil {
			return
		}
		articleList = append(articleList, dailyArticleList...)

		targetNum := account.TargetNum

		diffNum := targetNum - len(dailyArticleList)

		if diffNum > 0 {
			for i := 0; i < diffNum; i++ {
				// 找到主题相关的文章素材
				article, err := ArticleServiceApp.FindHotArticleByTopic(account.Topic)
				if err != nil {
					continue
				}

				aiDailyArticle := ai.DailyArticle{
					Title:             article.Title,
					PortalName:        article.PortalName,
					Topic:             article.Topic,
					AuthorName:        article.AuthorName,
					Tags:              article.Tags,
					Content:           article.Content,
					BatchId:           batchId,
					ReadNum:           article.ReadNum,
					LikeNum:           article.LikeNum,
					CommentNum:        article.CommentNum,
					Link:              article.Link,
					UseTimes:          0,
					HotspotId:         article.HotspotId,
					TargetAccountId:   account.AppId,
					TargetAccountName: account.AccountName,
				}

				article.UseTimes = article.UseTimes + 1

				err = global.GVA_DB.Save(article).Error
				if err != nil {
					continue
				}

				aiDailyArticle.BASEMODEL = BaseModel()
				err = global.GVA_DB.Model(&ai.DailyArticle{}).Create(&aiDailyArticle).Error
				if err != nil {
					continue
				}

				articleList = append(articleList, aiDailyArticle)
			}

		}
	}

	articleMap := make(map[string][]ai.DailyArticle)

	batchIdList := make([]string, 0)
	for _, item := range articleList {
		articleMap[item.BatchId] = append(articleMap[item.BatchId], item)
		batchIdList = append(batchIdList, item.BatchId)
	}

	batchIdList = utils.RemoveRepByMap(batchIdList)

	// 排序 batchIdList
	sort.Strings(batchIdList)

	result := make([]aiRes.DailyArticleParentResponse, 0)
	for _, batchId := range batchIdList {
		subList := articleMap[batchId]

		children := make([]ai.DailyArticle, 0)
		if len(subList) > 1 {
			children = subList[1:]
		}
		item := aiRes.DailyArticleParentResponse{
			DailyArticle: subList[0],
			Children:     children,
		}
		result = append(result, item)

	}

	return result, int64(len(result)), err
}
