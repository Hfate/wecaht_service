package utils

import "unicode"

func RemoveChinese(s string) string {
	result := ""
	for _, r := range s {
		if !unicode.Is(unicode.Scripts["Han"], r) { // 判断字符是否为中文字符
			result += string(r)
		}
	}
	return result
}
