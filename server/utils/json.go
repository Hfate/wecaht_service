package utils

import "encoding/json"

func JsonStrToStruct(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}
