// @Author xuanshuiyuan
package delay_queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xuanshuiyuan/delay_queue"
	"reflect"
)

var (
	DelayQueueGeneralDelayTime = 10
	DelayQueueGeneralName      = "delay-queue-General"
)

type DelayQueueGeneralConsumer struct {
}

func NewGeneral() delay_queue.Task {
	return delay_queue.Task{
		Name:        DelayQueueGeneralName,       // 通用
		DelayTime:   DelayQueueGeneralDelayTime,  // 延迟时间 10秒
		Limit:       5,                           // 单个consumer每次最大的消费数量
		Consumer:    DelayQueueGeneralConsumer{}, // 消费者处理程序
		ConsumerNum: 3,                           // 消费者数量
		Redis:       newRedis(),                  //
		AckType:     delay_queue.AckTypeAuto,     // ack类型，自动、手动、禁止
		AckTimeout:  60,                          // 当数据取出来后，如超过此时间还未被ack，数据会被重新消费
	}
}

//MethodByName value
func (d DelayQueueGeneralConsumer) Deal(ctx context.Context, task delay_queue.Task, messages []string) error {
	for _, v := range messages {
		var val = make(map[string]interface{})
		if err := json.Unmarshal([]byte(v), &val); err != nil {
			return err
		}
		res := reflect.ValueOf(d).MethodByName(val["method_name"].(string)).Call([]reflect.Value{reflect.ValueOf(val["value"].(string))})
		if res[0].Interface() != nil && res[0].Interface().(error) != nil {
			err := res[0].Interface().(error)
			return err
		}
	}
	return nil
}

func (d DelayQueueGeneralConsumer) Error(ctx context.Context, task delay_queue.Task, err *delay_queue.Error) {
	fmt.Println("DelayQueueGeneralConsumer Error", *err)
}

func Add(method_name, message string, time int64) {
	var val = make(map[string]interface{})
	val["method_name"] = method_name
	val["value"] = message
	params, _ := json.Marshal(val)
	PushT(DelayQueueGeneralName, string(params), int64(int(time)-DelayQueueGeneralDelayTime))
}

func Deletes(method_name, message string) (err error) {
	var val = make(map[string]interface{})
	val["method_name"] = method_name
	val["value"] = message
	params, _ := json.Marshal(val)
	Delete(DelayQueueGeneralName, string(params))
	return
}
