package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

type TopicService struct {
}

//@function: CreateTopic
//@description: 创建主题
//@param: e model.Topic
//@return: err error

func (exa *TopicService) CreateTopic(e ai.Topic) (err error) {
	e.BASEMODEL = BaseModel()
	err = global.GVA_DB.Create(&e).Error
	return err
}

//@function: DeleteFileChunk
//@description: 删除主题
//@param: e model.Topic
//@return: err error

func (exa *TopicService) DeleteTopic(e ai.Topic) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@function: UpdateTopic
//@description: 更新主题
//@param: e *model.Topic
//@return: err error

func (exa *TopicService) UpdateTopic(e *ai.Topic) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

//@function: GetTopic
//@description: 获取主题信息
//@param: id uint
//@return: customer model.Topic, err error

func (exa *TopicService) GetTopic(id uint64) (topic ai.Topic, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&topic).Error
	return
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetTopicList
// @description: 分页获取主题列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *TopicService) GetTopicPage(sysUserAuthorityID uint, info request.PageInfo) (list interface{}, total int64, err error) {
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
	var topicList []ai.Topic

	db := global.GVA_DB.Model(&ai.Topic{})
	err = db.Count(&total).Error
	if err != nil {
		return topicList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("created_at").Find(&topicList).Error
	}
	return topicList, total, err
}

func (exa *TopicService) GetTopicList(sysUserAuthorityID uint, info request.PageInfo) (list interface{}, err error) {
	info.PageSize = 10000

	result := make([]string, 0)
	list, _, err = exa.GetTopicPage(sysUserAuthorityID, info)

	topicList := list.([]ai.Topic)

	for _, item := range topicList {
		result = append(result, item.Topic)
	}
	return result, err
}
