package utils

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
)

type BaiduService struct {
}

func SearchAndSave(keyword string) string {
	imgUrlList := make([]string, 0)

	unsplashImgUrlList := CollectUnsplashImgUrl(keyword)
	if len(unsplashImgUrlList) > 0 {
		imgUrlList = append(imgUrlList, unsplashImgUrlList...)
	}

	baiduImgUrlList := CollectBaiduImgUrl(keyword)
	if len(baiduImgUrlList) > 0 {
		imgUrlList = append(imgUrlList, baiduImgUrlList...)
	}

	// 通过第一张图片链接下载图片
	return saveImage(imgUrlList)
}

func saveImage(imgUrlList []string) string {
	// 通过第一张图片链接下载图片
	filePath := ""

	for _, imgUrl := range imgUrlList {

		filePath = global.GVA_CONFIG.Local.Path + "/" + cast.ToString(GenID()) + ".jpg"
		err := downloadImage(imgUrl, filePath)
		if err != nil {
			global.GVA_LOG.Info("downloadImage failed", zap.Any("err", err.Error()))
		} else {
			break
		}
	}
	return filePath
}

// downloadImage 函数用于下载图片并保存到指定路径
func downloadImage(url string, filePath string) error {
	// 发起 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	// 计算文件大小
	fileSize := resp.ContentLength
	if fileSize > 1*1024*1024 { // 1MB
		return errors.New("图片文件大小超过1MB:" + url)
	}

	// 尝试创建此路径
	mkdirErr := os.MkdirAll(global.GVA_CONFIG.Local.StorePath, os.ModePerm)
	if mkdirErr != nil {
		global.GVA_LOG.Error("function os.MkdirAll() failed", zap.Any("err", mkdirErr.Error()))
		return errors.New("function os.MkdirAll() failed, err:" + mkdirErr.Error())
	}

	// 打开文件用于写入
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将响应体内容复制到文件中
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}