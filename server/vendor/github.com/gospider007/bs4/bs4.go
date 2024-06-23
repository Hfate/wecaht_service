package bs4

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gospider007/re"
	"github.com/gospider007/tools"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// 文档树操作========================================================================= start
type Client struct {
	object  *goquery.Selection
	baseUrl string
}

// 创建一个文档树
func NewClientWithNode(node *html.Node, baseUrl ...string) *Client {
	html := goquery.NewDocumentFromNode(node)
	if html.Size() < 1 {
		return nil
	}
	cli := new(Client)
	if len(baseUrl) > 0 {
		cli.baseUrl = baseUrl[0]
		html.Url, _ = url.Parse(baseUrl[0])
	}
	return cli.newDocument(html.Eq(0))
}
func NewClient(txt string, baseUrl ...string) *Client {
	txt = re.SubFunc(`<\w+(\s+\w+(\s*?=\s*?(".*?"|'.*?'))?)*?\s*?/?>`, func(s string) string {
		if strings.HasSuffix(s, "/>") {
			tag := re.Search(`<(\w+)`, s).Group(1)
			if tag != "br" {
				return re.Sub("/>$", fmt.Sprintf("></%s>", tag), s)
			}
		}
		return s
	}, txt)
	html, err := goquery.NewDocumentFromReader(strings.NewReader(txt))
	if err != nil {
		return nil
	}
	if html.Size() < 1 {
		return nil
	}
	cli := new(Client)
	if len(baseUrl) > 0 {
		cli.baseUrl = baseUrl[0]
		html.Url, _ = url.Parse(baseUrl[0])
	}
	return cli.newDocument(html.Eq(0))
}
func (obj *Client) newDocument(selection *goquery.Selection) *Client {
	return &Client{object: selection, baseUrl: obj.baseUrl}
}

// 寻找一个节点
func (obj *Client) Find(election string) *Client {
	rs := obj.Finds(election)
	if len(rs) > 0 {
		return rs[0]
	}
	return nil
}

// 寻找多个节点
func (obj *Client) Finds(election string) []*Client {
	if obj.Name() == "iframe" {
		return NewClient(obj.Text()).Finds(election)
	}
	clients := obj.finds(election)
	for _, iframe := range obj.finds("iframe") {
		clients = append(clients, iframe.Finds(election)...)
	}
	return clients
}
func (obj *Client) finds(election string) []*Client {
	ll := []*Client{}
	rs := obj.object.Find(election)
	for i := 0; i < rs.Size(); i++ {
		ll = append(ll, obj.newDocument(rs.Eq(i)))
	}
	return ll
}

// 寻找下一个节点
func (obj *Client) Next(elections ...string) *Client {
	var election string
	if len(elections) > 0 {
		election = elections[0]
	}
	if election == "" {
		return obj.newDocument(obj.object.Next())
	} else {
		return obj.newDocument(obj.object.NextFiltered(election))
	}
}

// 寻找之后的所有节点
func (obj *Client) Nexts(elections ...string) []*Client {
	var election string
	if len(elections) > 0 {
		election = elections[0]
	}
	ll := []*Client{}
	var rs *goquery.Selection
	if election == "" {
		rs = obj.object.NextAll()
	} else {
		rs = obj.object.NextAllFiltered(election)
	}
	for i := 0; i < rs.Size(); i++ {
		ll = append(ll, obj.newDocument(rs.Eq(i)))
	}
	return ll
}

// 寻找上一个节点
func (obj *Client) Prev(elections ...string) *Client {
	var election string
	if len(elections) > 0 {
		election = elections[0]
	}

	if election == "" {
		return obj.newDocument(obj.object.Prev())
	} else {
		return obj.newDocument(obj.object.PrevFiltered(election))
	}
}

// 寻找之前的所有节点
func (obj *Client) Prevs(elections ...string) []*Client {
	var election string
	if len(elections) > 0 {
		election = elections[0]
	}
	ll := []*Client{}
	var rs *goquery.Selection
	if election == "" {
		rs = obj.object.PrevAll()
	} else {
		rs = obj.object.PrevAllFiltered(election)
	}
	for i := 0; i < rs.Size(); i++ {
		ll = append(ll, obj.newDocument(rs.Eq(i)))
	}
	return ll
}

