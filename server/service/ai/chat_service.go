package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timeutil"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type ChatService struct {
}

var ChatServiceApp = new(ChatService)

func (cs *ChatService) GetKeyWord(title string, chatModel config.ChatModel) string {
	chatGptPrompt := "你现在是一名爆文写手，特别擅长从文章标题中找到关键词。我将给一个文章标题，需要你帮忙提取标题中的一个关键词用以做图片搜索。如果找不到关键词，可以返回该标题的主题，例如：历史，职场，明星等等" +
		"\n举例   " +
		"\n文章标题：中瑙友谊再升华，开启双边合作新篇章。  关键词：友谊再升华" +
		"\n文章标题：周处传奇：除三害、转人生，英雄之路的跌宕起伏  关键词：周处除三害" +
		"\n文章标题：" + title

	chatMessageHistory := []*ChatMessage{ChatSystemMessage}

	resp, chatMessageHistory, err := ChatServiceApp.ChatWithModel(chatGptPrompt, chatMessageHistory, chatModel)
	if err != nil || len(resp) > 10 {
		resp = "夜晚的城市"
	}

	return resp
}

func (cs *ChatService) HotSpotWrite(context *ArticleContext, chatModel config.ChatModel) (*ArticleContext, error) {

	chatGptPromptList, err := parsePrompt(context, ai.ContentRecreation)
	if err != nil {
		return &ArticleContext{}, err
	}

	chatMessageHistory := []*ChatMessage{ChatSystemMessage}
	resp := ""

	for _, chatGptPrompt := range chatGptPromptList {
		resp, chatMessageHistory, err = ChatServiceApp.ChatWithModel(chatGptPrompt, chatMessageHistory, chatModel)
		if err != nil {
			return nil, err
		}
		context.Content = resp
	}

	chatGptPromptList, err = parsePrompt(context, ai.TitleCreate)
	if err != nil {
		return nil, err
	}

	for _, chatGptPrompt := range chatGptPromptList {
		resp, chatMessageHistory, err = ChatServiceApp.ChatWithModel(chatGptPrompt, chatMessageHistory, chatModel)
		if err != nil {
			return nil, err
		}
		context.Title = resp
	}

	chatGptPromptList, err = parsePrompt(context, ai.AddImage)
	if err != nil {
		return nil, err
	}

	for _, chatGptPrompt := range chatGptPromptList {
		resp, chatMessageHistory, err = ChatServiceApp.ChatWithModel(chatGptPrompt, chatMessageHistory, chatModel)
		if err != nil {
			return nil, err
		}

		context.Content = resp
	}

	context.Params = []string{chatModel.ModelType, "HotSpotWrite"}
	return context, nil
}

