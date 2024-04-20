package ai

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiRes "github.com/flipped-aurora/gin-vue-admin/server/model/ai/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type OfficialAccountService struct {
}

var OfficialAccountServiceApp = new(OfficialAccountService)

func (exa *OfficialAccountService) FindTopicList() []string {
	result := make([]string, 0)

	err := global.GVA_DB.Model(&ai.OfficialAccount{}).Distinct("topic").Find(&result).Error
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		return result
	}

	return result
}

func (exa *OfficialAccountService) UpdateCreateTypes(id uint64, createTypeList []int) (err error) {
	createTypes := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(createTypeList)), ","), "[]")
	err = global.GVA_DB.Model(&ai.OfficialAccount{}).Where("id = ?", id).Update("create_types", createTypes).Error
	return err
}

//@function: CreateOfficialAccount
//@description: 创建公众号
//@param: e model.OfficialAccount
//@return: err error

func (exa *OfficialAccountService) CreateOfficialAccount(e ai.OfficialAccount) (err error) {
	e.BASEMODEL = BaseModel()
	err = global.GVA_DB.Create(&e).Error
	return err
}

//@function: DeleteFileChunk
//@description: 删除公众号
//@param: e model.OfficialAccount
//@return: err error

func (exa *OfficialAccountService) DeleteOfficialAccount(e ai.OfficialAccount) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@function: UpdateOfficialAccount
//@description: 更新公众号
//@param: e *model.OfficialAccount
//@return: err error

func (exa *OfficialAccountService) UpdateOfficialAccount(e *ai.OfficialAccount) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *OfficialAccountService) CreateArticle(id uint64) error {
	officialAccount, err := exa.GetOfficialAccount(id)
	if err != nil {
		return err
	}

	go func() {
		err = AIArticleServiceApp.GenerateArticle(officialAccount)
		if err != nil {
			global.GVA_LOG.Error("GenerateArticle", zap.Error(err))
		}
	}()
	return nil
}

//@function: GetOfficialAccount
//@description: 获取公众号信息
//@param: id uint
//@return: customer model.OfficialAccount, err error

func (exa *OfficialAccountService) GetOfficialAccount(id uint64) (officialAccount *ai.OfficialAccount, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&officialAccount).Error
	return
}

//@function: GetOfficialAccount
//@description: 获取公众号信息
//@param: id uint
//@return: customer model.OfficialAccount, err error

func (exa *OfficialAccountService) FindByTopic(topic string) (officialAccount ai.OfficialAccount, err error) {
	err = global.GVA_DB.Where("topic = ?", topic).Last(&officialAccount).Error

	if err != nil {
		//找不到默认用时事主题
		err = global.GVA_DB.Where("topic = ?", "时事").Last(&officialAccount).Error
	}
	return
}

//@function: GetOfficialAccountByAppId
//@description: 获取公众号信息
//@param: id uint
//@return: customer model.OfficialAccount, err error

func (exa *OfficialAccountService) GetOfficialAccountByAppId(appId string) (officialAccount ai.OfficialAccount, err error) {
	err = global.GVA_DB.Where("app_id = ?", appId).Last(&officialAccount).Error
	return
}

//@function: GetOfficialAccount
//@description: 获取公众号信息
//@param: id uint
//@return: customer model.OfficialAccount, err error

func (exa *OfficialAccountService) List() (list []*ai.OfficialAccount, err error) {
	err = global.GVA_DB.Where("1=1").Find(&list).Error
	return
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetOfficialAccountList
// @description: 分页获取公众号列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *OfficialAccountService) GetOfficialAccountList(sysUserAuthorityID uint, info request.PageInfo) (list interface{}, total int64, err error) {
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
	var officialAccountList []*ai.OfficialAccount

	db := global.GVA_DB.Model(&ai.OfficialAccount{})
	err = db.Count(&total).Error
	if err != nil {
		return officialAccountList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&officialAccountList).Error
	}

	result := make([]*aiRes.OfficialAccountResponse, 0)

	for _, v := range officialAccountList {
		createTypeList := strings.Split(v.CreateTypes, ",")
		createTypeNumList := make([]int, 0)
		for _, v := range createTypeList {
			num, _ := strconv.Atoi(v)
			createTypeNumList = append(createTypeNumList, num)
		}

		result = append(result, &aiRes.OfficialAccountResponse{
			OfficialAccount: v,
			CreateTypeList:  createTypeNumList,
		})
	}

	return result, total, err
}