// 寻找所有兄弟节点
func (obj *Client) Sibs(elections ...string) []*Client {
	var election string
	if len(elections) > 0 {
		election = elections[0]
	}

	ll := []*Client{}
	var rs *goquery.Selection
	if election == "" {
		rs = obj.object.Siblings()
	} else {
		rs = obj.object.SiblingsFiltered(election)
	}
	for i := 0; i < rs.Size(); i++ {
		ll = append(ll, obj.newDocument(rs.Eq(i)))
	}
	return ll
}

// 寻找所有直接子节点
func (obj *Client) Childrens(elections ...string) []*Client {
	var election string
	if len(elections) > 0 {
		election = elections[0]
	}

	ll := []*Client{}
	var rs *goquery.Selection
	if election == "" {
		rs = obj.object.Children()
	} else {
		rs = obj.object.ChildrenFiltered(election)
	}
	for i := 0; i < rs.Size(); i++ {
		ll = append(ll, obj.newDocument(rs.Eq(i)))
	}
	return ll
}

// 寻找所有子节点
func (obj *Client) ChildrensAll(election ...string) []*Client {
	ll := []*Client{}
	for _, chs := range obj.Childrens(election...) {
		ll = append(ll, chs)
		ll = append(ll, chs.ChildrensAll()...)
	}
	return ll
}

// 寻找所有内容节点
func (obj *Client) Contents(elections ...string) []*Client {
	var election string
	if len(elections) > 0 {
		election = elections[0]
	}
	ll := []*Client{}
	var rs *goquery.Selection
	if election == "" {
		rs = obj.object.Contents()
	} else {
		rs = obj.object.ContentsFiltered(election)
	}
	for i := 0; i < rs.Size(); i++ {
		llt := obj.newDocument(rs.Eq(i))
		if llt.Name() == "#text" && re.Sub(`\s`, "", llt.Text()) == "" {
			continue
		}
		ll = append(ll, llt)
	}
	return ll
}

// 返回所有节点的字符串
func (obj *Client) Texts() []string {
	results := []string{}
	for _, kk := range obj.Contents() {
		results = append(results, kk.Text())
	}
	return results
}

// 返回父节点
func (obj *Client) Parent(elections ...string) *Client {
	var election string
	if len(elections) > 0 {
		election = elections[0]
	}

	if election == "" {
		return obj.newDocument(obj.object.Parent())
	} else {
		return obj.newDocument(obj.object.ParentFiltered(election))
	}
}

// 返回所有父节点
func (obj *Client) Parents(elections ...string) []*Client {
	var election string
	if len(elections) > 0 {
		election = elections[0]
	}

	ll := []*Client{}
	var rs *goquery.Selection
	if election == "" {
		rs = obj.object.Parents()
	} else {
		rs = obj.object.ParentsFiltered(election)
	}
	for i := 0; i < rs.Size(); i++ {
		ll = append(ll, obj.newDocument(rs.Eq(i)))
	}
	return ll
}

// 判断元素是否包含节点
func (obj *Client) Has(obj2 *Client) bool {
	if obj.object.Size() < 1 {
		return false
	} else {
		return obj.object.Contains(obj2.object.Nodes[0])
	}
}

// 在节点中的头部添加节点
func (obj *Client) Prepend(str string) *Client {
	return obj.newDocument(obj.object.PrependHtml(str))
}

// 在节点中的末尾添加节点
func (obj *Client) AppendHtml(str string) *Client {
	return obj.newDocument(obj.object.AppendHtml(str))
}

func (obj *Client) AppendSelection(c *Client) *Client {
	return obj.newDocument(obj.object.AppendSelection(c.object))
}

// 在节点之后添加节点
func (obj *Client) After(str string) *Client {
	return obj.newDocument(obj.object.AfterHtml(str))
}

// 在节点之前添加节点
func (obj *Client) Before(str string) *Client {
	return obj.newDocument(obj.object.BeforeHtml(str))
}

