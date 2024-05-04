package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timeutil"
	"go.uber.org/zap"
	"strings"
	"time"
)

type CreatePipeline interface {
	Write(context *ArticleContext) error
}

var ArticlePipelineApp = new(ArticlePipeline)

type ArticlePipeline struct {
}

func (*ArticlePipeline) Run(model string, context *ArticleContext) error {
	switch model {
	case "hotspot":
		return HotSpotArticlePipelineApp.Write(context)
	default:
		return DefaultArticlePipelineApp.Write(context)
	}
}

var HotSpotArticlePipelineApp = new(HotSpotArticlePipeline)

type HotSpotArticlePipeline struct {
	ArticleWriteHandleList []ArticleWriteHandle
}

func (da *HotSpotArticlePipeline) init() {
	da.ArticleWriteHandleList = []ArticleWriteHandle{
		&HotSpotWriteArticle{},
	}
}

func (da *HotSpotArticlePipeline) Write(context *ArticleContext) error {
	// 初始化
	da.init()

	size := len(da.ArticleWriteHandleList)

	for i := 0; i < size; i++ {
		handle := da.ArticleWriteHandleList[i]
		err := handle.Handle(context)
		if err != nil {
			continue
		}
		// 完成写作
		if context.Content != "" && len(context.Content) > 500 {
			global.GVA_LOG.Info("完成AI创作", zap.String("AccountName", context.Account.AccountName), zap.String("title", context.Title))
			break
		}
	}

	return nil
}

var DefaultArticlePipelineApp = new(DefaultArticlePipeline)

type DefaultArticlePipeline struct {
	ArticleWriteHandleList []ArticleWriteHandle
}

func (da *DefaultArticlePipeline) Write(context *ArticleContext) error {

	articleWriteHandleList := make([]ArticleWriteHandle, 0)

	if strings.Contains(context.CreateTypes, "2") {
		articleWriteHandleList = append(articleWriteHandleList, &RecreationArticle{})
	}

	if strings.Contains(context.CreateTypes, "1") {
		articleWriteHandleList = append(articleWriteHandleList, &HotSpotWriteArticle{})
	}

	for _, handle := range articleWriteHandleList {
		err := handle.Handle(context)
		if err != nil {
			global.GVA_LOG.Error("文章写作失败", zap.Error(err))
			continue
		}
		// 完成写作
		if context.Content != "" && len(context.Content) > 1000 && len(context.Params) > 0 {
			global.GVA_LOG.Info("完成AI创作", zap.String("accountName", context.Account.AccountName), zap.String("title", context.Title))
			break
		}
	}

	return nil
}

type ArticleWriteHandle interface {
	Handle(context *ArticleContext) error
}

type AddImagesHandle interface {
	Handle(context *ArticleContext) error
}

type ArticleContext struct {
	Account     *ai.OfficialAccount
	Article     *ai.AIArticle
	Title       string
	Content     string
	Topic       string
	HotspotId   uint64
	Link        string
	Tags        []string
	Params      []string
	CreateTypes string
}

// RecreationArticle 文章改写 or  二创
type RecreationArticle struct {
}

func (r *RecreationArticle) Handle(context *ArticleContext) error {
	batchId := timeutil.GetCurDate() + context.Account.AppId

	article := ai.DailyArticle{}
	err := global.GVA_DB.Where("batch_id = ?", batchId).Where("use_times=0").Order("publish_time desc").Last(&article).Error
	if err != nil {
		return err
	}

	// 更新使用次数
	article.UseTimes = article.UseTimes + 1
	err = global.GVA_DB.Save(&article).Error
	if err != nil {
		return err
	}

	// 生成文章初稿
	aiArticle := &ai.AIArticle{
		BatchId:           batchId,
		Title:             article.Title,
		TargetAccountName: context.Account.AccountName,
		TargetAccountId:   context.Account.AppId,
		Topic:             context.Account.Topic,
		AuthorName:        context.Account.DefaultAuthorName,
		Link:              article.Link,
		Content:           article.Content,
		OriginContent:     article.Content,
		PublishTime:       time.Now(),
		ProcessParams:     "任务新创建，正在等待执行..",
	}
	aiArticle.BASEMODEL = BaseModel()
	err = global.GVA_DB.Model(&ai.AIArticle{}).Create(aiArticle).Error
	if err != nil {
		return err
	}

	context.Title = article.Title
	context.Content = article.Content
	context.Link = article.Link
	context.Article = aiArticle

	PoolServiceApp.SubmitBizTask(func() {

		aiArticle.ProcessStatus = ai.ProcessCreateIng
		// 更新进度
		global.GVA_DB.Save(&aiArticle)

		result, err2 := ChatModelServiceApp.Recreation(context)
		if err2 != nil {
			global.GVA_LOG.Error("Recreation", zap.Error(err2))
			return
		}

		if result.Title == article.Title {
			aiArticle.ProcessParams = "创作失败"
			aiArticle.ProcessStatus = ai.ProcessFail
			aiArticle.Percent = 100
		}

		context.Content = result.Content
		context.Title = result.Title

		//
		aiArticle.Title = ArticleContentHandlerApp.HandleTitle(result.Title)
		//  处理排版
		aiArticle.Content = ArticleContentHandlerApp.Handle(context.Account, result.Content)

		// 更新文章内容
		err2 = global.GVA_DB.Save(&aiArticle).Error
		if err2 != nil {
			global.GVA_LOG.Error("Recreation", zap.Error(err2))
		}
	})

	return err
}

type HotSpotWriteArticle struct {
}

func (r *HotSpotWriteArticle) Handle(context *ArticleContext) error {

	hotspotList := make([]ai.Hotspot, 0)

	if context.HotspotId != 0 {
		hotspot := ai.Hotspot{}
		err := global.GVA_DB.Where("id = ?", context.HotspotId).Order("created_at desc ,trending desc").Last(&hotspot).Error
		if err != nil {
			return err
		}
		hotspotList = append(hotspotList, hotspot)
	} else {
		err := global.GVA_DB.Where("topic like ?", "%"+context.Topic).Where("use_times=0").Order("created_at desc ,trending desc").Limit(10).Find(&hotspotList).Error
		if err != nil {
			return err
		}
	}

	for _, hotspot := range hotspotList {

		// 更新使用次数
		hotspot.UseTimes = hotspot.UseTimes + 1
		global.GVA_DB.Save(&hotspot)

		oldTopic := context.Topic

		// 找有没有相关的热点文章素材
		article := &ai.Article{}

		err := global.GVA_DB.Model(&ai.Article{}).Where("hotspot_id = ?", hotspot.ID).Last(&article).Error
		if err != nil {
			global.GVA_LOG.Info("没有找到热点词条相关的文章[" + hotspot.Headline + "]")
			continue
		}

		context.Content = article.Content

		result, err := ChatModelServiceApp.HotSpotWrite(context)
		if err != nil {
			global.GVA_LOG.Info("热点词条创作失败" + hotspot.Headline)
			context.Topic = oldTopic
			continue
		}

		context.Content = result.Content
		context.Title = result.Title
		context.Topic = oldTopic

		return nil
	}

	return nil
}
