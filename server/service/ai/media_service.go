package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"mime/multipart"
	"strings"
)

type MediaService struct {
}

var MediaServiceApp = new(MediaService)

//@function: CreateMedia
//@description: 创建素材
//@param: e model.Media
//@return: err error

func (exa *MediaService) CreateMedia(targetAccountId string, header *multipart.FileHeader) (err error) {
	officialAccount, err := OfficialAccountServiceApp.GetOfficialAccountByAppId(targetAccountId)
	if err != nil {
		return err
	}

	tagArr := strings.Split(header.Filename, ".")

	oss := upload.NewOss()
	filePath, _, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		return uploadErr
	}

	mediaID, url, err := WechatServiceApp.AddMaterial(officialAccount, filePath)
	if err != nil {
		return err
	}

	count := exa.CountByAccountId(targetAccountId)

	media := &ai.Media{
		Topic:             officialAccount.Topic,
		MediaID:           mediaID,
		Link:              url,
		FileName:          header.Filename,
		Tag:               tagArr[len(tagArr)-1],
		TargetAccountId:   targetAccountId,
		TargetAccountName: officialAccount.AccountName,
		SeqNum:            int(count) + 1,
	}

	media.BASEMODEL = BaseModel()
	err = global.GVA_DB.Create(&media).Error
	return err
}

func (exa *MediaService) ImageUpload(targetAccountId string, filePath string) (url string, err error) {
	officialAccount, err := OfficialAccountServiceApp.GetOfficialAccountByAppId(targetAccountId)
	if err != nil {
		return "", err
	}
	url, err = WechatServiceApp.ImageUpload(officialAccount, filePath)

	return

}

func (exa *MediaService) CreateMediaByPath(targetAccountId string, filePath string) (mediaID, url string, err error) {
	officialAccount, err := OfficialAccountServiceApp.GetOfficialAccountByAppId(targetAccountId)
	if err != nil {
		return "", "", err
	}

	mediaID, url, err = WechatServiceApp.AddMaterial(officialAccount, filePath)
	if err != nil {
		return "", "", err
	}

	tagArr := strings.Split(filePath, ".")
	pathArr := strings.Split(filePath, "/")

	count := exa.CountByAccountId(targetAccountId)

	media := &ai.Media{
		Topic:             officialAccount.Topic,
		MediaID:           mediaID,
		Link:              url,
		FileName:          pathArr[len(pathArr)-1],
		Tag:               tagArr[len(tagArr)-1],
		TargetAccountId:   targetAccountId,
		TargetAccountName: officialAccount.AccountName,
		SeqNum:            int(count) + 1,
	}

	media.BASEMODEL = BaseModel()
	err = global.GVA_DB.Create(&media).Error
	return mediaID, url, err

}

func (exa *MediaService) CountByAccountId(targetAccountId string) int {
	var count int64
	global.GVA_DB.Model(&ai.Media{}).Where("target_account_id=?", targetAccountId).Count(&count)
	return int(count)
}

func (exa *MediaService) FindByAccountId(targetAccountId string, seqNum int) *ai.Media {
	var media *ai.Media

	err := global.GVA_DB.Model(&ai.Media{}).Where("target_account_id=?", targetAccountId).Where("seq_num=?", seqNum).Last(&media).Error
	if err != nil {
		global.GVA_DB.Model(&ai.Media{}).Where("target_account_id=?", targetAccountId).Last(&media)
	}

	return media
}

func (exa *MediaService) FindLast() *ai.Media {
	var media *ai.Media

	err := global.GVA_DB.Model(&ai.Media{}).Last(&media).Error
	if err != nil {

	}

	return media
}

func (exa *MediaService) RandomByAccountId(targetAccountId string) *ai.Media {
	count := exa.CountByAccountId(targetAccountId)
	if count > 1 {
		randomInt := utils.RandomInt(0, count)
		return exa.FindByAccountId(targetAccountId, randomInt)
	}

	return exa.FindLast()

}

//@function: DeleteFileChunk
//@description: 删除素材
//@param: e model.Media
//@return: err error

func (exa *MediaService) DeleteMedia(e ai.Media) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@function: UpdateMedia
//@description: 更新素材
//@param: e *model.Media
//@return: err error

func (exa *MediaService) UpdateMedia(e *ai.Media) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

//@function: GetMedia
//@description: 获取素材信息
//@param: id uint
//@return: customer model.Media, err error

func (exa *MediaService) GetMedia(id uint64) (media ai.Media, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&media).Error
	return
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetMediaList
// @description: 分页获取素材列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *MediaService) GetMediaPage(sysUserAuthorityID uint, info request.PageInfo) (list interface{}, total int64, err error) {
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
	var mediaList []ai.Media

	db := global.GVA_DB.Model(&ai.Media{})
	err = db.Count(&total).Error
	if err != nil {
		return mediaList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("created_at").Find(&mediaList).Error
	}
	return mediaList, total, err
}