// 替换节点
func (obj *Client) ReplaceWithHtml(str string) *Client {
	return obj.newDocument(obj.object.ReplaceWithHtml(str))
}

// 替换节点
func (obj *Client) ReplaceWithSelection(c *Client) *Client {
	return obj.newDocument(obj.object.ReplaceWithSelection(c.object))
}

// 复制节点
func (obj *Client) Copy() *Client {
	return obj.newDocument(obj.object.Clone())
}

// 删除节点
func (obj *Client) Remove() *Client {
	return obj.newDocument(obj.object.Remove())
}

// 清空节点内容
func (obj *Client) Clear() *Client {
	return obj.newDocument(obj.object.Empty())
}

// 返回节点内容或设置节点内容
func (obj *Client) Text(str ...string) string {
	if len(str) != 0 {
		obj.object.SetText(str[0])
	}
	var buf bytes.Buffer
	// Slightly optimized vs calling Each: no single selection object created
	var f func(*html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			switch n.DataAtom {
			case atom.Input:
				var isText bool
				var textValue string
				for _, attr := range n.Attr {
					if attr.Key == "type" && attr.Val == "text" {
						isText = true
					}
					switch attr.Key {
					case "type":
						if attr.Val == "text" {
							isText = true
						}
					case "value":
						textValue = attr.Val
					}
				}
				if isText {
					buf.WriteString(textValue)
				}
			case atom.Br, atom.P, atom.Li, atom.Ul, atom.Div, atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6,
				atom.Header, atom.Form, atom.Table, atom.Tr, atom.Tbody, atom.Iframe:
				buf.WriteString("\n")
			case atom.Td:
				buf.WriteString(" ")
			}
		case html.TextNode:
			buf.WriteString(strings.TrimSpace(n.Data))
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	for _, n := range obj.object.Nodes {
		f(n)
	}
	return buf.String()
}

// 返回节点名称或设置节点名称
func (obj *Client) Name(str ...string) string {
	if len(str) == 0 {
		return goquery.NodeName(obj.object)
	}
	if len(obj.object.Nodes) > 0 {
		if obj.Name() != str[0] {
			if obj.Name() == "iframe" {
				txt := obj.Text()
				obj.object.Nodes[0].Data = str[0]
				obj.object.Nodes[0].DataAtom = atom.Lookup(tools.StringToBytes(str[0]))
				obj.SetHtml(txt)
			} else {
				obj.object.Nodes[0].Data = str[0]
				obj.object.Nodes[0].DataAtom = atom.Lookup(tools.StringToBytes(str[0]))
			}
		}
	}
	return ""
}

// 返回节点的dom 字符串
func (obj *Client) String() string {
	html, _ := goquery.OuterHtml(obj.object)
	return html
}

// 设置节点的dom
func (obj *Client) SetHtml(content string) {
	obj.object.SetHtml(content)
}

// 获取节点的属性
func (obj *Client) Get(key string, defaultValue ...string) string {
	if len(defaultValue) == 0 {
		val, ok := obj.object.Attr(key)
		if !ok {
			return val
		}
		switch key {
		case "href", "src":
			val2, err := tools.UrlJoin(obj.baseUrl, val)
			if err != nil {
				return val
			}
			return val2
		default:
			return val
		}
	}
	return obj.object.AttrOr(key, defaultValue[0])
}

// 设置节点的属性
func (obj *Client) Set(key string, val string) *Client {
	return obj.newDocument(obj.object.SetAttr(key, val))
}

// 删除节点的属性
func (obj *Client) Del(key string) *Client {
	return obj.newDocument(obj.object.RemoveAttr(key))
}

// 返回节点的所有属性
func (obj *Client) Attrs() map[string]string {
	if len(obj.object.Nodes) == 0 {
		return nil
	}
	result := map[string]string{}
	for _, node := range obj.object.Nodes[0].Attr {
		result[node.Key] = node.Val
	}
	return result
}

// 返回节点的所有属性
func (obj *Client) Nodes() []*html.Node {
	return obj.object.Nodes
}
