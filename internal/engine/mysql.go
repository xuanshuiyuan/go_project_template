// @Author  xuanshuiyuan
package engine

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"go_project_template/internal/conf"
	"time"
)

//初始化
func newMysql(mysqlBase *conf.Mysql) (*gorm.DB, error) {
	var DbLogMode logger.Interface
	if mysqlBase.DbLogMode {
		DbLogMode = logger.Default.LogMode(logger.Info)
	} else {
		DbLogMode = logger.Default.LogMode(logger.Silent)
	}
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=%s&parseTime=%d&loc=%s", mysqlBase.Username, mysqlBase.Password, mysqlBase.Network, mysqlBase.Hostname, mysqlBase.Port, mysqlBase.DataBase, mysqlBase.Charset, 1, mysqlBase.TimeZone)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       conn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: DbLogMode,
	})
	if err != nil {
		return nil, fmt.Errorf("mysql connect failed: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("mysql get sql.DB failed: %w", err)
	}
	maxIdle := 10
	maxOpen := 100
	if mysqlBase.MaxIdleConns > 0 {
		maxIdle = mysqlBase.MaxIdleConns
	}
	if mysqlBase.MaxOpenConns > 0 {
		maxOpen = mysqlBase.MaxOpenConns
	}
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}
