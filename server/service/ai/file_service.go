package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"mime/multipart"
)

type FileService struct {
}

var FileServiceApp = new(FileService)

//@function: CreateFile
//@description: 创建素材
//@param: e model.File
//@return: err error

func (exa *FileService) CreateFile(header *multipart.FileHeader) (filePathUrl string, err error) {
	oss := upload.NewOss()

	filePathUrl, _, uploadErr := oss.UploadMultipartFile(header)
	if uploadErr != nil {
		return "", uploadErr
	}
	return filePathUrl, err
}
