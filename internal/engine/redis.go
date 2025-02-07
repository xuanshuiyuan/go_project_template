// @Author  xuanshuiyuan
package engine

import (
	"github.com/gomodule/redigo/redis"
	"go_project_template/internal/conf"
	"time"
)

//初始化
func newRedis(redisBase *conf.Redis) *redis.Pool {
	//redisBase := conf.Config.Redis
	return &redis.Pool{
		MaxIdle:     redisBase.Idle,
		MaxActive:   redisBase.Active,
		IdleTimeout: time.Duration(redisBase.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(redisBase.Network, redisBase.Addr,
				redis.DialConnectTimeout(time.Duration(redisBase.DialTimeout)),
				redis.DialReadTimeout(time.Duration(redisBase.ReadTimeout)),
				redis.DialWriteTimeout(time.Duration(redisBase.WriteTimeout)),
				redis.DialPassword(redisBase.Auth),
				redis.DialDatabase(redisBase.DataBase),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}

// Close close the resource.
func RedisClose() error {
	return DB.Redis.Close()
}

// Ping dao ping.
func RedisPing() (err error) {
	conn := DB.Redis.Get()
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
