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
func newMysql(mysqlBase *conf.Mysql) *gorm.DB {
	var DbLogMode logger.Interface
	if mysqlBase.DbLogMode == true {
		DbLogMode = logger.Default.LogMode(logger.Info)
	} else {
		DbLogMode = logger.Default.LogMode(logger.Silent)
	}
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=%s&parseTime=%d&loc=%s", mysqlBase.Username, mysqlBase.Password, mysqlBase.Network, mysqlBase.Hostname, mysqlBase.Port, mysqlBase.DataBase, mysqlBase.Charset, 1, mysqlBase.TimeZone)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       conn,  // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: DbLogMode,
	})
	if err != nil {
		panic(err)
	}
	//db.Debug(mysqlBase.DbLogMode)
	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
