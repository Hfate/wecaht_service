package utils

import (
	"github.com/russross/blackfriday"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var CSSStyleMap = map[string]string{}

func init() {

	CSSStyleMap = make(map[string]string)
	//CSSStyleMap["<h1>"] = "<h1 style=\"letter-spacing: normal;text-wrap: wrap;text-align: center;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;font-size: 1.2em;font-weight: bold;display: table;margin: 2em auto 1em;padding-right: 1em;padding-left: 1em;border-bottom: 2px solid rgb(15, 76, 129);color: rgb(63, 63, 63);\">"
	//CSSStyleMap["<h2>"] = "<h2 style=\"letter-spacing: normal;text-wrap: wrap;text-align: center;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;font-size: 1.2em;font-weight: bold;display: table;margin: 4em auto 2em;padding-right: 0.2em;padding-left: 0.2em;background: rgb(15, 76, 129);color: rgb(255, 255, 255);\">"
	//CSSStyleMap["<h2>"] = "<h2 style=\"letter-spacing: normal;text-wrap: wrap;text-align: center;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;font-size: 1.2em;font-weight: bold;display: table;margin: 4em auto 2em;padding-right: 0.2em;padding-left: 0.2em;background: rgb(15, 76, 129);color: rgb(255, 255, 255);\">"
	//CSSStyleMap["<h3>"] = "<h3 style=\"letter-spacing: normal;text-wrap: wrap;text-align: left;line-height: 1.2;font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;font-size: 1.1em;font-weight: bold;margin-top: 2em;margin-right: 8px;margin-bottom: 0.75em;padding-left: 8px;border-left: 3px solid rgb(15, 76, 129);color: rgb(63, 63, 63);\">"
	//CSSStyleMap["<h4>"] = "<h4 style=\"font-size: 1em;letter-spacing: normal;text-wrap: wrap;text-align: left;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;font-weight: bold;margin: 2em 8px 0.5em;color: rgb(15, 76, 129);\">"
	CSSStyleMap["<p>"] = "<p style=\"line-height: 2;margin-bottom: 24px;margin-left: 8px;margin-right: 8px;\">"
	CSSStyleMap["<strong>"] = "<strong style=\"line-height: 2;color: rgb(15, 76, 129);\">"
	//CSSStyleMap["<ul>"] = "<ul style=\"color: rgb(63, 63, 63);\" >"
	//CSSStyleMap["<li>"] = "<li style=\"text-align: left;line-height: 1.75;text-indent: -1em;display: block;margin: 0.2em 8px;\">"
	CSSStyleMap["<span>"] = "<span style=\"line-height: 2;color: rgb(87, 107, 149);\">"
}

// InsertLineBreaks 会在 <p> 标签之间没有换行的情况下插入 <br> 标签
func InsertLineBreaks(html string) string {

	html = strings.ReplaceAll(html, "</p><p", "</p><br><p")

	return html
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

// RemoveSections 函数接受一个包含文章的字符串，并移除其中的特定段落。
func RemoveSections(text string) string {
	// 定义要移除的段落标记
	markers := []string{
		"开头（\\d+字)",
		"发展 (\\d+字)",
		"顶峰 (\\d+字)",
		// 可以继续添加更多的标记
	}

	// 构建用于匹配包含特定标记的行的正则表达式
	markerPattern := regexp.QuoteMeta(strings.Join(markers, "|"))
	regex := regexp.MustCompile("(?m)^(" + markerPattern + ")(\\s+.*)?$")

	// 使用正则表达式替换掉包含特定标记的行
	cleanedText := regex.ReplaceAllString(text, "")

	return cleanedText
}

// InsertTextAtThirds 在文本的1/3和2/3位置插入指定的文本
func InsertTextAtThirds(text string, insertText1, insertText2 string) string {
	// 将原始文本按行分割
	lines := strings.Split(text, "\n")
	// 计算总行数
	totalLines := len(lines)
	// 计算1/3和2/3的位置
	oneThird := totalLines / 3
	twoThirds := oneThird * 2

	// 在1/3和2/3的位置插入文本
	if oneThird > 0 {
		lines = append(lines[:oneThird], append([]string{insertText1}, lines[oneThird:]...)...)
	}
	if twoThirds > 0 && twoThirds < totalLines {
		lines = append(lines[:twoThirds], append([]string{insertText2}, lines[twoThirds:]...)...)
	}

	// 将行重新组合成字符串，并用换行符连接
	return strings.Join(lines, "\n")
}

var removeWords = []string{"标题：", "原创性", "二创", "Prompt", "占位符", "原文素材", "配图", "小标题", "正文"}

func RemoveSpecialWord(content string) string {
	// 以换行符为分隔符，将文章内容拆分成多行
	lines := strings.Split(content, "\n")

	// 排除标题行
	var contentLines []string
	for _, line := range lines {
		isContinue := false
		for _, item := range removeWords {
			if strings.Contains(line, item) {
				isContinue = true
				break
			}
		}

		if isContinue {
			continue
		}

		contentLines = append(contentLines, line)
	}

	// 将剩余的行重新连接成一篇文章
	markdownContent := strings.Join(contentLines, "\n")

	//```markdown
	markdownContent = strings.ReplaceAll(markdownContent, "```markdown", "")
	markdownContent = strings.ReplaceAll(markdownContent, "```", "")
	markdownContent = strings.ReplaceAll(markdownContent, "<li><p>", "<li>")
	return markdownContent
}

func RemoveNonsense(content string) string {
	contentArr := strings.Split(content, "---")
	maxLength := 0
	result := ""
	for _, item := range contentArr {
		if len(item) > maxLength {
			result = item
			maxLength = len(item)
		}
	}
	return result
}
