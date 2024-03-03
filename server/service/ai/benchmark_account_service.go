package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/wechat/request"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"time"
)

type BenchmarkAccountService struct {
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreatePortal
//@description: 创建门户
//@param: e model.Portal
//@return: err error

func (exa *BenchmarkAccountService) CreateBenchmarkAccount(e wechat.BenchmarkAccount) (err error) {
	e.GVA_MODEL = global.GVA_MODEL{ID: utils.GenID(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err = global.GVA_DB.Create(&e).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteFileChunk
//@description: 删除对标账号
//@param: e model.BenchmarkAccount
//@return: err error

func (exa *BenchmarkAccountService) DeleteBenchmarkAccount(e wechat.BenchmarkAccount) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetBenchmarkAccount
//@description: 获取对标账号信息
//@param: id uint
//@return: customer model.BenchmarkAccount, err error

func (exa *BenchmarkAccountService) GetBenchmarkAccount(id uint64) (portal wechat.BenchmarkAccount, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&portal).Error
	return
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetBenchmarkAccountList
// @description: 分页获取对标账号列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *BenchmarkAccountService) GetBenchmarkAccountList(sysUserAuthorityID uint, info aiReq.BenchmarkAccountSearch) (list interface{}, total int64, err error) {
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
	var benchmarkAccountList []wechat.BenchmarkAccount

	db := global.GVA_DB.Model(&wechat.BenchmarkAccount{})

	// 如果有条件搜索 下方会自动创建搜索语句
	if info.AccountName != "" {
		db = db.Where("account_name LIKE ?", "%"+info.AccountName+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return benchmarkAccountList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Find(&benchmarkAccountList).Error
	}
	return benchmarkAccountList, total, err
}
