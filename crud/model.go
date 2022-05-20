package crud

import "time"

const (
	FlagTrue  = "T"
	FlagFalse = "F"
)

const (
	FlagNo  = "0"
	FlagYes = "1"
)

type ModelInterface interface {
	// Table 返回表名
	Table() string
	// NewModel 返回实例
	NewModel() interface{}
	// NewModels 返回实例数组
	NewModels() interface{}
}

// TailColumns 通用的列,一般放在末尾
type TailColumns struct {
	Remark    string    `gorm:"TYPE:VARCHAR(100);COMMENT:'备注'" json:"remark" form:"remark"`
	Status    string    `gorm:"TYPE:CHAR(1);NOT NULL;DEFAULT:'1';COMMENT:'状态:0=NO,不可用 1=YES,正常'"  json:"status" form:"status" query:"eq"`
	Deleted   string    `gorm:"TYPE:CHAR(1);NOT NULL;DEFAULT:'0';COMMENT:'删除标识:0=NO,正常 1=YES,删除'" json:"-" query:"eq"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt"`
	CreatedBy int64     `gorm:"NOT NULL;DEFAULT:0" json:"-"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt"`
	UpdatedBy int64     `gorm:"NOT NULL;DEFAULT:0" json:"-"`
}

type BaseModel struct {
	Id int64 `json:"id" form:"id"`
}
