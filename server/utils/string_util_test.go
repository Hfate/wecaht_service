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
