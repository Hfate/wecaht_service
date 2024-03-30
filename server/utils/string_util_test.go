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
