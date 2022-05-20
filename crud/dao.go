package crud

type DAOInterface interface {
	GetModel() ModelInterface
	NewWrapper(modelParams ModelInterface, baseParams *BaseQueryParams) QueryWrapperInterface

	Insert(e interface{}) bool
	Update(e interface{}) bool
	// Delete 逻辑删除
	Delete(pk interface{}, operator int64) bool
	// Remove 物理删除
	Remove(pk interface{}) bool
	FindByPk(pk interface{}) interface{}
	FindOneByColumn(column string, value interface{}) interface{}
	CountByPk(pk interface{}) int64
	CountByColumn(column string, value interface{}) int64
	FindPage(modelParams ModelInterface, baseParams *BaseQueryParams) (interface{}, PageData)
	FindList(modelParams ModelInterface, baseParams *BaseQueryParams) interface{}
	FindAll(modelParams ModelInterface, baseParams *BaseQueryParams) interface{}
	FindListByColumn(column string, value interface{}) interface{}
}

type BaseDao struct {
	Model ModelInterface
}

// GetModel 取模型
func (dao *BaseDao) GetModel() ModelInterface {
	return dao.Model
}

// NewWrapper 取模型
func (dao *BaseDao) NewWrapper(modelParams ModelInterface, baseParams *BaseQueryParams) QueryWrapperInterface {
	return &BaseQueryWrapper{
		ModelParams: modelParams,
		BaseParams:  baseParams,
	}
}

// Insert 插入数据
func (dao *BaseDao) Insert(e interface{}) bool {
	if DbSess().Create(e).Error == nil {
		return true
	}
	return false
}

// Update 更新数据
func (dao *BaseDao) Update(e interface{}) bool {
	if DbSess().Updates(e).Error == nil {
		return true
	}
	return false
}

// Delete 删除数据(逻辑)
func (dao *BaseDao) Delete(pk interface{}, operator int64) bool {
	if dao.UpdateColumn(pk, "deleted", FlagYes, operator) == nil {
		return true
	}
	return false
}

// Remove 删除数据(物理)
func (dao *BaseDao) Remove(pk interface{}) bool {
	if DbSess().Table(dao.Model.Table()).Delete("id = ?", pk).Error == nil {
		return true
	}
	return false
}

// FindByPk 根据主键查询
func (dao *BaseDao) FindByPk(pk interface{}) interface{} {
	if dao.CountByPk(pk) == 0 {
		return nil
	}
	dst := dao.Model.NewModel()
	DbSess().Where("id = ?", pk).First(dst)
	return dst
}

// FindOneByColumn 根据某列查询
func (dao *BaseDao) FindOneByColumn(column string, value interface{}) interface{} {
	if dao.CountByColumn(column, value) == 0 {
		return nil
	}
	dst := dao.Model.NewModel()
	DbSess().Where(column+" = ?", value).First(dst)
	return dst
}

// CountByPk 根据主键查询
func (dao *BaseDao) CountByPk(pk interface{}) int64 {
	var cnt int64
	DbSess().Table(dao.Model.Table()).Where("id = ?", pk).Count(&cnt)
	return cnt
}

// CountByColumn 根据某列查询
func (dao *BaseDao) CountByColumn(column string, value interface{}) int64 {
	var cnt int64
	DbSess().Table(dao.Model.Table()).Where(column+" = ?", value).Count(&cnt)
	return cnt
}

// UpdateColumn 更新字段
func (dao *BaseDao) UpdateColumn(pk interface{}, column string, v interface{}, operator int64) error {
	return DbSess().Table(dao.Model.Table()).Where("id = ?", pk).Updates(map[string]interface{}{
		column:       v,
		"updated_by": operator,
	}).Error
}

// UpdateStatus 更新状态
func (dao *BaseDao) UpdateStatus(pk int64, status interface{}, operator int64) error {
	return dao.UpdateColumn(pk, "status", status, operator)
}

// FindPage 查询
func (dao *BaseDao) FindPage(modelParams ModelInterface, baseParams *BaseQueryParams) (interface{}, PageData) {
	rows := dao.Model.NewModels()
	var total int64
	db := DbSess()
	wrapper := dao.NewWrapper(modelParams, baseParams)
	db.Scopes(wrapper.QueryScope()).Table(dao.Model.Table()).Count(&total)
	db.Scopes(wrapper.QueryScope(), wrapper.PageScope()).Find(&rows)
	return rows, wrapper.PageResult(total)
}

// FindList 查询
func (dao *BaseDao) FindList(modelParams ModelInterface, baseParams *BaseQueryParams) interface{} {
	rows := dao.Model.NewModels()
	wrapper := dao.NewWrapper(modelParams, baseParams)
	DbSess().Scopes(wrapper.QueryScope(), wrapper.PageScope()).Find(&rows)
	return rows
}

// FindAll 查询
func (dao *BaseDao) FindAll(modelParams ModelInterface, baseParams *BaseQueryParams) interface{} {
	rows := dao.Model.NewModels()
	wrapper := dao.NewWrapper(modelParams, baseParams)
	DbSess().Scopes(wrapper.QueryScope()).Find(&rows)
	return rows
}

// FindListByColumn 查询
func (dao *BaseDao) FindListByColumn(column string, value interface{}) interface{} {
	rows := dao.Model.NewModels()
	DbSess().Where(column+" = ?", value).Find(&rows)
	return rows
}
