package utils

import (
	"bufio"
	"bytes"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/tdewolff/minify/v2"
	"go.uber.org/zap"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	mhtml "github.com/tdewolff/minify/v2/html"
	"github.com/vanng822/go-premailer/premailer"
)

func WriteCodeCss(theme *chroma.Style) string {
	// write css
	hlbuf := bytes.Buffer{}
	hlw := bufio.NewWriter(&hlbuf)
	formatter := html.New(html.WithClasses(true))
	if err := formatter.WriteCSS(hlw, theme); err != nil {
		panic(err)
	}
	hlw.Flush()
	return hlbuf.String()
}
func ReplaceCodeParts(doc *goquery.Document) (string, error) {
	// find code-parts via selector and replace them with highlighted versions
	var hlErr error
	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
		if hlErr != nil {
			return
		}
		class, _ := s.Attr("class")
		lang := strings.TrimPrefix(class, "language-")
		oldCode := s.Text()
		lexer := lexers.Get(lang)
		formatter := html.New(html.WithClasses(true))
		iterator, err := lexer.Tokenise(nil, string(oldCode))
		if err != nil {
			hlErr = err
			return
		}
		b := bytes.Buffer{}
		buf := bufio.NewWriter(&b)
		if err := formatter.Format(buf, styles.GitHub, iterator); err != nil {
			hlErr = err
			return
		}
		if err := buf.Flush(); err != nil {
			hlErr = err
			return
		}
		s.SetHtml(b.String())
	})
	if hlErr != nil {
		return "", hlErr
	}
	new, err := doc.Html()
	if err != nil {
		return "", err
	}
	// replace unnecessarily added html tags
	return new, nil
}

func AddHtmlTag(input string) string {
	//拼凑HTML页面，需要先导入Strings包
	s1 := "<html><head><meta charset=\"UTF-8\"><title></title></head><body>"
	s2 := "</body></html>"
	var build strings.Builder
	build.WriteString(s1)
	build.WriteString(input)
	build.WriteString(s2)
	s3 := build.String()
	return s3
}

func ParseInlineCss(content string) string {
	prem, err := premailer.NewPremailerFromString(content, premailer.NewOptions())
	if err != nil {
		global.GVA_LOG.Error("ParseInlineCss", zap.Error(err))
	}

	html, err := prem.Transform()
	if err != nil {
		global.GVA_LOG.Error("ParseInlineCss", zap.Error(err))
	}
	return html
}
func RepImage(htmls string) string {
	var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	imgs := imgRE.FindAllStringSubmatch(htmls, -1)
	return imgs[0][1]
}
func ChangeLine(content string) string {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		panic(err)
	}
	dom.Find("pre>code").Each(func(i int, selection *goquery.Selection) {
		textstring := selection.Text()
		textstring = strings.Replace(textstring, "\n", "g^g+;", -1)
		selection.SetText(textstring)
	})
	str, _ := dom.Html()
	//移除所有换行
	str = strings.Replace(str, "\n", "", -1)
	return str
}
func RemoveLine(dom *goquery.Document) *goquery.Document {
	dom.Find("pre>code").Each(func(i int, selection *goquery.Selection) {
		textstring := selection.Text()
		textstring = strings.Replace(textstring, "g^g+;", "\n", -1)
		selection.SetText(textstring)
	})
	return dom
}
func UssCssFormat(htmlContent string, cssFormat string) string {

	content := AddHtmlTag(htmlContent)

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		panic(err)
	}
	// dom = RemoveLine(dom)
	//将代码块部分修改
	ReplaceCodeParts(dom)

	//代码高亮css
	code_css := WriteCodeCss(styles.MonokaiLight)

	//正常css
	dom.Find("title").Each(func(i int, selection *goquery.Selection) {
		selection.AfterHtml("<style>" + code_css + "\n" + cssFormat + "</style>")
	})

	dom.Find("br:not(code)").Each(func(i int, selection *goquery.Selection) {
		selection.Remove()
	})

	dom_content, _ := dom.Html()
	parse_inline_html := ParseInlineCss(dom_content)

	parsedom, err := goquery.NewDocumentFromReader(strings.NewReader(parse_inline_html))
	if err != nil {
		panic(err)
	}
	parsedom.Find("style").Each(func(i int, selection *goquery.Selection) {
		selection.Remove()
	})

	str, _ := parsedom.Html()
	return str

}

// 压缩html文件
func HtmlMinifyFile(filename string) (string, error) {

	m := minify.New()
	m.Add("text/html", &mhtml.Minifier{
		KeepDefaultAttrVals: true,
		KeepDocumentTags:    true,
		KeepEndTags:         true,
	})

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	mb, err := m.String("text/html", string(b))
	if err != nil {
		return "", err
	}

	return mb, err

}
func HtmlMinify(htmlstring string) (string, error) {

	m := minify.New()
	m.Add("text/html", &mhtml.Minifier{
		KeepDefaultAttrVals: true,
		KeepDocumentTags:    true,
		KeepEndTags:         true,
	})

	mb, err := m.String("text/html", htmlstring)
	if err != nil {
		return "", err
	}

	return mb, err

}
func IsHttp(text string) bool {
	myRegex, _ := regexp.Compile("^(http|https)://")
	found := myRegex.FindStringIndex(text)
	if found == nil {
		return false
	} else {
		return true
	}

}
