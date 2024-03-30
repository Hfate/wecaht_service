package ai

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

type CreatePipeline interface {
	Write(context *ArticleContext) error
}

var ArticlePipelineApp = new(ArticlePipeline)

type ArticlePipeline struct {
}

func (*ArticlePipeline) Run(appId, topic string) *ArticleContext {
	context := &ArticleContext{}
	context.Topic = topic
	context.AppId = appId

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
	AddImageHandleList     []AddImagesHandle
}

func (da *DefaultArticlePipeline) init() {
	da.ArticleWriteHandleList = []ArticleWriteHandle{
		&HotSpotWriteArticle{},
		&RecreationArticle{},
		&AIWriteArticle{},
	}
	da.AddImageHandleList = []AddImagesHandle{
		&BaiduAddImage{},
	}
}

type BaiduAddImage struct {
}

func (*BaiduAddImage) Handle(context *ArticleContext) error {

	// 正则表达式匹配Markdown中的图片占位符描述
	re := regexp.MustCompile(`\[img\](.*?)\[/img\]`)
	matches := re.FindAllStringSubmatch(context.Content, -1)

	if len(matches) == 0 {
		return nil
	}
	context.Content = strings.ReplaceAll(context.Content, "[img]", "")
	context.Content = strings.ReplaceAll(context.Content, "[/img]", "")

	for _, match := range matches {
		// 搜索图片
		filePath := utils.SearchAndSave(match[1])

		if filePath == "" {
			continue
		}

		link, err := MediaServiceApp.ImageUpload(context.AppId, filePath)

		global.GVA_LOG.Info("公众号配图", zap.String("URL", link), zap.String("文案", match[1]),
			zap.String("appId", context.AppId), zap.Error(err))

		if err != nil {
			fmt.Println(err)
			continue
		}

		wxImgFmt := "<img src=\"%s\">"
		wxImgUrl := fmt.Sprintf(wxImgFmt, link)

		context.Content = strings.ReplaceAll(context.Content, match[1], "<p>"+wxImgUrl+"</p>")
	}

	return nil
}

func (da *DefaultArticlePipeline) Write(context *ArticleContext) error {
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
			global.GVA_LOG.Info("完成AI创作", zap.String("appID", context.AppId), zap.String("title", context.Title))
			break
		}
	}

	for _, handle := range da.AddImageHandleList {
		err := handle.Handle(context)
		if err != nil {
			fmt.Println("DefaultArticlePipeline Add Image", err)
			continue
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
	AppId   string
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
	err := global.GVA_DB.Where("topic like ?", "%"+context.Topic).Where("use_times=0").Order("publish_time desc").Last(&article).Error
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
	err := global.GVA_DB.Where("topic like ?", "%"+context.Topic).Where("use_times=0").Order("trending desc").Last(&hotspot).Error
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
