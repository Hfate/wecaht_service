package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	aiRes "github.com/flipped-aurora/gin-vue-admin/server/model/ai/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timeutil"
	"go.uber.org/zap"
	"sort"
	"time"
)

type AIArticleService struct {
}

var AIArticleServiceApp = new(AIArticleService)

//@function: DeleteFileChunk
//@description: 删除文章
//@param: e model.Article
//@return: err error

func (exa *AIArticleService) DeleteAIArticle(e ai.AIArticle) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@function: DeleteArticlesByIds
//@description: 删除选中文章
//@param: ids []wechat.Article
//@return: err error

func (exa *AIArticleService) DeleteAIArticlesByIds(ids request.IdsReq) (err error) {
	var articles []ai.AIArticle
	err = global.GVA_DB.Find(&articles, "id in ?", ids.Ids).Delete(&articles).Error
	return err
}

func (exa *AIArticleService) PublishArticle(ids []string) error {
	var articles []ai.AIArticle
	err := global.GVA_DB.Model(&ai.AIArticle{}).Where("id in ?", ids).Find(&articles).Error

	batchAIArticleMap := make(map[string][]ai.AIArticle)
	for _, item := range articles {
		batchAIArticleMap[item.BatchId] = append(batchAIArticleMap[item.BatchId], item)
	}

	go func() {
		for _, list := range batchAIArticleMap {
			err = exa.Publish1Article(list)
			if err != nil {

				for _, item := range list {

					global.GVA_LOG.Error("发布失败!", zap.Error(err), zap.String("item", item.Title))
					item.ArticleStatus = 4
					item.ErrMessage = err.Error()
					global.GVA_DB.Save(&item)

				}

			}

		}
	}()

	return err
}

func (exa *AIArticleService) UpdateArticle(aiArticle ai.AIArticle) (err error) {
	// 更新
	err = global.GVA_DB.Save(&aiArticle).Error
	return err
}

// @function: ApprovalArticle
// @description: 发布文章
// @param: e model.AIArticle
// @return: err error

func (exa *AIArticleService) Publish1Article(aiArticleList []ai.AIArticle) (err error) {
	if len(aiArticleList) == 0 {
		return err
	}
	officialAccount, err := OfficialAccountServiceApp.GetOfficialAccountByAppId(aiArticleList[0].TargetAccountId)
	if err != nil {
		return err
	}

	// 发布文章
	publishId, mediaID, msgId, msgDataID, err := WechatServiceApp.PublishArticle(officialAccount, aiArticleList)
	if err != nil {
		return err
	}

	for _, aiArticle := range aiArticleList {
		// 更新发布状态
		aiArticle.TargetAccountName = officialAccount.AccountName
		aiArticle.PublishTime = time.Now()
		aiArticle.ArticleStatus = 1
		aiArticle.MediaId = mediaID
		aiArticle.MsgId = msgId
		aiArticle.PublishId = publishId
		aiArticle.MsgDataID = msgDataID
		aiArticle.ErrMessage = ""
		err = global.GVA_DB.Save(&aiArticle).Error
	}

	return err
}

func (exa *AIArticleService) CreateAIArticle(e ai.AIArticle) (err error) {
	e.PublishTime = time.Now()
	err = global.GVA_DB.Create(&e).Error
	return err
}

//@function: GetAIArticle
//@description: 获取文章信息
//@param: id uint
//@return: customer model.Article, err error

func (exa *AIArticleService) GetAIArticle(id uint64) (aiArticle ai.AIArticle, err error) {
	err = global.GVA_DB.Where("id = ?", id).Last(&aiArticle).Error
	return
}

func (exa *AIArticleService) GenerateArticleById(hotspotId uint64, account *ai.OfficialAccount) error {
	context := &ArticleContext{
		Topic:     account.Topic,
		Account:   account,
		HotspotId: hotspotId,
	}

	ArticlePipelineApp.Run("hotspot", context)

	return nil
}

