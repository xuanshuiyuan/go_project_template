// @Author xuanshuiyuan
package delay_queue

import (
	"context"
	"go_project_template/internal/engine"
	"time"
)

//Redis redis的配置
type defaultRedis struct {
}

//初始化
func newRedis() *defaultRedis {
	return &defaultRedis{}
}

func (d defaultRedis) execCommand(command string, args ...interface{}) (interface{}, error) {
	conn := engine.DB.Redis.Get()
	defer conn.Close() //需要close
	reply, err := conn.Do(command, args...)
	return reply, err
}

func (d defaultRedis) ZAddT(ctx context.Context, key string, time int64, messages ...string) error {
	args := []interface{}{
		key,
		"NX",
		time,
	}
	for _, message := range messages {
		args = append(args, message)
	}
	_, err := d.execCommand("ZADD", args...)
	return err
}

func (d defaultRedis) ZAdd(ctx context.Context, key string, messages ...string) error {
	args := []interface{}{
		key,
		"NX",
		time.Now().Unix(),
	}
	for _, message := range messages {
		args = append(args, message)
	}
	_, err := d.execCommand("ZADD", args...)
	return err
}

func (d defaultRedis) ZRem(ctx context.Context, key string, messages ...string) error {
	args := []interface{}{
		key,
	}
	for _, message := range messages {
		args = append(args, message)
	}
	_, err := d.execCommand("ZREM", args...)
	return err
}

func (d defaultRedis) EvalSha(ctx context.Context, sha1 string, values []interface{}) (interface{}, error) {
	args := []interface{}{
		sha1,
	}
	args = append(args, values...)
	res, err := d.execCommand("EVALSHA", args...)
	return res, err
}

func (d defaultRedis) LoadScript(ctx context.Context, script string) error {
	args := []interface{}{
		"LOAD",
		script,
	}
	_, err := d.execCommand("SCRIPT", args...)
	return err
}
