package utils

import (
	"bytes"
	"fmt"
	"github.com/ledongthuc/pdf"
	"os"
	"path/filepath"
	"strings"
)

type Article struct {
	Number  string
	Title   string
	Content string
}

// PDF格式的所有文本
func ReadPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

// ReadPDFFiles 函数读取指定目录下的所有PDF文件，并将文件名作为标题，内容作为文章内容。
func ReadPDFFiles(dirPath string) (map[string]string, error) {
	pdfFiles := make(map[string]string)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".pdf" {
			title := filepath.Base(path) // PDF文件名作为标题
			content, _ := ReadPdf(path)
			if err != nil {
				fmt.Printf("Error reading content from %s: %v\n", path, err)
				return nil // 继续处理其他文件
			}

			// 清理内容，移除空格和换行符
			content = strings.TrimSpace(content)
			content = strings.ReplaceAll(content, "\n", " ")

			pdfFiles[title] = content // 将标题和内容添加到映射中
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return pdfFiles, nil
}
