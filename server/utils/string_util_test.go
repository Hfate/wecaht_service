package utils

import (
	"fmt"
	"testing"
)

var origin = "https://image.baidu.com/search/acjson?tn=resultjson_com&logid=7812647307136634725&ipn=rj&ct=201326592&is=&fp=result&fr=&word=%E9%A9%AC%E6%96%AF%E5%85%8B&queryWord=%E9%A9%AC%E6%96%AF%E5%85%8B&cl=2&lm=-1&ie=utf-8&oe=utf-8&adpicid=&st=-1&z=&ic=0&hd=&latest=&copyright=&s=&se=&tab=&width=&height=&face=0&istype=2&qc=&nc=1&expermode=&nojc=&isAsync=&pn=1&rn=1&gsm=3c&1711641530617="

func TestRemoveBrackets(t *testing.T) {
	resutl, _ := NormalGetStr(origin)
	fmt.Println(resutl)
}

func TestRenderMarkdownContent(t *testing.T) {

	got, err := RenderMarkdownContent("关于“台当局被批劫贫济富助绿金权贵”的文章，其主题敏感且具有争议性。正如您所说，我们在讨论这类话题时，需要保持客观和理性，避免传播不实信息。考虑到这一点，我将尝试以一种更加中立和建设性的方式来探讨这个话题，并在适当的位置添加配图占位符以增强文章的视觉效果。\n\n台湾，作为中国不可分割的一部分，其社会、政治动态一直备受关注。近期，有关台当局某些政策的争议不断升温，尤其是涉及经济分配和权贵利益的问题，更是引发了民众的广泛关注。\n\n\n<img data-s=\\\"300,640\\\" data-galleryid=\\\"\\\" data-type=\\\"png\\\"  class=\\\"rich_pages wxw-img\\\" data-src=\\\"http://mmbiz.qpic.cn/sz_mmbiz_png/uO29ibicRxJ0Q0q1oP0raoHEGuUJv74Lt9aeM1LnlDhzrxt3Wz0LaKQ5MGzgVDLk0ILE03JHGWeThraRSy7VE5yw/0?wx_fmt=png\\\" style=\\\"\\\" data-ratio=\\\"0.8264840182648402\\\" data-w=\\\"438\\\">\n\n\n在批评声音中，有人指责台当局的政策偏向于“劫贫济富”，即牺牲普通民众的利益来满足少数权贵的需求。这种指责是否成立，需要我们深入分析相关政策的制定背景、实施效果以及受益群体。\n\n然而，不论争议如何，一个不容忽视的事实是，台湾经济的发展与民众的生活息息相关。任何政策的出台，都应该以增进民众福祉为出发点和落脚点。\n\n\n<img data-s=\\\"300,640\\\" data-galleryid=\\\"\\\" data-type=\\\"jpeg\\\"  class=\\\"rich_pages wxw-img\\\" data-src=\\\"http://mmbiz.qpic.cn/sz_mmbiz_jpg/uO29ibicRxJ0Q0q1oP0raoHEGuUJv74Lt9KibicvBy7ESqlCQYTtgibCMZmPJ0FicryCnPULHgIs1XMQ7mo4yfyBcoIA/0?wx_fmt=jpeg\\\" style=\\\"\\\" data-ratio=\\\"0.8264840182648402\\\" data-w=\\\"438\\\">\n\n\n中国政府一直致力于维护国家主权和领土完整，这其中包括对台湾地区的关注和重视。在两岸关系上，中国政府始终坚持和平发展、互利共赢的原则，推动两岸经济文化交流合作，增进两岸同胞的相互了解和信任。\n\n面对台湾当局的政策和做法，我们作为普通民众，有责任也有权利通过合法渠道表达意见和建议。在表达观点时，我们应该保持理性和客观，避免使用攻击性和煽动性的言辞。\n\n\n<img data-s=\\\"300,640\\\" data-galleryid=\\\"\\\" data-type=\\\"jpeg\\\"  class=\\\"rich_pages wxw-img\\\" data-src=\\\"http://mmbiz.qpic.cn/sz_mmbiz_jpg/uO29ibicRxJ0Q0q1oP0raoHEGuUJv74Lt9ZiaTjfzCxnAMldaWK03zWGQ0UeDPuC6AYVNNTTZjdicgRu75zrBVhXDw/0?wx_fmt=jpeg\\\" style=\\\"\\\" data-ratio=\\\"0.8264840182648402\\\" data-w=\\\"438\\\">\n\n\n最后，无论是对于台湾问题还是其他社会争议话题，我们都应该以建设性的态度来参与讨论。通过理性对话和协商解决问题，才是推动社会进步和发展的重要途径。让我们共同努力，为两岸关系的和平发展贡献一份力量。")
	fmt.Println(got, err)
}

