package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type FileRouter struct {
}

func (e *FileRouter) InitFileRouter(Router *gin.RouterGroup) {
	fileRouter := Router.Group("file").Use(middleware.OperationRecord())
	fileApi := v1.ApiGroupApp.AIApiGroup.FileApi
	{
		fileRouter.POST("upload", fileApi.UploadFile) //上传文件
	}

}
