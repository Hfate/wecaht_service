package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
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

func (exa *AIArticleService) PublishArticle(ids []int) error {
	var articles []ai.AIArticle
	err := global.GVA_DB.Model(&ai.AIArticle{}).Where("id in ?", ids).Find(&articles).Error

	for _, item := range articles {
		err = exa.Publish1Article(item)
		if err != nil {
			global.GVA_LOG.Error("发布失败!", zap.Error(err), zap.String("item", item.Title))
			item.ArticleStatus = 4
			item.ErrMessage = err.Error()
		}

		item.ArticleStatus = 2

	}
	return err
}

func (exa *AIArticleService) UpdateArticle(aiArticle ai.AIArticle) (err error) {
	// 更新
	aiArticle.PublishTime = time.Now()
	err = global.GVA_DB.Save(&aiArticle).Error
	return err
}

// @function: ApprovalArticle
// @description: 发布文章
// @param: e model.AIArticle
// @return: err error

func (exa *AIArticleService) Publish1Article(aiArticle ai.AIArticle) (err error) {
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
			articleContext := ArticlePipelineApp.Run(item.AppId, item.Topic)
			if articleContext.Content == "" {
				return
			}

			aiArticle := ai.AIArticle{
				Title:             exa.parseTitle(articleContext.Title),
				TargetAccountName: item.AccountName,
				TargetAccountId:   item.AppId,
				Topic:             articleContext.Topic,
				AuthorName:        item.DefaultAuthorName,
				Tags:              strings.Join(articleContext.Tags, ","),
				Content:           exa.parseContent(articleContext.Content),
			}
			aiArticle.BASEMODEL = BaseModel()

			err = AIArticleServiceApp.CreateAIArticle(aiArticle)

			global.GVA_LOG.Info("AI创作完成", zap.String("appId", item.AppId), zap.String("topic", item.Topic))
		}()

	}

	return nil
}

func (exa *AIArticleService) parseTitle(title string) string {
	title = strings.ReplaceAll(title, "标题：", "")
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
		if !strings.HasPrefix(line, "标题：") {
			contentLines = append(contentLines, line)
		}
	}

	// 将剩余的行重新连接成一篇文章
	markdownContent := strings.Join(contentLines, "\n")

	//```markdown
	markdownContent = strings.ReplaceAll(markdownContent, "```markdown", "")
	markdownContent = strings.ReplaceAll(markdownContent, "```", "")

	htmlContent, _ := utils.RenderMarkdownContent(markdownContent)

	return htmlContent
}

func (exa *AIArticleService) Recreation(id uint64) (err error) {
	aiArticle := ai.AIArticle{}
	err = global.GVA_DB.Where("id = ?", id).First(&aiArticle).Error
	if err != nil {
		return err
	}

	article := ai.Article{
		Title:      aiArticle.Title,
		PortalName: aiArticle.PortalName,
		Topic:      aiArticle.Topic,
		AuthorName: aiArticle.AuthorName,
		Content:    aiArticle.Content,
	}

	chatGptResp, err := QianfanServiceApp.Recreation(article)
	if err != nil {
		return err
	}

	aiArticle.Topic = chatGptResp.Topic
	aiArticle.Content = chatGptResp.Content
	aiArticle.Title = chatGptResp.Title
	aiArticle.Tags = strings.Join(chatGptResp.Tags, ",")

	global.GVA_DB.Save(&aiArticle)

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
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&articleList).Error
	}
	return articleList, total, err
}
