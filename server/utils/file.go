package utils

import (
	"bytes"
	"errors"
	"github.com/spf13/cast"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func CreateTempImgFile(fileUrl string) (string, error) {
	// 发起 HTTP GET 请求
	resp, err := http.Get(fileUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to download image, status code: " + resp.Status)
	}

	// 将图片内容读入内存
	imageBuffer := new(bytes.Buffer)
	_, err = imageBuffer.ReadFrom(resp.Body)
	if err != nil {
		log.Println("Error reading image:", err)
		return "", err
	}

	tempFileName := cast.ToString(GenID()) + ".jpg"
	// 将内存中的图片数据写入临时文件
	file, err := os.CreateTemp("", tempFileName) // 假设图片格式为 jpg
	if err != nil {
		log.Println("Error creating temp file:", err)
		return "", err
	}
	defer file.Close()

	//写入临时文件
	_, err = file.Write(imageBuffer.Bytes())
	if err != nil {
		log.Println("Error writing to temp file:", err)
		return "", err
	}

	return file.Name(), nil
}

func MultiFile2File(formFile *multipart.FileHeader) (*os.File, error) {
	// 使用 fileHeader.Open() 获取 multipart.File
	file, err := formFile.Open()
	if err != nil {
		// 处理错误
		return nil, err
	}
	defer file.Close() // 记得关闭文件

	tempFileName := cast.ToString(GenID()) + ".jpg"
	// 将内存中的图片数据写入临时文件
	osFile, err := os.CreateTemp("", tempFileName) // 假设图片格式为 jpg
	if err != nil {
		// 处理错误
		return nil, err
	}
	// 将 multipart.File 的内容复制到 os.File 中
	_, err = io.Copy(osFile, file)
	if err != nil {
		// 处理错误
		return nil, err
	}
	return osFile, nil
}
