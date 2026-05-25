// @Author  xuanshuiyuan
// 数据引擎初始化包：统一管理 MySQL、MongoDB、Redis 连接
package engine

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"gorm.io/gorm"
)

// database 数据库连接集合，统一管理 MySQL 和 Redis 连接
type database struct {
	Mysql *gorm.DB    // MySQL 连接（GORM）
	Redis *redis.Pool // Redis 连接池
}

// DB 全局数据库实例，初始化后可全局访问
var DB *database

var log *goxy.Logs

// NewEngine 初始化所有数据引擎（MySQL、Redis）
// 必须在 conf.ConfInit() 之后调用
func NewEngine() error {
	mysqlDB, err := newMysql(conf.Config.Mysql)
	if err != nil {
		return fmt.Errorf("mysql init failed: %w", err)
	}
	DB = &database{
		Mysql: mysqlDB,
		//Mongo: newMongoDXJ(conf.Config.Mongodb),
		Redis: newRedis(conf.Config.Redis),
	}
	return nil
}
