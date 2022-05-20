package crud

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type Migration struct {
	BaseModel
	Code string `gorm:"TYPE:VARCHAR(100);UNIQUE;COMMENT:'编号'"`
	TailColumns
}

func SyncTables(tx *gorm.DB, tables []interface{}) error {
	for _, table := range tables {
		strings.TrimPrefix(fmt.Sprintf("%T", table), "*db.")
		err := tx.Migrator().AutoMigrate(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func MigrateTables(tx *gorm.DB, tables []interface{}) error {
	for _, table := range tables {
		if tx.Migrator().HasTable(table) {
			continue
		}
		strings.TrimPrefix(fmt.Sprintf("%T", table), "*db.")
		err := tx.Migrator().AutoMigrate(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func DoMigrate(code string, f func() error) {
	var cnt int64
	err := DbSess().Model(new(Migration)).Where("code = ?", code).Count(&cnt).Error
	if err != nil {
		return
	}
	if cnt > 0 {
		return
	}

	err = f()
	if err != nil {
		return
	}

	e := new(Migration)
	e.Code = code
	DbSess().Create(e)
}
