package ai

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/convertor"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"strings"
)

type ArticleService struct {
}

var ArticleServiceApp = new(ArticleService)

func (exa *ArticleService) FindHotArticleByTopic(topic string, portNameList []string) (ai.Article, error) {
	var article ai.Article
	err := global.GVA_DB.Model(&ai.Article{}).Where("topic = ?", topic).
		Where("portal_name not in ?", portNameList).
		Where("use_times = 0").
		Where("content != ''").
		Order("read_num desc,publish_time desc").First(&article).Error
	return article, err
}

func (exa *ArticleService) Create(article ai.Article) error {
	return global.GVA_DB.Model(&ai.Article{}).Create(article).Error
}

func (exa *ArticleService) Update(article ai.Article) error {
	return global.GVA_DB.Save(article).Error
}

func (exa *ArticleService) CheckUrl(url string) bool {
	var count int64
	err := global.GVA_DB.Model(&ai.Article{}).Where("link = ?", url).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (exa *ArticleService) UploadArticle(file *multipart.FileHeader) error {

	fileReader, err := file.Open()
	if err != nil {
		global.GVA_LOG.Error("Error opening file header", zap.Error(err))
		return err
	}
	defer fileReader.Close()

	list := make([]ai.ArticleExclUpload, 0)

	objHeaders, headers, err := utils.ReadExcelByReader(
		fileReader, // Path of the csv file
		&list,      // A pointer to the create slice )
	)
	if err != nil {
		return err
	}

	//  Header 验证
	objHeaderSet := convertor.StringListToSet(objHeaders)
	for _, header := range headers {
		if header == "" {
			continue
		}
		if !objHeaderSet.Contains(header) {
			return errors.New("文件表头不正确")
		}
	}

	articleList := make([]*ai.Article, 0)
	for _, item := range list {

		article := &ai.Article{
			Title:       item.Title,
			Topic:       item.Topic,
			Link:        item.Link,
			PublishTime: item.PublishTime,
			Content:     item.Content,
			Comment:     item.Comment,
			LikeNum:     cast.ToInt(item.LikeNum),
			ReadNum:     cast.ToInt(item.ReadNum),
			PortalName:  item.PortalName,
		}

		article.BASEMODEL = BaseModel()
		articleList = append(articleList, article)
	}

	return global.GVA_DB.Create(&articleList).Error
}

//@function: DeleteFileChunk
//@description: 删除文章
//@param: e model.Article
//@return: err error

func (exa *ArticleService) DeleteArticle(e ai.Article) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@function: DeleteArticlesByIds
//@description: 删除选中文章
//@param: ids []wechat.Article
//@return: err error

func (exa *ArticleService) DeleteArticlesByIds(ids request.IdsReq) (err error) {
	var articles []ai.Article
	err = global.GVA_DB.Find(&articles, "id in ?", ids.Ids).Delete(&articles).Error
	return err
}

//@function: GetArticle
//@description: 获取文章信息
//@param: id uint
//@return: customer model.Article, err error

func (exa *ArticleService) GetArticle(id uint64) (article ai.Article, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&article).Error
	return
}

func (exa *ArticleService) Recreation(id uint64) error {
	article := ai.Article{}
	err := global.GVA_DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		return err
	}

	context := &ArticleContext{
		Topic: article.Topic,
		Link:  article.Link,
		Title: article.Title,
	}

	chatGptResp, err := ChatModelServiceApp.Recreation(context)
	if err != nil {
		return err
	}

	aiArticle := ai.AIArticle{
		OriginId:   cast.ToUint64(article.ID),
		Title:      chatGptResp.Title,
		PortalName: article.PortalName,
		Topic:      chatGptResp.Topic,
		AuthorName: article.AuthorName,
		Tags:       strings.Join(chatGptResp.Tags, ","),
		Content:    chatGptResp.Content,
	}
	aiArticle.BASEMODEL = BaseModel()

	err = AIArticleServiceApp.CreateAIArticle(aiArticle)

	if err != nil {
		return err
	}
	article.UseTimes++
	err = global.GVA_DB.Save(&article).Error

	return err
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetArticleList
// @description: 分页获取文章列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *ArticleService) GetArticleList(sysUserAuthorityID uint, info aiReq.ArticleSearch) (list interface{}, total int64, err error) {
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
	var articleList []ai.Article

	db := global.GVA_DB.Model(&ai.Article{})

	// 如果有条件搜索 下方会自动创建搜索语句
	if info.PortalName != "" {
		db = db.Where("portal_name LIKE ?", "%"+info.PortalName+"%")
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
		err = db.Limit(limit).Offset(offset).Order("publish_time desc,read_num desc").Find(&articleList).Error
	}
	return articleList, total, err
}

func (exa *ArticleService) Download(c *gin.Context, sysUserAuthorityID uint, info aiReq.ArticleSearch) {
	list, _, err := exa.GetArticleList(sysUserAuthorityID, info)
	if err != nil {
		fmt.Println(err)
		return
	}
	articleList := list.([]ai.Article)
	result := make([]*ai.ArticleExcl, 0)
	for _, item := range articleList {
		result = append(result, &ai.ArticleExcl{
			AuthorName: item.AuthorName,
			PortalName: item.PortalName,
			Topic:      item.Topic,
			Title:      item.Title,
			Link:       item.Link,
			ReadNum:    item.ReadNum,
			CommentNum: item.CommentNum,
			LikeNum:    item.LikeNum,
			Content:    item.Content,
		})
	}

	excelFile := excelize.NewFile()

	utils.WriteDefaultExcelSheet(excelFile, result)

	fileName := "article.xlsx"
	filePath := "./tmp/" + fileName

	err = excelFile.SaveAs(filePath)

	if err != nil {
		fmt.Println(err)
		return
	}
	// 中文编码
	fileName = utils.EncodeFilename(fileName)

	//返回文件流
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.File(filePath)

	return

}

func (exa *ArticleService) Spider() {
	// 处理记录，批处理大小为100
	results := make([]*ai.Article, 0)

	index := 0

	result := global.GVA_DB.Where("content = ''").Where("tags=''").Order("publish_time desc").FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			// 对批中的每条记录进行操作

			content := utils.CollectWechatArticle(result.Link)

			result.Content = content

			if len(content) == 0 {
				result.Tags = "delete"
			}

			index++

			fmt.Println(index)
		}

		// 保存对当前批记录的修改
		tx.Save(&results)

		fmt.Println("当前处理批次" + cast.ToString(batch))
		return nil
	})

	fmt.Println(result.Error.Error())
}
