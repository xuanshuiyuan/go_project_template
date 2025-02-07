// @Author xuanshuiyuan 2023/4/26 10:31:00
package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"go_project_template/internal/third_party/alisms"
	"reflect"
	"strings"
)

type Alisms struct {
	push *alisms.Alisms
}

func NewAlisms() *Alisms {
	return &Alisms{}
}

//若使用其他短信服务商或者增加需求，直接自行实现接口即可
func (a *Alisms) NewAlismsXcc(p *PushService) {
	p.Channel = &Alisms{
		push: alisms.Client.Xcc,
	}
	p.Sender = "阿里云短信"
}

func (a *Alisms) NewAlismsGgc(p *PushService) {
	p.Channel = &Alisms{
		push: alisms.Client.Ggc,
	}
	p.Sender = "阿里云短信"
}

func LoadSmsConfig() *PushService {
	return &PushService{
		Channel:               NewAlisms(),
		channel:               1, //1:sms,2:微信
		MessagePushType:       conf.MessagePushType,
		MessagePushTypeParams: conf.Sms,
	}
}

func (a *Alisms) Send(args ...interface{}) (err error) {
	if err = goxy.WithTimeout(func() error {
		return a.push.SendSms(args[0].(string), args[1].(string), args[2].(string), args[3].(string))
	}); err != nil {
		return err
	}
	return
}

func (a *Alisms) SendBatch(args [][]string) (err error) {
	var phoneNumbers, signName, templateParam []string
	var templateCode string
	for _, v := range args {
		templateCode = v[2]
		phoneNumbers = append(phoneNumbers, v[0])
		signName = append(signName, v[1])
		templateParam = append(templateParam, v[3])
	}
	phoneNumbersJson, _ := json.Marshal(phoneNumbers)
	signNameJson, _ := json.Marshal(signName)
	templateParamJson, _ := json.Marshal(templateParam)
	if err = a.push.SendBatchSms(string(phoneNumbersJson), string(signNameJson), templateCode, string(templateParamJson)); err != nil {
		return
	}
	return
}

//phoneNumbers string, signName string, templateCode string, templateParam string, message string
func (a *Alisms) parseTemplate(p *PushService) (result []string, err error) {
	defer func() {
		if errs := recover(); errs != nil {
			err = errors.New("参数错误")
			return
		}
	}()
	messagePushType := p.MessagePushType[p.MessagePushCode]
	if messagePushType == "" {
		err = errors.New("参数错误")
		return
	}
	messagePushTypeArr := strings.Split(messagePushType, "-")
	var config = p.MessagePushTypeParams[messagePushTypeArr[0]]
	if len(config) == 0 {
		err = errors.New("参数错误")
		return
	}
	signName := config["SignName"].(string)
	result = append(result, signName)
	var template = config[messagePushTypeArr[1]].(map[string]interface{})
	if len(template) == 0 {
		err = errors.New("参数错误")
		return
	}
	templateCode := template["Code"].(string)
	result = append(result, templateCode)
	if template["Params"].(string) != "" {
		params := fmt.Sprintf(template["Params"].(string), goxy.StringToInterface(p.Content[0][1:])...)
		result = append(result, params)
	} else {
		result = append(result, "")
	}
	result = append(result, template["Code"].(string))
	message := fmt.Sprintf(template["Message"].(string), goxy.StringToInterface(p.Content[0][1:])...)
	result = append(result, message)
	channel := config["Channel"].(string)
	reflect.ValueOf(a).MethodByName(channel).Call([]reflect.Value{reflect.ValueOf(p)})
	return
}
