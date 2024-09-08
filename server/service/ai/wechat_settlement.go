package ai

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

type WechatSettlementService struct {
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetSettlementList
// @description: 分页获取门户列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *WechatSettlementService) GetSettlementList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	var settlementList []ai.WechatSettlement

	db := global.GVA_DB.Model(&ai.WechatSettlement{})
	err = db.Count(&total).Error
	if err != nil {
		return settlementList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("date desc,account_id desc").Find(&settlementList).Error
	}
	return settlementList, total, err
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetSettlementList
// @description: 分页获取门户列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *WechatSettlementService) Download(c *gin.Context, info request.PageInfo) {
	list, _, err := exa.GetSettlementList(info)
	if err != nil {
		fmt.Println(err)
		return
	}
	articleList := list.([]ai.WechatSettlement)
	result := make([]*ai.WechatSettlementExcel, 0)
	for _, item := range articleList {
		result = append(result, &ai.WechatSettlementExcel{
			AccountName:    item.AccountName,
			Date:           item.Date,
			Zone:           item.Zone,
			Month:          item.Month,
			Order:          item.Order,
			SettStatus:     exa.parseSettleStatus(item.SettStatus),
			SettledRevenue: utils.FloatDiv(cast.ToString(item.SettledRevenue), "100"),
			SettNo:         item.SettNo,
			MailSendCnt:    item.MailSendCnt,
			SlotRevenue:    item.SlotRevenue,
		})
	}

	excelFile := excelize.NewFile()

	utils.WriteDefaultExcelSheet(excelFile, result)

	fileName := "settlement.xlsx"
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

func (exa *WechatSettlementService) parseSettleStatus(settStatus int) string {

	switch settStatus {
	case 1:
		return "结算中"
	case 2:
	case 3:
		return "已结算"
	case 4:
		return "付款中"
	case 5:
		return "已付款"
	default:
		return "-"
	}

	return ""
}
