// @Author xuanshuiyuan
package delay_queue

import (
	"context"
	"fmt"
	"go_project_template/internal/conf"

	"github.com/xuanshuiyuan/delay_queue"
	"github.com/xuanshuiyuan/goxy"
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
	ctx := context.Background()
	if err = delay_queue.NewProducer().RegisterRedis(newRedis()).PushMessage(ctx, taskName, message); err != nil {
		logDelayQueueError(err)
		return err
	}
	return
}

func PushT(taskName, message string, time int64) (err error) {
	ctx := context.Background()
	if err = delay_queue.NewProducer().RegisterRedis(newRedis()).PushMessageT(ctx, taskName, time, message); err != nil {
		logDelayQueueError(err)
		return err
	}
	return
}

func Delete(taskName, message string) (err error) {
	ctx := context.Background()
	if err = delay_queue.NewProducer().RegisterRedis(newRedis()).DeleteMessage(ctx, taskName, message); err != nil {
		logDelayQueueError(err)
		return err
	}
	return
}

func logDelayQueueError(err error) {
	if err == nil {
		return
	}
	if log != nil && conf.Config != nil && conf.Config.Base != nil {
		log.Error(conf.Config.Base.LogFileName, "DelayQueue.log").Println(err)
		return
	}
	fmt.Println("delay queue error:", err)
}
