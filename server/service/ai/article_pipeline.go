package ai

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timeutil"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"go.uber.org/zap"
	"log"
	"os"
	"regexp"
	"strings"
)

type CreatePipeline interface {
	Write(context *ArticleContext) error
}

var ArticlePipelineApp = new(ArticlePipeline)

type ArticlePipeline struct {
}

func (*ArticlePipeline) Run(model string, context *ArticleContext) *ArticleContext {
	switch model {
	case "hotspot":
		err := HotSpotArticlePipelineApp.Write(context)
		if err != nil {
			global.GVA_LOG.Error("文章写作失败", zap.Error(err))
		}
	default:
		err := DefaultArticlePipelineApp.Write(context)
		if err != nil {
			global.GVA_LOG.Error("文章写作失败", zap.Error(err))
		}
	}
	return context
}

var HotSpotArticlePipelineApp = new(HotSpotArticlePipeline)

type HotSpotArticlePipeline struct {
	ArticleWriteHandleList []ArticleWriteHandle
	AddImageHandleList     []AddImagesHandle
}

func (da *HotSpotArticlePipeline) init() {
	da.ArticleWriteHandleList = []ArticleWriteHandle{
		&HotSpotWriteArticle{},
	}
	da.AddImageHandleList = []AddImagesHandle{
		&BaiduAddImage{},
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

var DefaultArticlePipelineApp = new(DefaultArticlePipeline)

type DefaultArticlePipeline struct {
	ArticleWriteHandleList []ArticleWriteHandle
	AddImageHandleList     []AddImagesHandle
}

func (da *DefaultArticlePipeline) init() {
	da.ArticleWriteHandleList = []ArticleWriteHandle{
		&HotSpotWriteArticle{},
		&RecreationArticle{},
	}
	da.AddImageHandleList = []AddImagesHandle{
		&BaiduAddImage{},
	}
}

type BaiduAddImage struct {
}

func (ba *BaiduAddImage) Handle(context *ArticleContext) error {

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
		filePath := ba.SearchAndSave(match[1])

		if filePath == "" {
			continue
		}

		if !strings.Contains(filePath, "http") {
			filePath = "https://" + filePath
		}

		wxImgFmt := "<img src=\"%s\">"
		wxImgUrl := fmt.Sprintf(wxImgFmt, filePath)

		context.Content = strings.ReplaceAll(context.Content, match[1], "\n"+wxImgUrl+"\n")
	}

	return nil
}

func (ba *BaiduAddImage) SearchAndSave(keyword string) string {
	imgUrlList := make([]string, 0)

	baiduImgUrlList := utils.CollectBaiduImgUrl(keyword)
	if len(baiduImgUrlList) > 0 {
		imgUrlList = append(imgUrlList, baiduImgUrlList...)
	}

	unsplashImgUrlList := utils.CollectUnsplashImgUrl(keyword)
	if len(unsplashImgUrlList) > 0 {
		imgUrlList = append(imgUrlList, unsplashImgUrlList...)
	}

	// 通过第一张图片链接下载图片
	return ba.saveImage(imgUrlList)
}

func (ba *BaiduAddImage) saveImage(imgUrlList []string) string {
	// 通过第一张图片链接下载图片
	filePath := ""

	for _, imgUrl := range imgUrlList {
		ossFilePath, err := ba.downloadImage(imgUrl)
		if err != nil {
			global.GVA_LOG.Info("downloadImage failed", zap.Any("err", err.Error()))
		} else {
			filePath = ossFilePath
			break
		}
	}
	return filePath
}

// DownloadImage 从 URL 下载图片并上传到 OSS
func (ba *BaiduAddImage) downloadImage(imageUrl string) (string, error) {
	// 发起 HTTP GET 请求
	tempFilePath, err := utils.CreateTempImgFile(imageUrl)
	if err != nil {
		return "", err
	}

	defer os.Remove(tempFilePath)

	// 使用 multipart.FileHeader 封装文件信息
	fileHeader, err := os.Open(tempFilePath)
	if err != nil {
		log.Println("Error opening file header:", err)
		return "", err
	}

	defer fileHeader.Close() // 创建文件 defer 关闭

	// 调用 OSS 上传方法
	oss := upload.NewOss()
	uploadUrl, _, err := oss.UploadFile(fileHeader)
	if err != nil {
		log.Println("Error uploading to OSS:", err)
		return "", err
	}

	// 返回上传的 URL 和 OSS 路径
	return uploadUrl, nil
}

func (da *DefaultArticlePipeline) Write(context *ArticleContext) error {
	// 初始化
	da.init()

	articleWriteHandleList := make([]ArticleWriteHandle, 0)

	if strings.Contains(context.CreateTypes, "2") {
		articleWriteHandleList = append(articleWriteHandleList, &RecreationArticle{})
	}

	if strings.Contains(context.CreateTypes, "1") {
		articleWriteHandleList = append(articleWriteHandleList, &HotSpotWriteArticle{})
	}

	for _, handle := range da.ArticleWriteHandleList {
		err := handle.Handle(context)
		if err != nil {
			global.GVA_LOG.Error("文章写作失败", zap.Error(err))
			continue
		}
		// 完成写作
		if context.Content != "" && len(context.Content) > 1000 && len(context.Params) > 0 {
			global.GVA_LOG.Info("完成AI创作", zap.String("appID", context.AppId), zap.String("title", context.Title))
			break
		}
	}

	for _, handle := range da.AddImageHandleList {
		err := handle.Handle(context)
		if err != nil {
			global.GVA_LOG.Error("文章配图失败", zap.Error(err))
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
	AppId       string
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
	batchId := timeutil.GetCurDate() + context.AppId

	article := ai.DailyArticle{}
	err := global.GVA_DB.Where("batch_id = ?", batchId).Where("use_times=0").Order("publish_time desc").Last(&article).Error
	if err != nil {
		return err
	}

	// 更新使用次数
	article.UseTimes = article.UseTimes + 1
	err = global.GVA_DB.Save(&article).Error

	context.Title = article.Title
	context.Content = article.Content
	context.Link = article.Link

	result, err := ChatModelServiceApp.Recreation(context)
	if err != nil {
		return err
	}

	context.Content = result.Content
	context.Title = result.Title

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
