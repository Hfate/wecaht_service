package utils

import "testing"

func TestCollectArticle(t *testing.T) {
	CollectArticle("以色列：伊朗将承受冲突升级的后果")
}

func TestCollectWechatArticle(t *testing.T) {
	CollectWechatArticle("https://mp.weixin.qq.com/s?__biz=MzUyNzE4OTE1Mw==&mid=2247756365&idx=1&sn=abd56014b9839b7d434f4af12dd4df9b&chksm=fb0e06df07aab431bc5afd2f17a1d027efd09fa770cab2f515bac138d00a59181398a4909bbd&scene=0&xtrack=1#rd")
}
