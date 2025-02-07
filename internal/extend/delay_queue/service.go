// @Author xuanshuiyuan
package delay_queue

import (
	"context"
	"github.com/xuanshuiyuan/delay_queue"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
)

var (
	log *goxy.Logs
)

func New() {
	delay_queue.NewServer().AddTasks([]delay_queue.Task{
		NewGeneral(),
	}).Start()
}

func Push(taskName, message string) (err error) {
	var ctx context.Context
	if err = delay_queue.NewProducer().RegisterRedis(newRedis()).PushMessage(ctx, taskName, message); err != nil {
		log.Error(conf.Config.Base.LogFileName, "DelayQueue.log").Println(err)
		return err
	}
	return
}

func PushT(taskName, message string, time int64) (err error) {
	var ctx context.Context
	if err = delay_queue.NewProducer().RegisterRedis(newRedis()).PushMessageT(ctx, taskName, time, message); err != nil {
		log.Error(conf.Config.Base.LogFileName, "DelayQueue.log").Println(err)
		return err
	}
	return
}

func Delete(taskName, message string) (err error) {
	var ctx context.Context
	if err = delay_queue.NewProducer().RegisterRedis(newRedis()).DeleteMessage(ctx, taskName, message); err != nil {
		log.Error(conf.Config.Base.LogFileName, "DelayQueue.log").Println(err)
		return err
	}
	return
}
