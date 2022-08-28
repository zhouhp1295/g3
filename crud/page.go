// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package crud

import "math"

type PageData struct {
	PageNum   int `json:"pageNum"`
	PageTotal int `json:"pageTotal"`
	PageSize  int `json:"pageSize"`
	Total     int `json:"total"`
}

func PageResult(page, size, total int) PageData {
	return PageData{
		PageNum:   page,
		PageTotal: int(math.Ceil(float64(total) / float64(size))),
		PageSize:  size,
		Total:     total,
	}
}
