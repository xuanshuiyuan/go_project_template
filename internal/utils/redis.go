// @Author  xuanshuiyuan
package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go_project_template/internal/conf"
	"go_project_template/internal/engine"
	"time"
)

//RedisService 定义了redis的常用参数
type RedisService struct {
	Key   string
	Value interface{}
	Exp   string
}

func NewRedis() *RedisService {
	return &RedisService{}
}

func (r *RedisService) SetKey(key string) *RedisService {
	r.Key = key
	return r
}

func (r *RedisService) SetValue(value interface{}) *RedisService {
	r.Value = value
	return r
}

func (r *RedisService) SetExp(exp string) *RedisService {
	r.Exp = exp
	return r
}

func (r *RedisService) execCommand(command string, args ...interface{}) (interface{}, error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close() //需要close
	reply, err := conn.Do(command, args...)
	return reply, err
}

func (r *RedisService) Incr(key string) (res int64, err error) {
	result, err := r.execCommand("Incr", key)
	if err != nil {
		return
	}
	var v, ok = result.(int64)
	if ok {
		res = int64(v)
	}
	return
}

func (r *RedisService) Decr(key string) (res int64, err error) {
	result, err := r.execCommand("Decr", key)
	if err != nil {
		return
	}
	var v, ok = result.(int64)
	if ok {
		res = int64(v)
	}
	return
}

func (r *RedisService) EvalSha(sha1 string, values []interface{}) (interface{}, error) {
	args := []interface{}{
		sha1,
	}
	args = append(args, values...)
	res, err := r.execCommand("EVALSHA", args...)
	return res, err
}

func (r *RedisService) LoadScript(script string) error {
	args := []interface{}{
		"LOAD",
		script,
	}
	_, err := r.execCommand("SCRIPT", args...)
	return err
}

//Zscore 命令返回有序集中，成员的分数值。 如果成员元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (r *RedisService) Zscore(key, args string) (res string, err error) {
	result, err := r.execCommand("Zscore", key, args)
	if err != nil {
		return
	}
	var v, ok = result.([]byte)
	if ok {
		res = string(v)
	}
	return
}

//Zscore 命令返回有序集中，成员的分数值。 如果成员元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (r *RedisService) ZscoreDelayQueue(key, args string) (res string, err error) {
	result, err := r.execCommand("Zscore", fmt.Sprintf("{%s}:waiting", key), args)
	if err != nil {
		return
	}
	var v, ok = result.([]byte)
	if ok {
		res = string(v)
	}
	return
}

func (r *RedisService) ZAdd(key string, messages ...interface{}) error {
	args := []interface{}{
		key,
		//"NX",
	}
	for _, message := range messages {
		args = append(args, message)
	}
	_, err := r.execCommand("ZAdd", args...)
	return err
}

//添加集合元素
func (r *RedisService) SAdd(key string, messages ...interface{}) error {
	args := []interface{}{
		key,
		//"NX",
	}
	for _, message := range messages {
		args = append(args, message)
	}
	_, err := r.execCommand("SADD", args...)
	return err
}

//获取集合元素个数
func (r *RedisService) SCard(key string) (size interface{}, err error) {
	size, err = r.execCommand("SCard", key)
	if err != nil {
		return
	}
	return
}

//判断元素是否在集合中 1:在 0:不存在
func (r *RedisService) SIsMember(key string, message interface{}) (res int64, err error) {
	result, err := r.execCommand("SIsMember", key, message)
	if err != nil {
		return
	}
	res = result.(int64)
	return
}

//获取集合中所有的元素
func (r *RedisService) SMembers(key string) (result []string, err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close() //需要close
	err = conn.Send("SMembers", key)
	conn.Flush()
	reply, err := redis.MultiBulk(conn.Receive())
	if err != nil {
		return
	}
	for _, x := range reply {
		var v, ok = x.([]byte)
		if ok {
			result = append(result, string(v))
		}
	}
	return
}

//删除集合元素 1:成功 0:失败
func (r *RedisService) SRem(key string, messages ...interface{}) (res int64, err error) {
	args := []interface{}{
		key,
	}
	for _, message := range messages {
		args = append(args, message)
	}
	result, err := r.execCommand("SRem", args...)
	if err != nil {
		return
	}
	res = result.(int64)
	return
}

//随机返回集合中的元素，并且删除返回的元素
func (r *RedisService) SPop(key string) (res string, err error) {
	result, err := r.execCommand("SPop", key)
	if err != nil {
		return
	}
	var v, ok = result.([]byte)
	if ok {
		res = string(v)
	}
	return
}

