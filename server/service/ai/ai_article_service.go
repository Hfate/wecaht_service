package ai

import (
	"errors"
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
	"strings"
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
	aiArticle.PublishTime = time.Now()
	aiArticle.Content = exa.parseContent(aiArticle.Content)
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

func (exa *AIArticleService) GenerateArticleById(hotspotId uint64, account ai.OfficialAccount) error {
	context := &ArticleContext{
		Topic:     account.Topic,
		AppId:     account.AppId,
		HotspotId: hotspotId,
	}

	articleContext := ArticlePipelineApp.Run("hotspot", context)
	if articleContext.Content == "" {
		return errors.New("AI创作失败")
	}

	aiArticle := ai.AIArticle{
		Title:             exa.parseTitle(articleContext.Title),
		TargetAccountName: account.AccountName,
		TargetAccountId:   account.AppId,
		Topic:             articleContext.Topic,
		AuthorName:        account.DefaultAuthorName,
		Tags:              strings.Join(articleContext.Tags, ","),
		OriginalContent:   articleContext.Content,
		Content:           exa.parseContent(articleContext.Content),
	}
	aiArticle.BASEMODEL = BaseModel()

	err := AIArticleServiceApp.CreateAIArticle(aiArticle)

	global.GVA_LOG.Info("AI创作完成", zap.String("appId", account.AppId), zap.String("topic", account.Topic))

	return err
}

func (exa *AIArticleService) GenerateArticle(account *ai.OfficialAccount) error {
	targetNum := account.TargetNum

	i := 0

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

	for i < targetNum {
		context := &ArticleContext{
			Topic:       account.Topic,
			AppId:       account.AppId,
			Params:      []string{},
			CreateTypes: account.CreateTypes,
		}

		articleContext := ArticlePipelineApp.Run("", context)
		if articleContext.Content == "" || len(articleContext.Params) == 0 {
			return errors.New("AI创作失败")
		}

		aiArticle := ai.AIArticle{
			BatchId:           batchId,
			Title:             exa.parseTitle(articleContext.Title),
			TargetAccountName: account.AccountName,
			TargetAccountId:   account.AppId,
			Topic:             account.Topic,
			AuthorName:        account.DefaultAuthorName,
			Link:              articleContext.Link,
			Tags:              strings.Join(articleContext.Tags, ","),
			OriginalContent:   articleContext.Content,
			Content:           articleContext.Content,
			Params:            strings.Join(articleContext.Params, ","),
		}
		aiArticle.BASEMODEL = BaseModel()

		// 获取历史已发布消息5条图文消息
		publishArticleList, err2 := WechatServiceApp.BatchGetHistoryArticleList(account)
		if err2 != nil || len(publishArticleList.Item) == 0 {
			global.GVA_LOG.Error("BatchGetHistoryArticleList", zap.Error(err2))
		} else {
			originalContent := articleContext.Content
			originalContent += "---\n"
			originalContent += "#### 推荐阅读\n"
			for _, item := range publishArticleList.Item {

				if len(item.Content.NewsItem) == 0 {
					continue
				}

				originalContent += "-[" + item.Content.NewsItem[0].Title + "](" + item.Content.NewsItem[0].URL + ")\n"
			}
		}

		//  处理排版
		aiArticle.Content = exa.parseContent(articleContext.Content)

		err = AIArticleServiceApp.CreateAIArticle(aiArticle)
		if err != nil {
			global.GVA_LOG.Error("CreateAIArticle", zap.Error(err))
			continue
		}

		global.GVA_LOG.Info("AI创作完成", zap.String("appId", account.AppId), zap.String("topic", account.Topic))

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

		time.Sleep(500 * time.Millisecond)

		go func() {

			err2 := exa.GenerateArticle(item)
			if err2 != nil {
				global.GVA_LOG.Error("GenerateArticle With err", zap.Error(err2))
			}

		}()

	}

	return nil
}

func (exa *AIArticleService) parseTitle(title string) string {
	title = strings.ReplaceAll(title, "#", "")
	title = strings.ReplaceAll(title, "*", "")
	title = strings.ReplaceAll(title, "标题：", "")
	title = strings.ReplaceAll(title, "#", "")
	title = utils.RemoveBookTitleBrackets(title)
	title = strings.ReplaceAll(title, "标题建议：", "")
	return title
}

func (exa *AIArticleService) parseContent(content string) string {
	// 以换行符为分隔符，将文章内容拆分成多行
	lines := strings.Split(content, "\n")

	// 排除标题行
	var contentLines []string
	for _, line := range lines {
		if !strings.Contains(line, "标题：") &&
			!strings.Contains(line, "占位符") &&
			!strings.Contains(line, "配图") {
			contentLines = append(contentLines, line)
		}
	}

	// 将剩余的行重新连接成一篇文章
	markdownContent := strings.Join(contentLines, "\n")

	//```markdown
	markdownContent = strings.ReplaceAll(markdownContent, "```markdown", "")
	markdownContent = strings.ReplaceAll(markdownContent, "```", "")
	markdownContent = strings.ReplaceAll(markdownContent, "<li><p>", "<li>")

	htmlContent, _ := utils.RenderMarkdownContent(markdownContent)

	htmlContent = utils.HtmlAddStyle(htmlContent)

	return htmlContent
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
		err = db.Limit(limit).Offset(offset).Order("batch_id desc,target_account_id desc").Find(&articleList).Error
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
