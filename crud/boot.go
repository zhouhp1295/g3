package crud

import "gorm.io/gorm"

var dbEngine *gorm.DB

func InitDbEngine(engine *gorm.DB) {
	dbEngine = engine

	//初始化数据库迁移记录表
	err := MigrateTables(engine, []interface{}{new(Migration)})
	if err != nil {
		panic("Database init error" + err.Error())
	}
}

func DbSess() *gorm.DB {
	return dbEngine.Session(&gorm.Session{})
}
