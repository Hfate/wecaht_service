package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
)

type WxTokenService struct {
}

func NewWxTokenService() *WxTokenService {
	return &WxTokenService{}
}

func (*WxTokenService) UpdateWxToken(wxToken ai.WxToken) error {

	dbWxToken := &ai.WxToken{}
	er := global.GVA_DB.Model(&ai.WxToken{}).Where("1=1").Last(&dbWxToken).Error
	if er != nil {
		return er
	}
	dbWxToken.RandInfo = wxToken.RandInfo
	dbWxToken.Token = wxToken.Token
	dbWxToken.DataTicket = wxToken.DataTicket
	dbWxToken.BizUin = wxToken.BizUin
	dbWxToken.SlaveSid = wxToken.SlaveSid

	er = global.GVA_DB.Save(dbWxToken).Error

	return er
}
