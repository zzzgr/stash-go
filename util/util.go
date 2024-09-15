package util

import "strings"

// IsMatch 是否匹配, s待匹配字符串, p规则字符串  *代表0个或多个字符 ?代表一个字符
func IsMatch(s string, p string) bool {
	// s 的索引位置
	i := 0
	// p 的索引位置
	j := 0
	// 通配符时回溯的位置
	ii := -1
	jj := -1
	for i < len(s) {
		if j < len(p) && p[j] == '*' {
			// 遇到通配符了,记录下位置,规则字符串+1,定位到非通配字符串
			ii = i
			jj = j
			j++
		} else if j < len(p) && (s[i] == p[j] || p[j] == '?') {
			// 匹配到了
			i++
			j++
		} else {
			// 匹配失败,需要判断 s 是否被 p 的 * 号匹配着
			if jj == -1 { // 前面没有通配符
				return false
			}
			// 回到之前记录通配符的位置
			j = jj
			// 带匹配字符串也回到记录的位置,并后移一位
			i = ii + 1
		}
	}
	// 当 s 的每一个字段都匹配成功以后,判断 p 剩下的串,是*则放行
	for j < len(p) && p[j] == '*' {
		j++
	}
	// 检测到最后就匹配成功
	return j == len(p)
}

// GetCookieFieldMap cookie转成kv的map
func GetCookieFieldMap(cookie string) map[string]string {
	pairs := strings.Split(cookie, ";")
	ckMap := make(map[string]string, 0)
	for _, pair := range pairs {

		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) == 2 {
			ckMap[kv[0]] = kv[1]
		}
	}
	return ckMap
}