func TestHtmlAddStyle(t *testing.T) {

	html := "<h1 style=\"letter-spacing: normal;text-wrap: wrap;text-align: center;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1.2em;font-weight: bold;display: table;margin: 2em auto 1em;padding-right: 1em;padding-left: 1em;border-bottom: 2px solid rgb(15, 76, 129);color: rgb(63, 63, 63);\">2025国考新风向：半月谈60天备考计划全面升级</h1>\n\n<h2 style=\"letter-spacing: normal;text-wrap: wrap;text-align: center;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1.2em;font-weight: bold;display: table;margin: 4em auto 2em;padding-right: 0.2em;padding-left: 0.2em;background: rgb(15, 76, 129);color: rgb(255, 255, 255);\">引言</h2>\n\n<p style=\"line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1em;letter-spacing: 0.1em;color: rgb(80, 80, 80);\">随着2025年国考季的临近，备考的热潮再次掀起。半月谈公考系列的《60天上岸计划》迎来了它的第七个年头，今年更是带来了全面升级改版的好消息。今天，我们就来聊聊这个备受瞩目的备考计划，以及它为考生们带来的新机遇。</p>\n\n<h2 style=\"letter-spacing: normal;text-wrap: wrap;text-align: center;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1.2em;font-weight: bold;display: table;margin: 4em auto 2em;padding-right: 0.2em;padding-left: 0.2em;background: rgb(15, 76, 129);color: rgb(255, 255, 255);\">国考备考新动态</h2>\n\n<h3 style=\"letter-spacing: normal;text-wrap: wrap;text-align: left;line-height: 1.2;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1.1em;font-weight: bold;margin-top: 2em;margin-right: 8px;margin-bottom: 0.75em;padding-left: 8px;border-left: 3px solid rgb(15, 76, 129);color: rgb(63, 63, 63);\">1. 全新升级的配套讲义</h3>\n\n<ul>\n<li><strong style=\"line-height: 1.75;color: rgb(15, 76, 129);\">版本更新</strong>：《2025国考60天上岸计划》迎来了它的第七个版本，全面升级改版，为考生提供更全面的备考支持。</li>\n<li><strong style=\"line-height: 1.75;color: rgb(15, 76, 129);\">内容开箱</strong>：新版讲义已经寄达，让我们来看看它究竟带来了哪些新变化。</li>\n</ul>\n\n<h3 style=\"letter-spacing: normal;text-wrap: wrap;text-align: left;line-height: 1.2;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1.1em;font-weight: bold;margin-top: 2em;margin-right: 8px;margin-bottom: 0.75em;padding-left: 8px;border-left: 3px solid rgb(15, 76, 129);color: rgb(63, 63, 63);\">2. 灵活的学习方式</h3>\n\n<ul>\n<li><strong style=\"line-height: 1.75;color: rgb(15, 76, 129);\">课程解锁</strong>：购买课程后，考生可以立即解锁行测+申论的全套录播系统课，无需打卡，支持无限回看，自主掌握学习节奏。</li>\n<li><strong style=\"line-height: 1.75;color: rgb(15, 76, 129);\">物流同步</strong>：根据考生在半月谈APP上填写的地址，配套讲义将被寄送，同时后台会同步物流信息，确保考生及时收到学习资料。</li>\n</ul>\n\n<h3 style=\"letter-spacing: normal;text-wrap: wrap;text-align: left;line-height: 1.2;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1.1em;font-weight: bold;margin-top: 2em;margin-right: 8px;margin-bottom: 0.75em;padding-left: 8px;border-left: 3px solid rgb(15, 76, 129);color: rgb(63, 63, 63);\">3. 系统化的备考资料</h3>\n\n<ul>\n<li><strong style=\"line-height: 1.75;color: rgb(15, 76, 129);\">分册学习</strong>：讲义分为三册，分别对应申论、行测文、行测理，与视频课程配合使用，提高学习效率。</li>\n<li><strong style=\"line-height: 1.75;color: rgb(15, 76, 129);\">申论范文</strong>：独家的30篇五色申论范文，改编自半月谈时评文章，配备多巴胺配色解析，助力考生提升写作技能。</li>\n</ul>\n\n<h3 style=\"letter-spacing: normal;text-wrap: wrap;text-align: left;line-height: 1.2;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1.1em;font-weight: bold;margin-top: 2em;margin-right: 8px;margin-bottom: 0.75em;padding-left: 8px;border-left: 3px solid rgb(15, 76, 129);color: rgb(63, 63, 63);\">4. 实战演练与复盘</h3>\n\n<ul>\n<li><strong style=\"line-height: 1.75;color: rgb(15, 76, 129);\">真题练习</strong>：提供17套国考行测和申论真题，帮助考生熟悉考试题型和难度。</li>\n<li><strong style=\"line-height: 1.75;color: rgb(15, 76, 129);\">一对一测评</strong>：考生可以凭借购买记录，享受一对一的测评服务，由谈哥亲自指导。</li>\n<li><strong style=\"line-height: 1.75;color: rgb(15, 76, 129);\">复盘表</strong>：首发的半月谈公考行测、申论刷题复盘表，由资深讲师和高分学员共同研发，提供有效的复盘方法。</li>\n</ul>\n\n<h3 style=\"letter-spacing: normal;text-wrap: wrap;text-align: left;line-height: 1.2;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1.1em;font-weight: bold;margin-top: 2em;margin-right: 8px;margin-bottom: 0.75em;padding-left: 8px;border-left: 3px solid rgb(15, 76, 129);color: rgb(63, 63, 63);\">5. 额外福利</h3>\n\n<ul>\n<li><strong style=\"line-height: 1.75;color: rgb(15, 76, 129);\">杂志赠送</strong>：首发期间购买的考生，还将额外获得2024年全年的《半月谈》杂志，为时政学习和写作提供素材。</li>\n</ul>\n\n<h2 style=\"letter-spacing: normal;text-wrap: wrap;text-align: center;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1.2em;font-weight: bold;display: table;margin: 4em auto 2em;padding-right: 0.2em;padding-left: 0.2em;background: rgb(15, 76, 129);color: rgb(255, 255, 255);\">价格与购买</h2>\n\n<ul>\n<li><strong style=\"line-height: 1.75;color: rgb(15, 76, 129);\">价格调整</strong>：2025季的《60天上岸计划》调整为非打卡产品，价格下调至680元，提供更经济的备考选择。</li>\n</ul>\n\n<h2 style=\"letter-spacing: normal;text-wrap: wrap;text-align: center;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1.2em;font-weight: bold;display: table;margin: 4em auto 2em;padding-right: 0.2em;padding-left: 0.2em;background: rgb(15, 76, 129);color: rgb(255, 255, 255);\">展望与预测</h2>\n\n<p style=\"line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1em;letter-spacing: 0.1em;color: rgb(80, 80, 80);\">随着国考竞争的日益激烈，备考资源的质量和多样性成为了考生们关注的焦点。半月谈的这次全面升级，无疑为考生们提供了一个更加系统、灵活且经济的备考方案。预计这将吸引更多考生的关注，并可能引领公考备考的新趋势。</p>\n\n<h2 style=\"letter-spacing: normal;text-wrap: wrap;text-align: center;line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1.2em;font-weight: bold;display: table;margin: 4em auto 2em;padding-right: 0.2em;padding-left: 0.2em;background: rgb(15, 76, 129);color: rgb(255, 255, 255);\">结语</h2>\n\n<p style=\"line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1em;letter-spacing: 0.1em;color: rgb(80, 80, 80);\">备考国考是一场持久战，而一个好的备考计划就像是一盏明灯，照亮前行的路。半月谈的《2025国考60天上岸计划》全面升级，无疑为考生们提供了一个全新的起点。让我们一起期待，这将如何助力考生们在国考的征途上乘风破浪。</p>\n\n<hr />\n\n<p style=\"line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1em;letter-spacing: 0.1em;color: rgb(80, 80, 80);\"><br><p style=\"line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1em;letter-spacing: 0.1em;color: rgb(80, 80, 80);\"><img src=\"http://mmbiz.qpic.cn/sz_mmbiz_jpg/uO29ibicRxJ0QibM4iaN36DhfCfQ8D6k5tOEqWicRf1QHXKueqfrdE3xxhHokplkXN4KvO0QtQYfIHA8F00Qqbiaib0VA/0?from=appmsg\"></p><br></p>\n\n<p style=\"line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1em;letter-spacing: 0.1em;color: rgb(80, 80, 80);\"><br><p style=\"line-height: 1.75;font-family: -apple-system-font, BlinkMacSystemFont, \"Helvetica Neue\", \"PingFang SC\", \"Hiragino Sans GB\", \"Microsoft YaHei UI\", \"Microsoft YaHei\", Arial, sans-serif;font-size: 1em;letter-spacing: 0.1em;color: rgb(80, 80, 80);\"><img src=\"http://mmbiz.qpic.cn/sz_mmbiz_jpg/uO29ibicRxJ0QibM4iaN36DhfCfQ8D6k5tOEaWdRN6VWo2jaKM1GlfZpHFX8ggMK4jbu7WEUy5mCfCmLRDmycFJCzg/0?from=appmsg\"></p><br></p>\n\n"
	fmt.Println(HtmlAddStyle(html))

}

