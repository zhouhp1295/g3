// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package helpers

func IndexOf[T comparable](source []T, search T) int {
	for i, v := range source {
		if v == search {
			return i
		}
	}
	return -1
}

func Reverse[T any](source []T) {
	for i, j := 0, len(source)-1; i < j; i, j = i+1, j-1 {
		source[i], source[j] = source[j], source[i]
	}
}
