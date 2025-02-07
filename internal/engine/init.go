// @Author  xuanshuiyuan
package engine

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"gorm.io/gorm"
)

type database struct {
	//Mongo *mongoDXJ
	Mysql *gorm.DB
	Redis *redis.Pool
}

var DB *database

var log *goxy.Logs

//初始化
func NewEngine() {
	DB = &database{
		Mysql: newMysql(conf.Config.Mysql),
		//Mongo: newMongoDXJ(conf.Config.Mongodb),
		Redis: newRedis(conf.Config.Redis),
	}
}
