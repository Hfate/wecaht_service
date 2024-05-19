package utils

import (
	"fmt"
	"testing"
)

func Test_extractArticles(t *testing.T) {
	paths, _ := ReadPDFFiles("/Users/changqi.huang/Documents/hao/xinli/")

	fmt.Println(len(paths))

}