func (cs *ChatService) Recreation(articleContext *ArticleContext, chatModel config.ChatModel) (*ArticleContext, error) {
	starTime := timeutil.GetCurTime()

	// 重置进度
	aiArticle := articleContext.Article
	aiArticle.CreatedAt = time.Now()
	aiArticle.Percent = 0
	aiArticle.ProcessStatus = ai.ProcessCreateIng
	aiArticle.ProcessParams = "开始执行文章改写"
	global.GVA_DB.Save(&aiArticle)

	// 文章改写
	chatGptPromptList, err := parsePrompt(articleContext, ai.ContentRecreation)
	if err != nil {
		return &ArticleContext{}, err
	}

	chatMessageHistory := []*ChatMessage{ChatSystemMessage}
	resp := ""

	size := cast.ToString(len(chatGptPromptList))
	for index, chatGptPrompt := range chatGptPromptList {
		aiArticle.ProcessParams = "【" + chatModel.ModelType + "】文章改写:prompt执行进度[" + cast.ToString(index+1) + "/" + size + "]"
		global.GVA_DB.Save(&aiArticle)
		resp, chatMessageHistory, err = ChatServiceApp.ChatWithModel(chatGptPrompt, chatMessageHistory, chatModel)
		if err != nil {
			return nil, err
		}
		articleContext.Content = resp
	}

	//if len(articleContext.Content) < 1200 {
	//	// 文章扩写
	//	chatGptPromptList, err = parsePrompt(articleContext, ai.ArticleExpansion)
	//	if err != nil {
	//		return nil, err
	//	}
	//	for index, chatGptPrompt := range chatGptPromptList {
	//		aiArticle.ProcessParams = "【" + chatModel.ModelType + "】文章扩写:正在执行第" + cast.ToString(index+1) + "条prompt"
	//		global.GVA_DB.Save(&aiArticle)
	//		resp, chatMessageHistory, err = ChatServiceApp.ChatWithModel(chatGptPrompt, chatMessageHistory, chatModel)
	//		if err != nil {
	//			return nil, err
	//		}
	//		articleContext.Content = resp
	//	}
	//}

	//// 标题创建
	//chatGptPromptList, err = parsePrompt(articleContext, ai.TitleCreate)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, chatGptPrompt := range chatGptPromptList {
	//	aiArticle.ProcessParams = "【" + chatModel.ModelType + "】标题创建prompt执行ing"
	//	global.GVA_DB.Save(&aiArticle)
	//	resp, chatMessageHistory, err = ChatServiceApp.ChatWithModel(chatGptPrompt, chatMessageHistory, chatModel)
	//	if err != nil {
	//		return nil, err
	//	}
	//	articleContext.Title = resp
	//}

	// 文章配图
	aiArticle.ProcessStatus = ai.ProcessAddImgIng
	// 更新进度
	global.GVA_DB.Save(&aiArticle)

	// 文章配图
	chatGptPromptList, err = parsePrompt(articleContext, ai.AddImage)
	if err != nil {
		return nil, err
	}

	size = cast.ToString(len(chatGptPromptList))
	for index, chatGptPrompt := range chatGptPromptList {
		aiArticle.ProcessParams = "【" + chatModel.ModelType + "】文章配图:prompt执行进度[" + cast.ToString(index+1) + "/" + size + "]"
		global.GVA_DB.Save(&aiArticle)
		resp, chatMessageHistory, err = ChatServiceApp.ChatWithModel(chatGptPrompt, chatMessageHistory, chatModel)
		if err != nil {
			return nil, err
		}

		resp = strings.ReplaceAll(resp, "```json", "")
		resp = strings.ReplaceAll(resp, "```", "")

		addImgResp := &AddImgResp{}

		err = utils.JsonStrToStruct(resp, addImgResp)
		if err == nil && addImgResp.Image1Description != "" && addImgResp.Image2Description != "" {
			img1 := cs.SearchAndSave(addImgResp.Image1Description)
			img2 := cs.SearchAndSave(addImgResp.Image2Description)

			if strings.Contains(img1, "http") {
				articleContext.Content = "![" + addImgResp.Image1Description + "](" + img1 + ")" + "\n" + articleContext.Content
			}

			if strings.Contains(img2, "http") {
				articleContext.Content = articleContext.Content + "\n" + "![" + addImgResp.Image2Description + "](" + img2 + ")"
			}
		}

	}

	aiArticle.ProcessStatus = ai.ProcessCreated
	aiArticle.ProcessParams = "创作完成"
	aiArticle.Percent = 100
	aiArticle.Params = chatModel.ModelType + "," + "Recreation"
	aiArticle.Content = articleContext.Content
	aiArticle.Context = utils.Parse2Json(chatMessageHistory)

	// 更新进度
	global.GVA_DB.Save(&aiArticle)

	go cs.GetSimilarity(aiArticle)

	articleContext.Params = []string{chatModel.ModelType, "Recreation"}

	endTime := timeutil.GetCurTime()

	// 更新时长
	AvgTimeServiceApp.UpdateAvgTime(endTime - starTime)

	return articleContext, nil
}

func (cs *ChatService) GetSimilarity(aiArticle *ai.AIArticle) {
	req := &SimilarityReq{
		Text:    aiArticle.OriginContent,
		Compare: []string{aiArticle.Content},
	}

	resp, err := SiTongServiceApp.Similarity(req)
	if err != nil {
		global.GVA_LOG.Error("GetSimilarity", zap.Error(err))
		return
	}

	if resp != nil && len(resp.Results) > 0 {
		// 更新相似度
		similarity := resp.Results[0].Similarity
		aiArticle.Similarity = similarity
		global.GVA_DB.Save(&aiArticle)
	}

}

func (cs *ChatService) SearchAndSave(keyword string) string {
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
	return cs.saveImage(imgUrlList)
}

func (cs *ChatService) saveImage(imgUrlList []string) string {
	// 通过第一张图片链接下载图片
	filePath := ""

	for _, imgUrl := range imgUrlList {
		ossFilePath, err := cs.downloadImage(imgUrl)
		if err != nil {
			global.GVA_LOG.Info("downloadImage failed", zap.Any("err", err.Error()))
		} else {
			filePath = ossFilePath
			break
		}
	}

	if !strings.Contains(filePath, "http") {
		filePath = "https://" + filePath
	}

	return filePath
}

