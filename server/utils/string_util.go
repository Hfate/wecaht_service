package utils

import (
	"github.com/russross/blackfriday"
	"strings"
	"unicode"
)

func RemoveChinese(s string) string {
	result := ""
	for _, r := range s {
		if !unicode.Is(unicode.Scripts["Han"], r) { // 判断字符是否为中文字符
			result += string(r)
		}
	}
	return result
}

func RenderMarkdownContent(markdown string) (string, error) {
	// 渲染Markdown
	html := blackfriday.MarkdownCommon([]byte(markdown))
	return string(html), nil
}

// RemoveBookTitleBrackets 函数用于去除字符串中的书名号
func RemoveBookTitleBrackets(s string) string {
	// 替换书名号《》为空字符串
	temp := strings.ReplaceAll(s, "《", "")   // 去除左侧书名号
	temp = strings.ReplaceAll(temp, "》", "") // 去除右侧书名号
	return temp
}

// EscapeSpecialCharacters 函数接受一个字符串并转义其中的特定特殊字符。
func EscapeSpecialCharacters(input string) string {
	// 使用Go的字符串替换功能来转义特殊字符
	escaped := input
	escaped = strings.ReplaceAll(escaped, "\n", "\\n") // 转义换行符
	escaped = strings.ReplaceAll(escaped, "\t", "\\t") // 转义制表符

	return escaped
}
