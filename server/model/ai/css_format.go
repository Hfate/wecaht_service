package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type CssFormat struct {
	global.BASEMODEL
	FormatName string `json:"formatName"`
	CssCode    string `json:"cssCode"`
}

func (CssFormat) TableName() string {
	return "wechat_css_format"
}
