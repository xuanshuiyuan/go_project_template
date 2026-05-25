// @Author xuanshuiyuan
package delay_queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuanshuiyuan/delay_queue"
	"reflect"
)

var (
	DelayQueueGeneralDelayTime = 10
	DelayQueueGeneralName      = "delay-queue-General"

	ErrGeneralMessageInvalid   = errors.New("general delay queue message invalid")
	ErrGeneralMethodNameEmpty  = errors.New("general delay queue method_name is empty")
	ErrGeneralMethodNotFound   = errors.New("general delay queue method not found")
	ErrGeneralMethodSignatures = errors.New("general delay queue method signature invalid")
)

type DelayQueueGeneralConsumer struct {
}

type generalMessage struct {
	MethodName string `json:"method_name"`
	Value      string `json:"value"`
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

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
		payload, err := decodeGeneralMessage(v)
		if err != nil {
			return err
		}
		if err := d.call(payload.MethodName, payload.Value); err != nil {
			return err
		}
	}
	return nil
}

func (d DelayQueueGeneralConsumer) Error(ctx context.Context, task delay_queue.Task, err *delay_queue.Error) {
	fmt.Println("DelayQueueGeneralConsumer Error", *err)
}

func Add(method_name, message string, time int64) {
	params, err := json.Marshal(generalMessage{
		MethodName: method_name,
		Value:      message,
	})
	if err != nil {
		return
	}
	_ = PushT(DelayQueueGeneralName, string(params), int64(int(time)-DelayQueueGeneralDelayTime))
}

func Deletes(method_name, message string) (err error) {
	params, err := json.Marshal(generalMessage{
		MethodName: method_name,
		Value:      message,
	})
	if err != nil {
		return err
	}
	_ = Delete(DelayQueueGeneralName, string(params))
	return
}

func decodeGeneralMessage(raw string) (generalMessage, error) {
	var payload generalMessage
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return generalMessage{}, fmt.Errorf("%w: %v", ErrGeneralMessageInvalid, err)
	}
	if payload.MethodName == "" {
		return generalMessage{}, ErrGeneralMethodNameEmpty
	}
	return payload, nil
}

func (d DelayQueueGeneralConsumer) call(methodName, value string) error {
	method := reflect.ValueOf(d).MethodByName(methodName)
	if !method.IsValid() {
		return fmt.Errorf("%w: %s", ErrGeneralMethodNotFound, methodName)
	}
	methodType := method.Type()
	if methodType.NumIn() != 1 || methodType.In(0).Kind() != reflect.String {
		return fmt.Errorf("%w: %s", ErrGeneralMethodSignatures, methodName)
	}
	if methodType.NumOut() != 1 || !methodType.Out(0).Implements(errorType) {
		return fmt.Errorf("%w: %s", ErrGeneralMethodSignatures, methodName)
	}
	results := method.Call([]reflect.Value{reflect.ValueOf(value)})
	if results[0].IsNil() {
		return nil
	}
	return fmt.Errorf("general method %s: %w", methodName, results[0].Interface().(error))
}