func (exa *AIArticleService) GenerateArticle(account *ai.OfficialAccount) error {
	targetNum := account.TargetNum
	if targetNum == 0 {
		return nil
	}

	batchId := timeutil.GetCurDate() + account.AppId
	// 重置当天的素材池，将use time 更新为0
	err := global.GVA_DB.Model(&ai.DailyArticle{}).Where("batch_id = ?", batchId).Update("use_times", 0).Error
	if err != nil {
		return err
	}

	// 删除当天生成的文章
	err = global.GVA_DB.Where("batch_id = ?", batchId).Delete(&ai.AIArticle{}).Error
	if err != nil {
		return err
	}

	i := 0
	for i < targetNum && i < 10 {
		context := &ArticleContext{
			Topic:       account.Topic,
			Account:     account,
			Params:      []string{},
			CreateTypes: account.CreateTypes,
		}

		err = ArticlePipelineApp.Run("", context)
		if err != nil {
			return err
		}

		i++
	}

	return nil
}

// GenerateDailyArticle 生成每日文章
func (exa *AIArticleService) GenerateDailyArticle() error {
	// 获取公众号列表
	list, err := OfficialAccountServiceApp.List()

	if err != nil {
		return err
	}

	for _, account := range list {

		if account.AppId == "" {
			continue
		}

		item := account
		time.Sleep(5 * time.Second)

		err2 := exa.GenerateArticle(item)
		if err2 != nil {
			global.GVA_LOG.Error("GenerateArticle With err", zap.Error(err2))
		}

	}

	return nil
}

func (exa *AIArticleService) Recreation(id uint64) (err error) {
	aiArticle := ai.AIArticle{}
	err = global.GVA_DB.Where("id = ?", id).First(&aiArticle).Error
	if err != nil {
		return err
	}

	officeAccount, err := OfficialAccountServiceApp.GetOfficialAccountByAppId(aiArticle.TargetAccountId)
	if err != nil {
		return err
	}

	err = exa.GenerateArticle(officeAccount)

	return
}

// @function: GetArticleList
// @description: 分页获取文章列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *AIArticleService) GetAIArticleList(sysUserAuthorityID uint, info aiReq.AIArticleSearch) (list interface{}, total int64, err error) {

	batchId := timeutil.GetCurDate()

	var a system.SysAuthority
	a.AuthorityId = sysUserAuthorityID
	auth, err := systemService.AuthorityServiceApp.GetAuthorityInfo(a)
	if err != nil {
		return
	}
	var dataId []uint
	for _, v := range auth.DataAuthorityId {
		dataId = append(dataId, v.AuthorityId)
	}
	var articleList []ai.AIArticle

	db := global.GVA_DB.Model(&ai.AIArticle{})

	db = db.Where("batch_id LIKE ?", "%"+batchId+"%")

	// 如果有条件搜索 下方会自动创建搜索语句
	if info.PortalName != "" {
		db = db.Where("portal_name LIKE ?", "%"+info.PortalName+"%")
	}

	if info.TargetAccountName != "" {
		db = db.Where("target_account_name LIKE ?", "%"+info.TargetAccountName+"%")
	}

	if info.Topic != "" {
		db = db.Where("topic LIKE ?", "%"+info.Topic+"%")
	}
	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return articleList, total, err
	} else {
		err = db.Order("batch_id desc,target_account_id desc").Find(&articleList).Error
	}

	articleMap := make(map[string][]ai.AIArticle)

	batchIdList := make([]string, 0)
	for _, item := range articleList {
		articleMap[item.BatchId] = append(articleMap[item.BatchId], item)
		batchIdList = append(batchIdList, item.BatchId)
	}

	batchIdList = utils.RemoveRepByMap(batchIdList)

	// 排序 batchIdList
	sort.Strings(batchIdList)

	result := make([]aiRes.AIArticleParentResponse, 0)
	for _, batchId := range batchIdList {
		subList := articleMap[batchId]

		children := make([]ai.AIArticle, 0)
		if len(subList) > 1 {
			children = subList[1:]
		}
		item := aiRes.AIArticleParentResponse{
			AIArticle: subList[0],
			Children:  children,
		}
		result = append(result, item)

	}

	return result, total, err
}
