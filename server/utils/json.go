package utils

import "encoding/json"

func JsonStrToStruct(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}

func Parse2Json(v interface{}) string {
	// 将结构体转换为JSON格式
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}
