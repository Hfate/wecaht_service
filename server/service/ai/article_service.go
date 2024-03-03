package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/wechat/request"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

type ArticleService struct {
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteFileChunk
//@description: 删除文章
//@param: e model.Article
//@return: err error

func (exa *ArticleService) DeleteArticle(e wechat.Article) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetArticle
//@description: 获取文章信息
//@param: id uint
//@return: customer model.Article, err error

func (exa *ArticleService) GetArticle(id uint64) (portal wechat.Article, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&portal).Error
	return
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
	var articleList []wechat.Article

	db := global.GVA_DB.Model(&wechat.Article{})

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
