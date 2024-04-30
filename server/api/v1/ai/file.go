package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FileApi struct{}

// UploadFile
// @Tags      File
// @Summary   创建素材
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     data  body      wechat.File            true  "素材用户名, 素材手机号码"
// @Success   200   {object}  response.Response{msg=string}  "创建素材"
// @Router    /file/upload [post]
func (e *FileApi) UploadFile(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")

	fileUrl, err := fileService.CreateFile(header)

	if err != nil {
		global.GVA_LOG.Error("上传失败!", zap.Error(err))
		response.FileErr(err.Error(), c)
		return
	}

	response.File("https://"+fileUrl, "上传成功", c)
}
