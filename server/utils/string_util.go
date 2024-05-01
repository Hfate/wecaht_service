package utils

import (
	"github.com/russross/blackfriday"
	"strings"
	"unicode"
	"unicode/utf8"
)

var CSSStyleMap = map[string]string{}

func init() {

	CSSStyleMap = make(map[string]string)
	CSSStyleMap["<h1>"] = "<h1 style=\"letter-spacing: normal;text-wrap: wrap;text-align: center;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;font-size: 1.2em;font-weight: bold;display: table;margin: 2em auto 1em;padding-right: 1em;padding-left: 1em;border-bottom: 2px solid rgb(15, 76, 129);color: rgb(63, 63, 63);\">"
	CSSStyleMap["<h2>"] = "<h2 style=\"letter-spacing: normal;text-wrap: wrap;text-align: center;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;font-size: 1.2em;font-weight: bold;display: table;margin: 4em auto 2em;padding-right: 0.2em;padding-left: 0.2em;background: rgb(15, 76, 129);color: rgb(255, 255, 255);\">"
	CSSStyleMap["<h2>"] = "<h2 style=\"letter-spacing: normal;text-wrap: wrap;text-align: center;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;font-size: 1.2em;font-weight: bold;display: table;margin: 4em auto 2em;padding-right: 0.2em;padding-left: 0.2em;background: rgb(15, 76, 129);color: rgb(255, 255, 255);\">"
	CSSStyleMap["<h3>"] = "<h3 style=\"letter-spacing: normal;text-wrap: wrap;text-align: left;line-height: 1.2;font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;font-size: 1.1em;font-weight: bold;margin-top: 2em;margin-right: 8px;margin-bottom: 0.75em;padding-left: 8px;border-left: 3px solid rgb(15, 76, 129);color: rgb(63, 63, 63);\">"
	CSSStyleMap["<h4"] = "<h4 style=\"font-size: 1em;letter-spacing: normal;text-wrap: wrap;text-align: left;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;font-weight: bold;margin: 2em 8px 0.5em;color: rgb(15, 76, 129);\">"
	CSSStyleMap["<p>"] = "<p style=\"line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;font-size: 1em;letter-spacing: 0.1em;color: rgb(80, 80, 80);\">"
	CSSStyleMap["<strong>"] = "<strong style=\"line-height: 1.75;color: rgb(15, 76, 129);\">"
	CSSStyleMap["<ul>"] = "<ul style=\"font-size: 14px;letter-spacing: normal;text-wrap: wrap;text-align: left;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;padding-left: 1em;list-style: circle;color: rgb(63, 63, 63);\" class=\"list-paddingleft-1\">"
	CSSStyleMap["<li>"] = "<li style=\"text-align: left;line-height: 1.75;text-indent: -1em;display: block;margin: 0.2em 8px;\">"
	CSSStyleMap["<span>"] = "<span style=\"line-height: 1.75;color: rgb(87, 107, 149);\">"
}

func RemoveChinese(s string) string {
	result := ""
	for _, r := range s {
		if !unicode.Is(unicode.Scripts["Han"], r) { // 判断字符是否为中文字符
			result += string(r)
		}
	}
	return result
}

func RemoveQuotes(str string) string {
	// 检查字符串长度是否足够
	if len(str) < 2 {
		return str
	}

	// 读取第一个和最后一个字符
	firstRune, firstSize := utf8.DecodeRuneInString(str)
	lastRune, _ := utf8.DecodeLastRuneInString(str)

	// 检查首尾字符是否为双引号
	if (firstRune == '"' || firstRune == '“' || firstRune == '”') && (lastRune == '"' || lastRune == '“' || lastRune == '”') {
		// 去除首尾的双引号
		return strings.TrimFunc(str[firstSize:len(str)-utf8.RuneLen(lastRune)], func(r rune) bool {
			return r == '"' || r == '“' || r == '”'
		})
	}

	return str
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

func HtmlAddStyle(html string) string {
	// 使用Go的字符串替换功能来转义特殊字符
	for k, v := range CSSStyleMap {
		html = strings.ReplaceAll(html, k, v)
	}
	return html
}

func HtmlRemoveStyle(html string) string {
	// 使用Go的字符串替换功能来转义特殊字符
	for k, v := range CSSStyleMap {
		html = strings.ReplaceAll(html, v, k)
	}
	return html
}

func RemoveRepByMap(slc []string) []string {
	var result []string
	tempMap := map[string]struct{}{}
	for _, e := range slc {
		if _, ok := tempMap[e]; !ok {
			tempMap[e] = struct{}{}
			result = append(result, e)
		}
	}
	return result
}
