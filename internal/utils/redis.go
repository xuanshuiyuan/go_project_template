// @Author  xuanshuiyuan
// Redis 通用操作封装，基于连接池提供常用 Redis 命令的链式调用
package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_project_template/internal/conf"
	"go_project_template/internal/engine"
	"time"

	"github.com/gomodule/redigo/redis"
)

// RedisService Redis 通用操作服务
// 使用方式: NewRedis().SetKey("key").SetValue("value").SetExp("60").RedisSetAndEx()
type RedisService struct {
	Key   string      // Redis 键名
	Value interface{} // Redis 值
	Exp   string      // 过期时间（秒，字符串格式）
}

// NewRedis 创建 RedisService 实例
func NewRedis() *RedisService {
	return &RedisService{}
}

// SetKey 设置键名（链式调用）
func (r *RedisService) SetKey(key string) *RedisService {
	r.Key = key
	return r
}

// SetValue 设置值（链式调用）
func (r *RedisService) SetValue(value interface{}) *RedisService {
	r.Value = value
	return r
}

// SetExp 设置过期时间，单位秒（链式调用）
func (r *RedisService) SetExp(exp string) *RedisService {
	r.Exp = exp
	return r
}

// execCommand 从连接池获取连接并执行 Redis 命令
func (r *RedisService) execCommand(command string, args ...interface{}) (interface{}, error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	return conn.Do(command, args...)
}

// ==================== 自增/自减 ====================

// Incr 对 key 执行自增操作
func (r *RedisService) Incr(key string) (res int64, err error) {
	result, err := r.execCommand("INCR", key)
	if err != nil {
		return
	}
	if v, ok := result.(int64); ok {
		res = v
	}
	return
}

// Decr 对 key 执行自减操作
func (r *RedisService) Decr(key string) (res int64, err error) {
	result, err := r.execCommand("DECR", key)
	if err != nil {
		return
	}
	if v, ok := result.(int64); ok {
		res = v
	}
	return
}

// ==================== 有序集合 ====================

// Zscore 返回有序集中成员的分数值，成员不存在或 key 不存在返回空字符串
func (r *RedisService) Zscore(key, member string) (res string, err error) {
	result, err := r.execCommand("ZSCORE", key, member)
	if err != nil {
		return
	}
	if v, ok := result.([]byte); ok {
		res = string(v)
	}
	return
}

// ZscoreDelayQueue 查询延迟队列中有序集成员的分数值
// key 会自动拼接 {:key}:waiting 格式，与延迟队列内部结构对应
func (r *RedisService) ZscoreDelayQueue(key, member string) (res string, err error) {
	result, err := r.execCommand("ZSCORE", fmt.Sprintf("{%s}:waiting", key), member)
	if err != nil {
		return
	}
	if v, ok := result.([]byte); ok {
		res = string(v)
	}
	return
}

// ZAdd 向有序集合添加成员
// messages 格式: score1, member1, score2, member2, ...
func (r *RedisService) ZAdd(key string, messages ...interface{}) error {
	args := []interface{}{key}
	args = append(args, messages...)
	_, err := r.execCommand("ZADD", args...)
	return err
}

// ==================== 集合操作 ====================

// SAdd 向集合添加元素
func (r *RedisService) SAdd(key string, messages ...interface{}) error {
	args := []interface{}{key}
	args = append(args, messages...)
	_, err := r.execCommand("SADD", args...)
	return err
}

// SCard 获取集合元素个数
func (r *RedisService) SCard(key string) (size interface{}, err error) {
	return r.execCommand("SCARD", key)
}

// SIsMember 判断元素是否在集合中，返回 1 表示存在，0 表示不存在
func (r *RedisService) SIsMember(key string, message interface{}) (res int64, err error) {
	result, err := r.execCommand("SISMEMBER", key, message)
	if err != nil {
		return
	}
	v, ok := result.(int64)
	if !ok {
		return 0, fmt.Errorf("SIsMember 返回类型异常: %T", result)
	}
	res = v
	return
}

// SMembers 获取集合中所有元素
func (r *RedisService) SMembers(key string) (result []string, err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	reply, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		return
	}
	result = reply
	return
}

// SRem 删除集合元素，返回成功删除的个数
func (r *RedisService) SRem(key string, messages ...interface{}) (res int64, err error) {
	args := []interface{}{key}
	args = append(args, messages...)
	result, err := r.execCommand("SREM", args...)
	if err != nil {
		return
	}
	v, ok := result.(int64)
	if !ok {
		return 0, fmt.Errorf("SRem 返回类型异常: %T", result)
	}
	res = v
	return
}

// SPop 随机返回并删除集合中的一个元素
func (r *RedisService) SPop(key string) (res string, err error) {
	result, err := r.execCommand("SPOP", key)
	if err != nil {
		return
	}
	if v, ok := result.([]byte); ok {
		res = string(v)
	}
	return
}

// SPopN 随机返回并删除集合中的 N 个元素
func (r *RedisService) SPopN(key string, size int64) (result []string, err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	reply, err := redis.Strings(conn.Do("SPOP", key, size))
	if err != nil {
		return
	}
	result = reply
	return
}

// ==================== Lua 脚本 ====================

// EvalSha 执行已加载的 Lua 脚本
func (r *RedisService) EvalSha(sha1 string, values []interface{}) (interface{}, error) {
	args := []interface{}{sha1}
	args = append(args, values...)
	return r.execCommand("EVALSHA", args...)
}

