package ai

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"strings"
)

var OfficialAccountCard = "<section class='mp_profile_iframe_wrp'>" +
	"<mp-common-profile contenteditable='false' class='js_uneditable custom_select_card mp_profile_iframe' " +
	"data-pluginname='mpprofile' data-id='%s' " +
	"data-headimg='%s' " +
	"data-nickname='%s' data-alias='' " +
	"data-signature='%s' " +
	"data-from='0' data-is_biz_ban='0'></mp-common-profile></section>"

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
	return utils.RemoveQuotes(title)
}

func (ac *ArticleContentHandler) Handle(account *ai.OfficialAccount, content string) string {

	// 移除特殊字符
	content = ac.removeSpecialWord(content)

	// 移除文首文末废话
	content = ac.removeNonsense(content)

	//使用模板
	content = ac.useTemplate(account.AppId, content)

	// 添加推荐阅读
	content = ac.addRecommendedReading(account, content)

	// 添加公众号排版内容
	content = ac.addCard(account, content)

	// 富文本
	content = ac.addCss(content, account.CssFormat)

	return content

}

func (ac *ArticleContentHandler) useTemplate(appId, mdContent string) string {
	template := TemplateServiceApp.FindByAccountId(appId)

	htmlContent := template.TemplateValue
	htmlContent = strings.ReplaceAll(htmlContent, "<p style=\"text-align: left;\">{{.Content}}</p>", "{{.Content}}")

	md2Html, _ := utils.RenderMarkdownContent(mdContent)

	htmlContent = strings.ReplaceAll(htmlContent, "{{.Content}}", md2Html)

	return htmlContent
}

func (ac *ArticleContentHandler) removeNonsense(content string) string {
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

var removeWords = []string{"标题：", "原创性", "二创", "Prompt", "占位符", "原文素材", "配图", "小标题", "正文"}

func (ac *ArticleContentHandler) removeSpecialWord(content string) string {
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

func (ac *ArticleContentHandler) addRecommendedReading(account *ai.OfficialAccount, htmlContent string) string {
	if !strings.Contains(htmlContent, "{{.RecommendList}}") {
		return htmlContent
	}

	htmlContent = strings.ReplaceAll(htmlContent, "<p style=\"text-align: left;\">{{.RecommendList}}</p>", "{{.RecommendList}}")

	// 获取历史已发布消息5条图文消息
	articleList := ArticleServiceApp.FindLimit5ByPortalName(account.AccountName)

	if len(articleList) == 0 {
		htmlContent = strings.ReplaceAll(htmlContent, "{{.RecommendList}}", "")
		return htmlContent
	}

	titleSet := make(map[string]bool)
	recommendList := "<p>"

	for _, item := range articleList {
		if titleSet[item.Title] {
			continue
		}
		titleSet[item.Title] = true
		recommendList += "<a href='" + item.Link + "'>" + item.Title + "</a><br>"
	}
	recommendList += "</p>"

	htmlContent = strings.ReplaceAll(htmlContent, "{{.RecommendList}}", recommendList)

	return htmlContent
}

func (ac *ArticleContentHandler) addCard(account *ai.OfficialAccount, htmlContent string) string {
	if !strings.Contains(htmlContent, "{{.Card}}") {
		return htmlContent
	}

	htmlContent = strings.ReplaceAll(htmlContent, "<p style=\"text-align: left;\">{{.Card}}</p>", "{{.Card}}")

	accountCard := fmt.Sprintf(OfficialAccountCard, account.AccountId,
		account.HeadImgUrl, account.AccountName, account.Signature)

	htmlContent = strings.ReplaceAll(htmlContent, "{{.Card}}", accountCard)

	return htmlContent
}

func (ac *ArticleContentHandler) addCss(htmlContent string, cssFormat string) string {
	// 添加css
	if cssFormat != "" {
		htmlContent = utils.UssCssFormat(htmlContent, cssFormat)
	} else {
		htmlContent = utils.HtmlAddStyle(htmlContent)
	}
	return htmlContent
}
