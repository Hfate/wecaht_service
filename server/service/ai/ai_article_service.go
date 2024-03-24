package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
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

// @function: ApprovalArticle
// @description: 发布文章
// @param: e model.AIArticle
// @return: err error

func (exa *AIArticleService) PublishArticle(aiArticle ai.AIArticle) (err error) {
	officialAccount, err := OfficialAccountServiceApp.GetOfficialAccountByAppId(aiArticle.TargetAccountId)
	if err != nil {
		return err
	}

	// 发布文章
	publishId, mediaID, msgId, msgDataID, err := WechatServiceApp.PublishArticle(officialAccount, aiArticle)
	if err != nil {
		return err
	}

	// 更新发布状态
	aiArticle.TargetAccountName = officialAccount.AccountName
	aiArticle.PublishTime = time.Now()
	aiArticle.ArticleStatus = 3
	aiArticle.MediaId = mediaID
	aiArticle.MsgId = msgId
	aiArticle.PublishId = publishId
	aiArticle.MsgDataID = msgDataID
	err = global.GVA_DB.Save(&aiArticle).Error
	return err
}

func (exa *AIArticleService) CreateAIArticle(e ai.AIArticle) (err error) {
	err = global.GVA_DB.Where("origin_id=?", e.OriginId).Delete(&ai.AIArticle{}).Error

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

func (exa *AIArticleService) Recreation(id uint64) (err error) {
	article := ai.AIArticle{}
	err = global.GVA_DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		return err
	}
	return ArticleServiceApp.Recreation(article.OriginId)
}

// @function: GetArticleList
// @description: 分页获取文章列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *AIArticleService) GetAIArticleList(sysUserAuthorityID uint, info aiReq.AIArticleSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

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
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&articleList).Error
	}
	return articleList, total, err
}