func TestRemoveQuotes(t *testing.T) {
	fmt.Println(RemoveQuotes("“当爱已成往事：揭秘出轨男人的真心话”"))
}

func TestRemoveMarkedSections(t *testing.T) {
	// 示例文章
	article := `开头（200字）
李婷站在咖啡厅的门口，手里紧握着一杯已经冷却的拿铁...
发展（500字）
李婷和张涛的故事，始于大学校园...
顶峰（300字）
然而，就在李婷以为自己已经走出了张涛的阴影时...
结尾（200字）
李婷走在回家的路上，她的心情异常平静。`

	// 移除标记及其后面的内容
	modifiedArticle := RemoveSections(article)

	// 输出结果
	fmt.Println(modifiedArticle)
}

func TestInsertLineBreaks(t *testing.T) {
	// 示例HTML内容，实际应用中应从提供的网页源代码中获取
	html := `<p>这是第一段。</p><p>这是第二段。</p>`

	// 修正排版
	fixedHTML := InsertLineBreaks(html)

	// 打印修正后的HTML内容
	println(fixedHTML)
}

func TestRenderMarkdownContent1(t *testing.T) {
	str := "---\\n\\n# 工作与生活：寻找平衡的艺术\\n\\n在这个快节奏的时代，工作与生活的平衡成为了许多人心中的难题。对于中年人来说，这个问题尤为突出，因为他们不仅要面对职场的挑战，还要承担起家庭的责任。然而，真正的生活不应该只是工作的延伸，而是一个独立而丰富的存在。\\n\\n## 工作的意义：超越薪水的追求\\n\\n工作，不仅仅是为了薪水。它应该是一种实现自我价值、追求个人成长的方式。对于饱经沧桑的中年人来说，工作的意义更在于它能否带来内心的满足和成就感。**“一个人主业是好好生活，副业才是工作。”**这句话提醒我们，工作虽然重要，但它不应该占据我们生活的全部。\\n\\n## 休假的价值：生活的调味品\\n\\n休假，是生活中不可或缺的一部分。它让我们有机会暂时放下工作的重担，去享受生活的美好。然而，许多人却失去了对休假的兴趣，将其视为另一种形式的工作。这无疑是一种遗憾。休假应该是一个放松身心、充实自我的过程，而不是简单的休息。\\n\\n## 工作与热爱：寻找激情的火花\\n\\n对于工作，我们是否真正热爱它？我们的工作是否真正有意义？这些问题值得每个人深思。对于那些重复而枯燥的工作，我们是否应该寻找改变？对于那些能够激发我们激情的工作，我们是否应该更加投入？\\n\\n## 富人与穷人：行动力的差距\\n\\n穷人与富人之间的差距，并不在于他们是否有工作，而在于他们的行动力。有工作的人不一定具备工作能力，同样，没有工作的人也可能拥有强大的工作能力。关键在于，我们是否有勇气去行动，去改变现状。\\n\\n## 信息茧房：思维的局限\\n\\n信息茧房，是指人们只关注和接触自己感兴趣的信息，从而形成一种思维的局限。这种现象在当今社会尤为普遍。我们应该如何打破这种局限，拓宽自己的视野？这需要我们不断学习，不断尝试新的事物。\\n\\n## 刻意研究赚钱：财富的修行\\n\\n赚钱，是一个人最好的修行。它不仅需要智慧，更需要行动力。那些拥有财富的人，他们的付出和努力远超我们的想象。他们知道，要实现财务自由，需要有战略性的思考和规划。\\n\\n"
	fmt.Println(RenderMarkdownContent(str))
}