// LoadScript 加载 Lua 脚本到 Redis，返回 sha1 校验值
func (r *RedisService) LoadScript(script string) error {
	_, err := r.execCommand("SCRIPT", "LOAD", script)
	return err
}

// ==================== String 操作 ====================

// RedisSetAndEx 设置键值并指定过期时间
// 使用方式: NewRedis().SetKey("key").SetValue("value").SetExp("60").RedisSetAndEx()
func (r *RedisService) RedisSetAndEx() error {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	_, err := conn.Do("SET", r.Key, r.Value, "EX", r.Exp)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("set ", r.Key, "failed, err:", err)
		return err
	}
	return nil
}

// GetStringKey 根据 key 获取字符串值
func (r *RedisService) GetStringKey(key string) (string, error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	result, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("get ", key, " failed, err:", err)
		return result, err
	}
	return result, nil
}

// GetInfoByKey 根据 key 获取 JSON 对象并反序列化为 map
func (r *RedisService) GetInfoByKey(key string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	info, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("get ", key, " failed, err:", err)
		return result, err
	}
	if er := json.Unmarshal([]byte(info), &result); er != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("json unmarshal ", key, " failed, err:", er)
		return result, er
	}
	return result, nil
}

// GetArrByKey 根据 key 获取 JSON 数组并反序列化为切片
func (r *RedisService) GetArrByKey(key string) (result []interface{}, err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	info, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("get ", key, " failed, err:", err)
		return result, err
	}
	if er := json.Unmarshal([]byte(info), &result); er != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("json unmarshal ", key, " failed, err:", er)
		return result, er
	}
	return result, nil
}

// Get 根据 r.Key 获取字符串值，key 不存在时不报错
func (r *RedisService) Get() (res string, err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	res, err = redis.String(conn.Do("GET", r.Key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		}
	}
	return
}

// Set 设置键值对（无过期时间）
func (r *RedisService) Set() (err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	_, err = conn.Do("SET", r.Key, r.Value)
	return
}

// RedisGetSign 获取接口签名，用于验证签名是否已存在（防重复请求）
func (r *RedisService) RedisGetSign() error {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	_, err := redis.String(conn.Do("GET", r.Key))
	return err
}

// RedisVerification 验证 token 有效性，验证通过后自动续期
func (r *RedisService) RedisVerification(key string, token string) error {
	data, err := r.GetInfoByKey(key)
	if err != nil {
		return err
	}
	if data["token"] != token {
		return errors.New("Im Invalid Request!")
	}
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	n, err := conn.Do("EXPIRE", key, conf.RedisTokenExp)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("expiring ", key, " failed, err:", err)
		return err
	}
	if n != int64(1) {
		log.Error(conf.Config.Base.LogFileName, "").Println("expiring ", key, " failed, key may not exist")
		return errors.New("token 续期失败")
	}
	return nil
}

// ==================== 删除 ====================

// Del 删除指定 key
func (r *RedisService) Del() (err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	_, err = conn.Do("DEL", r.Key)
	return
}

// ==================== 分布式锁 ====================

// Lock 获取分布式锁（基于 SET NX EX 实现）
// 使用方式: isLock, err := NewRedis().SetKey("lock_key").SetExp("30").Lock()
// 返回 isLock=true 表示获取锁成功，需要在业务处理完后调用 UnLock 释放锁
func (r *RedisService) Lock() (isLock bool, err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	_, err = redis.String(conn.Do("SET", r.Key, 1, "EX", r.Exp, "NX"))
	if err != nil {
		if err == redis.ErrNil {
			// 锁已被其他请求持有，非错误
			return false, nil
		}
		log.Error(conf.Config.Base.LogFileName, "").Println("lock ", r.Key, " failed, err:", err)
		return false, err
	}
	return true, nil
}

// UnLock 释放分布式锁（删除 key）
// 使用方式: err := NewRedis().SetKey("lock_key").UnLock()
func (r *RedisService) UnLock() (err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	_, err = conn.Do("DEL", r.Key)
	return
}

// Spin 自旋获取分布式锁，在指定次数内反复尝试获取锁
// 参数:
//   - r: RedisService 实例
//   - key: 锁的键名
//   - exp: 锁的过期时间（秒），传空则默认 30 秒
//   - maxRetry: 最大重试次数，传 0 则默认 60 次
//
// 使用方式:
//   isLock := Spin(NewRedis(), "lock_key", "30", 60)  // 自定义参数
//   isLock := Spin(NewRedis(), "lock_key", "", 0)     // 使用默认值：30秒超时，60次重试
//
// 返回 isLock=true 表示在重试范围内成功获取锁
// 每 500ms 重试一次，maxRetry=60 时最长等待约 30 秒
func Spin(r *RedisService, key string) (isLock bool) {
	const (
		exp      = "30"
		maxRetry = 60
	)
	for i := int64(0); i < maxRetry; i++ {
		isLock, err := r.SetKey(key).SetExp(exp).Lock()
		if err != nil {
			log.Error(conf.Config.Base.LogFileName, "redis.log").Println("spin error:", err)
			return false
		}
		if isLock {
			return true
		}
		time.Sleep(time.Millisecond * 500)
	}
	return false
}
