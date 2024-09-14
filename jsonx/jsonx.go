package jsonx

import (
	"encoding/json"
)

// ToStr 转换为json字符串
func ToStr(val any) string {
	if strBytes, err := json.Marshal(val); err != nil {
		return ""
	} else {
		return string(strBytes)
	}
}

// ToMap json字符串转map
func ToMap(jsonStr string) map[string]any {
	return ToMapByBytes([]byte(jsonStr))
}

// To json字符串转结构体
func To[T any](jsonStr string, res T) (T, error) {
	return res, json.Unmarshal([]byte(jsonStr), &res)
}

// ToMapByBytes json字节数组转map
func ToMapByBytes(bytes []byte) map[string]any {
	var res map[string]any
	_ = json.Unmarshal(bytes, &res)
	return res
}