// DownloadImage 从 URL 下载图片并上传到 OSS
func (cs *ChatService) downloadImage(imageUrl string) (string, error) {
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

//func (ba *BaiduAddImage) SearchAndSave(keyword string) string {
//	imgUrlList := make([]string, 0)
//
//	baiduImgUrlList := utils.CollectBaiduImgUrl(keyword)
//	if len(baiduImgUrlList) > 0 {
//		imgUrlList = append(imgUrlList, baiduImgUrlList...)
//	}
//
//	unsplashImgUrlList := utils.CollectUnsplashImgUrl(keyword)
//	if len(unsplashImgUrlList) > 0 {
//		imgUrlList = append(imgUrlList, unsplashImgUrlList...)
//	}
//
//	// 通过第一张图片链接下载图片
//	return ba.saveImage(imgUrlList)
//}
//
//var BaiduAddImageApp = new(BaiduAddImage)
//
//type BaiduAddImage struct {
//}
//
//func (ba *BaiduAddImage) Handle(context *ArticleContext) error {
//
//	// 正则表达式匹配Markdown中的图片占位符描述
//	re := regexp.MustCompile(`\[img\](.*?)\[/img\]`)
//	matches := re.FindAllStringSubmatch(context.Content, -1)
//
//	if len(matches) == 0 {
//		return nil
//	}
//	context.Content = strings.ReplaceAll(context.Content, "[img]", "")
//	context.Content = strings.ReplaceAll(context.Content, "[/img]", "")
//
//	aiArticle := context.Article
//	size := cast.ToString(len(matches))
//	for index, match := range matches {
//		aiArticle.ProcessParams = "文章配图:正在下载.[" + cast.ToString(index+1) + "/" + size + "]"
//		global.GVA_DB.Save(&aiArticle)
//
//		// 搜索图片
//		filePath := ba.SearchAndSave(match[1])
//
//		if filePath == "" {
//			continue
//		}
//
//		if !strings.Contains(filePath, "http") {
//			filePath = "https://" + filePath
//		}
//
//		wxImgFmt := "<img src=\"%s\">"
//		wxImgUrl := fmt.Sprintf(wxImgFmt, filePath)
//
//		context.Content = strings.ReplaceAll(context.Content, match[1], "\n"+wxImgUrl+"\n")
//	}
//
//	aiArticle.Content = context.Content
//	global.GVA_DB.Save(&aiArticle)
//
//	return nil
//}
//
//func (ba *BaiduAddImage) saveImage(imgUrlList []string) string {
//	// 通过第一张图片链接下载图片
//	filePath := ""
//
//	for _, imgUrl := range imgUrlList {
//		ossFilePath, err := ba.downloadImage(imgUrl)
//		if err != nil {
//			global.GVA_LOG.Info("downloadImage failed", zap.Any("err", err.Error()))
//		} else {
//			filePath = ossFilePath
//			break
//		}
//	}
//	return filePath
//}

var ChatSystemMessage = &ChatMessage{
	Role:    "system",
	Content: "你是一个人工智能助手，你更擅长中文的对话。你会为用户提供安全，有帮助，准确的回答.",
}

func (cs *ChatService) ChatWithModel(message string, history []*ChatMessage, chatModel config.ChatModel) (string, []*ChatMessage, error) {
	apiUrl := chatModel.ApiUrl
	model := chatModel.Model
	refreshToken := chatModel.RefreshToken

	history = append(history, &ChatMessage{
		Role:    "user",
		Content: message,
	})

	chatReq := &ChatReq{
		Model:    model,
		Messages: history,
	}

	statusCode, respBody, err := utils.PostWithHeaders(apiUrl, utils.Parse2Json(chatReq), map[string]string{
		"Authorization": "Bearer " + refreshToken,
	})

	if err != nil {
		return "", history, err
	}

	if statusCode != 200 {
		return "", history, errors.New(string(respBody))
	}

	chatResp := &ChatResp{}
	err = utils.JsonStrToStruct(string(respBody), chatResp)
	if err != nil {
		return "", history, err
	}

	if len(chatResp.Choices) > 0 {
		history = append(history, &ChatMessage{
			Role:    "assistant",
			Content: chatResp.Choices[0].Message.Content,
		})
		return chatResp.Choices[0].Message.Content, history, err
	}

	return "", history, errors.New("chat回复为空")
}

type ChatReq struct {
	Model       string         `json:"model"`
	Messages    []*ChatMessage `json:"messages"`
	Temperature float64        `json:"temperature"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResp struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func parsePrompt(context *ArticleContext, promptType int) ([]string, error) {
	topic := context.Topic
	prompt, err := PromptServiceApp.FindPromptByTopicAndType(topic, promptType)
	if err != nil {
		return []string{}, err
	}

	promptList := make([]string, 0)
	// 使用json.Unmarshal将JSON字符串解析到字符串切片
	err = json.Unmarshal([]byte(utils.EscapeSpecialCharacters(prompt)), &promptList)

	if len(promptList) == 0 {
		promptList = []string{prompt}
	}

	result := make([]string, 0)

	for _, item := range promptList {
		temp := template.New("ChatGptPrompt")
		tmpl, err2 := temp.Parse(item)
		if err2 != nil {
			global.GVA_LOG.Info("模板解析失败")
			return []string{}, err2
		}

		// 创建一个缓冲区来保存模板生成的结果
		var buf bytes.Buffer
		// 使用模板和数据生成输出
		err = tmpl.Execute(&buf, context)

		chatGptPrompt := buf.String()

		result = append(result, chatGptPrompt)
	}

	return result, err
}

type AddImgResp struct {
	Image1Description string `json:"Image1Description"`
	Image2Description string `json:"Image2Description"`
}
