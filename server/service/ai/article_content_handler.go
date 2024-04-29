package ai

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"strings"
	"unicode"
)

var OfficialAccountCard = "<section class='mp_profile_iframe_wrp'>" +
	"<mp-common-profile contenteditable='false' class='js_uneditable custom_select_card mp_profile_iframe' " +
	"data-pluginname='mpprofile' data-id='%s' " +
	"data-headimg='%s' " +
	"data-nickname='%s' data-alias='' " +
	"data-signature='%s' " +
	"data-from='0' data-is_biz_ban='0'></mp-common-profile></section>"

var followA = "<section style='margin-bottom: 0px;outline: 0px;font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;background-color: rgb(255, 255, 255);letter-spacing: 0.578px;text-align: center;'><span style='outline: 0px;font-size: 14px;font-family: 'Helvetica Neue', Helvetica, 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;'>" +
	"关注不迷路 随时找得到</span>" +
	"</section>"
var followB = "<section style='margin-bottom: 0px;outline: 0px;font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;background-color: rgb(255, 255, 255);letter-spacing: 0.578px;text-align: center;'><span style='outline: 0px;font-size: 14px;font-family: 'Helvetica Neue', Helvetica, 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;'>" +
	"点赞、关注、转发</span>" +
	"</section>"
var followC = "<section style='margin-bottom: 0px;outline: 0px;font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;background-color: rgb(255, 255, 255);letter-spacing: 0.578px;text-align: center;'><span style='outline: 0px;font-size: 14px;font-family: 'Helvetica Neue', Helvetica, 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;'>" +
	"↓↓↓↓↓↓</span>" +
	"</section>"

type ArticleContentHandler struct {
}

var ArticleContentHandlerApp = new(ArticleContentHandler)

func (ac *ArticleContentHandler) HandleTitle(title string) string {
	title = strings.ReplaceAll(title, "#", "")
	title = strings.ReplaceAll(title, "*", "")
	title = strings.ReplaceAll(title, "标题：", "")
	title = strings.ReplaceAll(title, "#", "")
	title = utils.RemoveBookTitleBrackets(title)
	title = strings.ReplaceAll(title, "标题建议：", "")
	return removeQuotes(title)
}

func removeQuotes(str string) string {
	// 判断字符串是否为空或长度小于2，无法包含前后双引号
	if len(str) < 2 {
		return str
	}

	// 判断首尾字符是否为双引号，并进行去除
	firstChar := rune(str[0])
	lastChar := rune(str[len(str)-1])
	if (firstChar == '"' || firstChar == '“' || firstChar == '”') && (lastChar == '"' || lastChar == '“' || lastChar == '”') {
		return strings.TrimFunc(str, func(r rune) bool {
			return unicode.Is(unicode.Quotation_Mark, r)
		})
	}

	return str
}

func (ac *ArticleContentHandler) Handle(account *ai.OfficialAccount, content string) string {

	// 移除特殊字符
	content = ac.removeSpecialWord(content)

	// 添加推荐阅读
	content = ac.addRecommendedReading(account, content)

	// 富文本
	content = ac.handleRichContent(content, account.CssFormat)

	// 添加公众号排版内容
	content = ac.handleWechatSetType(account, content)

	return content

}

func (ac *ArticleContentHandler) removeSpecialWord(content string) string {
	// 以换行符为分隔符，将文章内容拆分成多行
	lines := strings.Split(content, "\n")

	// 排除标题行
	var contentLines []string
	for _, line := range lines {
		if !strings.Contains(line, "标题：") &&
			!strings.Contains(line, "占位符") &&
			!strings.Contains(line, "配图") &&
			!strings.Contains(line, "小标题") {
			contentLines = append(contentLines, line)
		}
	}

	// 将剩余的行重新连接成一篇文章
	markdownContent := strings.Join(contentLines, "\n")

	//```markdown
	markdownContent = strings.ReplaceAll(markdownContent, "```markdown", "")
	markdownContent = strings.ReplaceAll(markdownContent, "```", "")
	markdownContent = strings.ReplaceAll(markdownContent, "<li><p>", "<li>")
	return markdownContent
}

func (ac *ArticleContentHandler) addRecommendedReading(account *ai.OfficialAccount, mdContent string) string {
	// 获取历史已发布消息5条图文消息
	articleList := ArticleServiceApp.FindLimit5ByPortalName(account.AccountName)

	if len(articleList) == 0 {
		return mdContent
	}

	mdContent += "---\n"
	mdContent += "#### 推荐阅读\n"
	for _, item := range articleList {
		mdContent += "-[" + item.Title + "](" + item.Link + ")\n"
	}

	return mdContent
}

func (ac *ArticleContentHandler) handleWechatSetType(account *ai.OfficialAccount, htmlContent string) string {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		global.GVA_LOG.Error("handleWechatSetType", zap.Error(err))
		return htmlContent
	}

	// 往body中添加一个公众号名片
	accountCard := fmt.Sprintf(OfficialAccountCard, account.AccountId,
		account.HeadImgUrl, account.AccountName, account.Signature)

	body := dom.Find("body")
	body.AppendHtml(followA)
	body.AppendHtml(followB)
	body.AppendHtml(followC)
	body.AppendHtml(accountCard)

	// 输出整个修改后的HTML文档
	modifiedHtml, err := dom.Html()

	return modifiedHtml
}

func (ac *ArticleContentHandler) handleRichContent(content string, cssFormat string) string {
	// 转成html
	htmlContent, err := utils.RenderMarkdownContent(content)
	if err != nil {
		global.GVA_LOG.Error("handleRichContent", zap.Error(err))
		return content
	}

	// 添加css
	if cssFormat != "" {
		htmlContent = utils.UssCssFormat(htmlContent, cssFormat)
	} else {
		htmlContent = utils.HtmlAddStyle(htmlContent)
	}

	return htmlContent
}
