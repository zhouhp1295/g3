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
	GetId() int64
	// Table 返回表名
	Table() string
	// NewModel 返回实例
	NewModel() ModelInterface
	// NewModels 返回实例数组
	NewModels() interface{}
	// SetCreatedBy 设置操作人
	SetCreatedBy(operator int64)
	// SetUpdatedBy 设置操作人
	SetUpdatedBy(operator int64)
	// GetUpdateColumns 更新时的列
	GetUpdateColumns() []string
	// GetOmitColumns 更新忽略的列
	GetOmitColumns() []string
	// GetSelectColumns 查询时的列
	GetSelectColumns() []string
	// SetLastModel 设置last
	SetLastModel(last ModelInterface)
}

// TailColumns 通用的列,一般放在末尾
type TailColumns struct {
	Remark    string    `gorm:"TYPE:VARCHAR(100);COMMENT:备注" json:"remark" form:"remark"`
	Status    string    `gorm:"TYPE:CHAR(1);NOT NULL;DEFAULT:1;COMMENT:状态 0=NO,不可用 1=YES,正常"  json:"status" form:"status" query:"eq"`
	Deleted   string    `gorm:"TYPE:CHAR(1);NOT NULL;DEFAULT:0;COMMENT:删除标识 0=NO,正常 1=YES,删除" json:"-" query:"eq"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt"`
	CreatedBy int64     `gorm:"NOT NULL;DEFAULT:0" json:"-"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt"`
	UpdatedBy int64     `gorm:"NOT NULL;DEFAULT:0" json:"-"`
}

func (tailColumns *TailColumns) SetCreatedBy(operator int64) {
	tailColumns.CreatedBy = operator
}

func (tailColumns *TailColumns) SetUpdatedBy(operator int64) {
	tailColumns.UpdatedBy = operator
}

type BaseModel struct {
	Id   int64          `json:"id" form:"id"`
	Last ModelInterface `gorm:"-" json:"-"`
}

func (baseModel *BaseModel) GetId() int64 {
	return baseModel.Id
}

func (baseModel *BaseModel) GetUpdateColumns() []string {
	return nil
}

func (baseModel *BaseModel) GetOmitColumns() []string {
	return nil
}

func (baseModel *BaseModel) GetSelectColumns() []string {
	return nil
}

func (baseModel *BaseModel) SetLastModel(last ModelInterface) {
	baseModel.Last = last
}
