package crud

type DAOInterface interface {
	GetModel() ModelInterface
	NewWrapper(modelParams ModelInterface, baseParams *BaseQueryParams) QueryWrapperInterface
	// Insert 插入
	Insert(m ModelInterface, operator int64) bool
	// Update 更新
	Update(m ModelInterface, operator int64) bool
	// UpdateColumn 更新
	UpdateColumn(pk interface{}, column string, v interface{}, operator int64) bool
	// UpdateStatus 更新
	UpdateStatus(pk int64, status interface{}, operator int64) bool
	// Delete 删除
	Delete(m ModelInterface, operator int64) bool
	// DeleteByPk 逻辑删除
	DeleteByPk(pk interface{}, operator int64) bool
	// Remove 删除
	Remove(m ModelInterface, operator int64) bool
	// RemoveByPk 物理删除
	RemoveByPk(pk interface{}) bool
	FindByPk(pk interface{}) ModelInterface
	FindOneByColumn(column string, value interface{}) ModelInterface
	CountByPk(pk interface{}) int64
	CountByColumn(column string, value interface{}) int64
	FindPage(modelParams ModelInterface, baseParams *BaseQueryParams) (interface{}, PageData)
	FindList(modelParams ModelInterface, baseParams *BaseQueryParams) interface{}
	FindAll(modelParams ModelInterface, baseParams *BaseQueryParams) interface{}
	FindListByColumn(column string, value interface{}) interface{}
	// BeforeInsert 插入之前
	BeforeInsert(m ModelInterface) (ok bool, msg string)
	// AfterInsert 插入之后
	AfterInsert(m ModelInterface) (ok bool, msg string)
	// BeforeUpdate 更新之前
	BeforeUpdate(m ModelInterface) (ok bool, msg string)
	// AfterUpdate 更新之后
	AfterUpdate(m ModelInterface) (ok bool, msg string)
	// BeforeRemove 移除之前
	BeforeRemove(m ModelInterface) (ok bool, msg string)
	// AfterRemove 移除之后
	AfterRemove(m ModelInterface) (ok bool, msg string)
	// BeforeDelete 删除之前
	BeforeDelete(m ModelInterface) (ok bool, msg string)
	// AfterDelete 删除之后
	AfterDelete(m ModelInterface) (ok bool, msg string)
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
func (dao *BaseDao) Insert(m ModelInterface, operator int64) bool {
	m.SetCreatedBy(operator)
	m.SetUpdatedBy(operator)
	if DbSess().Create(m).Error == nil {
		return true
	}
	return false
}

// Update 更新数据
func (dao *BaseDao) Update(m ModelInterface, operator int64) bool {
	m.SetUpdatedBy(operator)

	sess := DbSess()

	updateCols := m.GetUpdateColumns()
	if updateCols != nil && len(updateCols) > 0 {
		sess = sess.Select(updateCols)
	}

	omitCols := m.GetOmitColumns()
	if omitCols != nil && len(omitCols) > 0 {
		sess = sess.Omit(omitCols...)
	}

	if sess.Updates(m).Error == nil {
		return true
	}
	return false
}

func (dao *BaseDao) Delete(m ModelInterface, operator int64) bool {
	return dao.DeleteByPk(m.GetId(), operator)
}

// DeleteByPk 删除数据(逻辑)
func (dao *BaseDao) DeleteByPk(pk interface{}, operator int64) bool {
	return dao.UpdateColumn(pk, "deleted", FlagYes, operator)
}

func (dao *BaseDao) Remove(m ModelInterface, operator int64) bool {
	return dao.RemoveByPk(m.GetId())
}

// RemoveByPk 删除数据(物理)
func (dao *BaseDao) RemoveByPk(pk interface{}) bool {
	if DbSess().Table(dao.Model.Table()).Delete("id = ?", pk).Error == nil {
		return true
	}
	return false
}

// FindByPk 根据主键查询
func (dao *BaseDao) FindByPk(pk interface{}) ModelInterface {
	if dao.CountByPk(pk) == 0 {
		return nil
	}
	dst := dao.Model.NewModel()
	DbSess().Where("id = ?", pk).First(dst)
	return dst
}

// FindOneByColumn 根据某列查询
func (dao *BaseDao) FindOneByColumn(column string, value interface{}) ModelInterface {
	if dao.CountByColumn(column, value) == 0 {
		return nil
	}
	dst := dao.Model.NewModel()
	DbSess().Where(column+" = ? and deleted = ?", value, FlagNo).First(dst)
	return dst
}

// CountByPk 根据主键查询
func (dao *BaseDao) CountByPk(pk interface{}) int64 {
	var cnt int64
	DbSess().Table(dao.Model.Table()).Where("id = ? and deleted = ?", pk, FlagNo).Count(&cnt)
	return cnt
}

// CountByColumn 根据某列查询
func (dao *BaseDao) CountByColumn(column string, value interface{}) int64 {
	var cnt int64
	DbSess().Table(dao.Model.Table()).Where(column+" = ? and deleted = ?", value, FlagNo).Count(&cnt)
	return cnt
}

// UpdateColumn 更新字段
func (dao *BaseDao) UpdateColumn(pk interface{}, column string, v interface{}, operator int64) bool {
	if DbSess().Table(dao.Model.Table()).Where("id = ?", pk).Updates(map[string]interface{}{
		column:       v,
		"updated_by": operator,
	}).Error != nil {
		return false
	}
	return true
}

// UpdateStatus 更新状态
func (dao *BaseDao) UpdateStatus(pk int64, status interface{}, operator int64) bool {
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

func (dao *BaseDao) BeforeInsert(m ModelInterface) (ok bool, msg string) {
	ok = true
	return
}
func (dao *BaseDao) AfterInsert(m ModelInterface) (ok bool, msg string) {
	ok = true
	return
}
func (dao *BaseDao) BeforeUpdate(m ModelInterface) (ok bool, msg string) {
	ok = true
	return
}
func (dao *BaseDao) AfterUpdate(m ModelInterface) (ok bool, msg string) {
	ok = true
	return
}
func (dao *BaseDao) BeforeRemove(m ModelInterface) (ok bool, msg string) {
	ok = true
	return
}
func (dao *BaseDao) AfterRemove(m ModelInterface) (ok bool, msg string) {
	ok = true
	return
}
func (dao *BaseDao) BeforeDelete(m ModelInterface) (ok bool, msg string) {
	ok = true
	return
}
func (dao *BaseDao) AfterDelete(m ModelInterface) (ok bool, msg string) {
	ok = true
	return
}
