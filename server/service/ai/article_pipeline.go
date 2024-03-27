package ai

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
)

type CreatePipeline interface {
	Write(context *ArticleContext) error
}

var ArticlePipelineApp = new(ArticlePipeline)

type ArticlePipeline struct {
}

func (*ArticlePipeline) Run(topic string) *ArticleContext {
	context := &ArticleContext{}
	context.Topic = topic
	switch topic {
	case "":

	default:
		err := DefaultArticlePipelineApp.Write(context)
		if err != nil {
			fmt.Println(err)
		}
	}
	return context
}

var DefaultArticlePipelineApp = new(DefaultArticlePipeline)

type DefaultArticlePipeline struct {
	ArticleWriteHandleList []ArticleWriteHandle
}

func (da *DefaultArticlePipeline) init() {
	da.ArticleWriteHandleList = []ArticleWriteHandle{
		&HotSpotWriteArticle{},
		&RecreationArticle{},
		&AIWriteArticle{},
	}
}

func (da *DefaultArticlePipeline) Write(context *ArticleContext) error {
	// 初始化
	da.init()

	for _, handle := range da.ArticleWriteHandleList {
		err := handle.Handle(context)
		if err != nil {
			fmt.Println("DefaultArticlePipeline Write", err)
			continue
		}
		// 完成写作
		if context.Content != "" && len(context.Content) > 500 {
			break
		}

	}

	return nil
}

type ArticleWriteHandle interface {
	Handle(context *ArticleContext) error
}

type AddImagesHandle interface {
	AddImages(context *ArticleContext) error
}

type ArticleContext struct {
	Title   string
	Content string
	Topic   string
	Tags    []string
}

// RecreationArticle 文章改写 or  二创
type RecreationArticle struct {
}

func (r *RecreationArticle) Handle(context *ArticleContext) error {
	article := ai.Article{}
	err := global.GVA_DB.Where("topic like ?", "%"+context.Topic+"%").Where("use_times=0").Order("publish_time desc").Last(&article).Error
	if err != nil {
		return err
	}

	result, err := ChatModelServiceApp.Recreation(article, AllModel)
	if err != nil {
		return err
	}

	context.Content = result.Content
	context.Title = result.Title

	// 更新使用次数
	article.UseTimes = article.UseTimes + 1
	err = global.GVA_DB.Save(&article).Error

	return err
}

type HotSpotWriteArticle struct {
}

func (r *HotSpotWriteArticle) Handle(context *ArticleContext) error {
	hotspot := ai.Hotspot{}
	err := global.GVA_DB.Where("topic like ?", "%"+context.Topic+"%").Where("use_times=0").Order("trending desc").Last(&hotspot).Error
	if err != nil {
		return err
	}

	result, err := ChatModelServiceApp.HotSpotWrite(hotspot.Headline, AllModel)
	if err != nil {
		return err
	}

	context.Content = result.Content
	context.Title = result.Title

	// 更新使用次数
	hotspot.UseTimes = hotspot.UseTimes + 1

	err = global.GVA_DB.Save(&hotspot).Error
	return err
}

// AIWriteArticle ai 写作
type AIWriteArticle struct {
}

func (r *AIWriteArticle) Handle(context *ArticleContext) error {

	result, err := ChatModelServiceApp.TopicWrite(context.Topic, AllModel)
	if err != nil {
		return err
	}

	context.Content = result.Content
	context.Title = result.Title

	return err
}
