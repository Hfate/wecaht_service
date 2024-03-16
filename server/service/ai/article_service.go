package ai

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ArticleService struct {
}

var ArticleServiceApp = new(ArticleService)

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

func (exa *ArticleService) GetArticle(id uint64) (portal ai.Article, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&portal).Error
	return
}

func (exa *ArticleService) Recreation(id uint64) error {
	article := ai.Article{}
	err := global.GVA_DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		return err
	}

	recreationTitle, recreationContent := QianfanServiceApp.Recreation(article)

	aiArticle := ai.AIArticle{
		OriginId:   article.ID,
		Title:      recreationTitle,
		PortalName: article.PortalName,
		Topic:      article.Topic,
		AuthorName: article.AuthorName,
		Tags:       article.Tags,
		Content:    recreationContent,
	}
	aiArticle.BASEMODEL = BaseModel()

	err = AIArticleServiceApp.CreateAIArticle(aiArticle)

	if err != nil {
		return err
	}
	article.RecreationNum++
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
		err = db.Limit(limit).Offset(offset).Order("publish_time desc").Find(&articleList).Error
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