//随机返回集合中的元素，并且删除返回的元素
func (r *RedisService) SPopN(key string, size int64) (result []string, err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close() //需要close
	err = conn.Send("SPop", key, size)
	conn.Flush()
	reply, err := redis.MultiBulk(conn.Receive())
	if err != nil {
		return
	}
	for _, x := range reply {
		var v, ok = x.([]byte)
		if ok {
			result = append(result, string(v))
		}
	}
	return
}

// @Title RedisSetAndEx
// @Description redis增加数据和过期时间
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param
// @Return error
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

// @Title GetStringKey
// @Description 根据key获得数据
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param
// @Return error
func (r RedisService) GetStringKey(key string) (string, error) {
	var result string
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	result, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("get ", key, " failed, err:", err)
		return result, err
	}
	return result, nil
}

// @Title GetInfoByKey
// @Description 根据key获得数据
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param
// @Return error
func (r RedisService) GetInfoByKey(key string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	info, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("get ", key, " failed, err:", err)
		return result, err
	}
	if er := json.Unmarshal([]byte(info), &result); er != nil {
		_, err = conn.Do("DEL", key)
		if err != nil {
			log.Error(conf.Config.Base.LogFileName, "").Println("del ", key, " failed, err:", err)
			return result, err
		}
		return result, err
	}
	return result, nil
}

func (r RedisService) GetArrByKey(key string) (result []interface{}, err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	info, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("get ", key, " failed, err:", err)
		return result, err
	}
	if er := json.Unmarshal([]byte(info), &result); er != nil {
		_, err = conn.Do("DEL", key)
		if err != nil {
			log.Error(conf.Config.Base.LogFileName, "").Println("del ", key, " failed, err:", err)
			return result, err
		}
		return result, err
	}
	return result, nil
}

func (r RedisService) Get() (res string, err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	res, err = redis.String(conn.Do("GET", r.Key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	return
}

func (r RedisService) Set() (err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	_, err = conn.Do("SET", r.Key, r.Value)
	if err != nil {
		return
	}
	return
}

// @Title RedisGetSign
// @Description 得到接口验证的sign
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param
// @Return error
func (r RedisService) RedisGetSign() error {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	_, err := redis.String(conn.Do("GET", r.Key))
	if err != nil {
		//log.Error(conf.Config.Base.LogFileName, "").Println("get %s failed, err:%v", r.Key, err)
		return err
	}
	return nil
}

// @Title RedisVerification
// @Description 接口token验证
// @Author xuanshuiyuan 2021-10-31 16:02
// @Param token
// @Return string,error
func (r *RedisService) RedisVerification(key string, token string) error {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	//
	data, err := r.GetInfoByKey(key)
	if err != nil { //登陆失效
		return err
	}
	if data["token"] != token {
		return errors.New("Im Invalid Request!")
	}
	n, err := conn.Do("EXPIRE", key, conf.RedisTokenExp) //重置token时间
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "").Println("expiring ", r.Key, " failed, err:", err)
		return err
	} else if n != int64(1) {
		log.Error(conf.Config.Base.LogFileName, "").Println("expiring ", r.Key, " failed")
		return err
	}
	return nil
}

// @Title Del
// @Description 删除数据
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param token
// @Return string,error
func (r *RedisService) Del() (err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	_, err = conn.Do("DEL", r.Key)
	if err != nil {
		return
	}
	return
}

// @Title Lock
// @Description 获取分布式锁
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param token
// @Return string,error
func (r *RedisService) Lock() (isLock bool, err error) {
	//conf.Mutex.Lock()
	//defer conf.Mutex.Unlock()
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	_, err = redis.String(conn.Do("set", r.Key, 1, "ex", r.Exp, "nx"))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			log.Error(conf.Config.Base.LogFileName, "").Println("set ", r.Key, " failed, err:", err)
			return false, err
		}
		//log.Error(conf.Config.Base.LogFileName, "").Println("set ", r.Key, " failed, err:", err)
		return
	}
	isLock = true
	return
}

// @Title UnLock
// @Description 删除分布式锁
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param token
// @Return string,error
func (r *RedisService) UnLock() (err error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close()
	_, err = conn.Do("DEL", r.Key)
	if err != nil {
		return
	}
	return
}

func Spin(redis *RedisService, key string, exp string, frequency int64) (isLock bool) {
	if frequency >= 60 {
		//自旋超过20次，退出
		return false
	}
	isLock, err := redis.SetKey(key).SetExp(exp).Lock()
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "redis.log").Println("spin ", err)
		return false
	}
	//获取了锁
	if isLock == true {
		return true
	} else { // 自旋 500ms一次
		log.Error(conf.Config.Base.LogFileName, "redis.log").Println("spin:", key, "第", frequency, "次", err)
		time.Sleep(time.Millisecond * 500)
		return Spin(redis, key, exp, frequency+1)
	}
	return false
}
