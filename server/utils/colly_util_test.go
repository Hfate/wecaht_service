package utils

import "testing"

func Test_collectImgUrl(t *testing.T) {
	CollectArticle("aa")
}

func TestCollectArticle(t *testing.T) {
	CollectArticle("以色列：伊朗将承受冲突升级的后果")
}
