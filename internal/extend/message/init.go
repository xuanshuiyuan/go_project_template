// @Author xuanshuiyuan
package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/engine"
	"go_project_template/internal/models"
	"strings"
	"time"
)

var log *goxy.Logs

//消息推送
type MessageEr interface {
	Send() error      //message_type,receiver,params
	SendBatch() error //message_type,[receiver,params]
}

//消息推送服务商
type ChannelEr interface {
	Send(...interface{}) error //message_type,receiver,params
	SendBatch([][]string) error
	parseTemplate(*PushService) ([]string, error) //格式化消息 eg: [Receiver,Params...,LogMessage]
}

type PushService struct {
	PushParams
	Channel               ChannelEr
	channel               int8
	Sender                string
	MessagePushType       map[string]string
	MessagePushTypeParams map[string]map[string]interface{}
}

type PushParams struct {
	MessagePushCode string
	Content         [][]string
}

func (p *PushService) InjectParams(ps *PushParams) *PushService {
	p.PushParams = *ps
	return p
}

func Push(m MessageEr) error {
	return m.Send()
}

func PushBatch(m MessageEr) error {
	return m.SendBatch()
}

// @Title Send
// @Description 消息推送
func (p *PushService) Send() (err error) {
	if len(p.Content) == 0 {
		return errors.New("参数不能为空")
	}
	content, err := p.Channel.parseTemplate(p)
	if err != nil {
		return
	}
	params, _ := json.Marshal(p.Content[0])
	logs := append([]models.MessagePushQueue{}, models.MessagePushQueue{
		Sender:       p.Sender,
		Receiver:     p.Content[0][0],
		Channel:      p.channel,
		TemplateCode: content[len(content)-2],
		Type:         p.MessagePushCode,
		Content:      content[len(content)-1],
		Params:       string(params),
		Status:       1,
	})
	if err = p.Channel.Send(goxy.StringToInterface(append([]string{p.Content[0][0]}, content...))...); err != nil {
		logs[0].Reason = err.Error()
		logs[0].Status = 2
		addLogs(logs)
		return err
	}
	//插入日志
	addLogs(logs)
	return nil
}

// @Title SendBatch
// @Description 批量发送
func (p *PushService) SendBatch() (err error) {
	if len(p.Content) == 0 {
		return errors.New("参数不能为空")
	}
	contents := [][]string{}
	var logs []models.MessagePushQueue
	for _, v := range p.Content {
		pv := &PushService{}
		pv = p
		pv.Content = [][]string{v}
		content, err := p.Channel.parseTemplate(pv)
		if err != nil {
			return err
		}
		contents = append(contents, append([]string{v[0]}, content...))
		params, _ := json.Marshal(v)
		logs = append(logs, models.MessagePushQueue{
			Sender:       p.Sender,
			Receiver:     v[0],
			Channel:      p.channel,
			Type:         p.MessagePushCode,
			TemplateCode: content[len(content)-2],
			Content:      content[len(content)-1],
			Params:       string(params),
			Status:       1,
		})
	}
	if err = p.Channel.SendBatch(contents); err != nil {
		for k, _ := range logs {
			logs[k].Status = 2
			logs[k].Reason = err.Error()
		}
		addLogs(logs)
		return
	}
	//插入日志
	addLogs(logs)
	return
}

// @Title addLogs
// @Description 消息推送队列
func addLogs(data []models.MessagePushQueue) {
	sql := []string{}
	for _, v := range data {
		sql = append(sql, fmt.Sprintf("('%s', '%s', %d, '%s', '%s','%s', '%s', %d, '%s', %d)", v.Sender, v.Receiver, v.Channel, v.Type, v.TemplateCode, v.Content, v.Params, v.Status, v.Reason, time.Now().Unix()))
	}
	sqls := strings.Join(sql, ",")
	exec := "INSERT INTO `message_push_queue` (sender,receiver,channel,type,template_code,content,params,status,reason,create_time) VALUES  " + sqls
	engine.DB.Mysql.Exec(exec)
}
