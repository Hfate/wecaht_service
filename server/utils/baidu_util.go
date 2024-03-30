package utils

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"os"
	"strings"
)

type BaiduService struct {
}

func SearchAndSave(keyword string) string {
	imgUrlList := CollectBaiduImgUrl(keyword)

	// 通过第一张图片链接下载图片
	return saveImage(imgUrlList)
}

func saveImage(imgUrlList []string) string {
	// 通过第一张图片链接下载图片
	filePath := ""

	for _, imgUrl := range imgUrlList {
		// 解析图片链接中的参数
		strArr := strings.Split(imgUrl, ".")

		filePath = global.GVA_CONFIG.Local.Path + "/" + cast.ToString(GenID()) + "." + strArr[len(strArr)-1]
		err := downloadImage(imgUrl, filePath)
		if err != nil {
			fmt.Println(err)
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
