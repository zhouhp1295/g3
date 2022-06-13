package crud

import (
	"fmt"
	"github.com/zhouhp1295/g3/helpers"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

// DefaultPageSize 默认分页大小
const DefaultPageSize = 20

type BaseQueryParams struct {
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
	OrderBy   string `form:"orderBy" json:"orderBy"`
}

type QueryWrapperInterface interface {
	QueryScope() func(db *gorm.DB) *gorm.DB
	WrapQuery(db *gorm.DB)
	PageScope() func(db *gorm.DB) *gorm.DB
	PageResult(total int64) PageData
}

type BaseQueryWrapper struct {
	ModelParams ModelInterface
	BaseParams  *BaseQueryParams
}

func (wrapper *BaseQueryWrapper) QueryScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if wrapper.BaseParams != nil {
			if len(wrapper.BaseParams.BeginTime) > 0 {
				db.Where("created_at > ?", wrapper.BaseParams.BeginTime+" 00:00:00")
			}
			if len(wrapper.BaseParams.EndTime) > 0 {
				db.Where("created_at < ?", wrapper.BaseParams.EndTime+" 23:59:59")
			}
		}
		wrapper.WrapQuery(db)

		db.Where("deleted = ?", FlagNo)

		return db
	}
}

//wrapperQuery 查询组装
func wrapperQuery(table string, params interface{}, db *gorm.DB) {
	var t reflect.Type
	var v reflect.Value
	if reflect.TypeOf(params).Kind() == reflect.Ptr {
		t = reflect.TypeOf(params).Elem()
		v = reflect.ValueOf(params).Elem()
	} else if reflect.TypeOf(params).Kind() == reflect.Struct {
		t = reflect.TypeOf(params)
		v = reflect.ValueOf(params)
	} else {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}
		if v.Field(i).Interface() == nil {
			continue
		}
		if reflect.TypeOf(v.Field(i).Interface()).Kind() == reflect.Struct {
			wrapperQuery(table, v.Field(i).Interface(), db)
			continue
		}
		sf := t.Field(i)
		if query := sf.Tag.Get("query"); len(query) > 0 {
			queryItems := strings.Split(strings.ToLower(query), ";")
			if helpers.IndexOf[string](queryItems, "skipnil") >= 0 {
				if str, ok := v.Field(i).Interface().(string); ok && len(str) == 0 {
					continue
				}
				if i, ok := helpers.Int64(v.Field(i).Interface()); ok && i == 0 {
					continue
				}
			}
			//目前只简单处理列的名字
			colName := db.Statement.NamingStrategy.ColumnName(table, sf.Name)
			if helpers.IndexOf[string](queryItems, "like") >= 0 {
				db.Where(colName+" like ?", fmt.Sprintf("%%%v%%", v.Field(i).Interface()))
			} else if helpers.IndexOf[string](queryItems, "eq") >= 0 {
				db.Where(colName+" = ?", v.Field(i).Interface())
			}
		}
	}
}

// WrapQuery 查询组装
// 标签: query
func (wrapper *BaseQueryWrapper) WrapQuery(db *gorm.DB) {
	wrapperQuery(wrapper.ModelParams.Table(), wrapper.ModelParams, db)
}

func (wrapper *BaseQueryWrapper) PageScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if wrapper.BaseParams.PageNum <= 0 {
			wrapper.BaseParams.PageNum = 1
		}
		if wrapper.BaseParams.PageSize <= 0 {
			wrapper.BaseParams.PageSize = DefaultPageSize
		}
		return db.
			Offset((wrapper.BaseParams.PageNum - 1) * wrapper.BaseParams.PageSize).
			Limit(wrapper.BaseParams.PageSize).
			Order(wrapper.BaseParams.OrderBy)
	}
}

func (wrapper *BaseQueryWrapper) PageResult(total int64) PageData {
	return PageResult(wrapper.BaseParams.PageNum, wrapper.BaseParams.PageSize, int(total))
}
