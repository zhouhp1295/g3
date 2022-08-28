// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package helpers

import "regexp"

// ParseQueryString 解析Query字符串
func ParseQueryString(queryString string) map[string]string {
	r := regexp.MustCompile(`(\w+)\=([\w\-\.\+]+)`)
	matches := r.FindAllStringSubmatch(queryString, -1)
	result := make(map[string]string, len(matches))
	for _, v := range matches {
		result[v[1]] = v[2]
	}
	return result
}
