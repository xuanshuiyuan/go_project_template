// @Author xuanshuiyuan
package alisms

import (
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"os"
)

type Alisms struct {
	Client          *dysmsapi20170525.Client
	AccessKeyId     *string `json:"access_key_id"`
	AccessKeySecret *string `json:"access_key_secret"`
}

var (
	log    *goxy.Logs
	Client *SmsType
)

type SmsType struct {
	Xcc *Alisms
	Ggc *Alisms
}

func NewAlisms() {
	Client = &SmsType{
		//Xcc: NewXccAlisms(),
		Ggc: NewGgcAlisms(),
	}
}

func NewGgcAlisms() *Alisms {
	var init = &Alisms{}
	smsConf, ok := conf.Sms["AlismsGgc"]
	if !ok {
		panic("AlismsGgc config not found")
	}
	accessKeyId, ok := smsConf["accessKeyId"].(string)
	if !ok {
		panic("AlismsGgc accessKeyId not found or not string")
	}
	accessKeySecret, ok := smsConf["accessKeySecret"].(string)
	if !ok {
		panic("AlismsGgc accessKeySecret not found or not string")
	}
	init.SetAccessKeyId(&accessKeyId).SetAccessKeySecret(&accessKeySecret)
	if err := init.NewAlismsClient(); err != nil {
		panic(err)
	}
	return init
}

func (a *Alisms) SetAccessKeyId(accessKeyId *string) *Alisms {
	a.AccessKeyId = accessKeyId
	return a
}

func (a *Alisms) SetAccessKeySecret(accessKeySecret *string) *Alisms {
	a.AccessKeySecret = accessKeySecret
	return a
}

// 使用AK&SK初始化账号Client
// @param accessKeyId
// @param accessKeySecret
// @return Client
// @throws Exception
//
func (a *Alisms) NewAlismsClient() error {
	config := &openapi.Config{
		AccessKeyId:     a.AccessKeyId,
		AccessKeySecret: a.AccessKeySecret,
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	client, err := dysmsapi20170525.NewClient(config)
	if err != nil {
		return fmt.Errorf("Alisms初始化失败: %w", err)
	}
	a.Client = client
	return nil
}

func (a *Alisms) SendSms(phoneNumbers string, signName string, templateCode string, templateParam string,
) (err error) {
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phoneNumbers),
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(templateParam),
	}
	//contain := goxy.Contain(phoneNumbers, conf.SmsTest)
	if os.Getenv("ENV") != "production" {
		return
	}
	// 复制代码运行请自行打印 API 的返回值
	res, _err := a.Client.SendSms(sendSmsRequest)
	if _err != nil {
		log.Data(conf.Config.Base.LogFileName, "Alisms.log").Println(goxy.FmtLog("params.title", sendSmsRequest, "result.title", res))
		return
	}
	if *res.Body.Code != "OK" {
		return errors.New(*res.Body.Message)
	}
	return
}

func (a *Alisms) SendBatchSms(phoneNumbers string, signName string, templateCode string, templateParam string,
) (err error) {
	sendBatchSmsRequest := &dysmsapi20170525.SendBatchSmsRequest{
		PhoneNumberJson:   tea.String(phoneNumbers),
		SignNameJson:      tea.String(signName),
		TemplateCode:      tea.String(templateCode),
		TemplateParamJson: tea.String(templateParam),
	}
	//contain := goxy.Contain(phoneNumbers, conf.SmsTest)
	if os.Getenv("ENV") != "production" {
		return
	}
	// 复制代码运行请自行打印 API 的返回值
	res, _err := a.Client.SendBatchSms(sendBatchSmsRequest)
	log.Data(conf.Config.Base.LogFileName, "Alisms.log").Println(goxy.FmtLog("params.title", sendBatchSmsRequest, "result.title", res))
	if _err != nil {
		return
	}
	if *res.Body.Code != "OK" {
		return errors.New(*res.Body.Message)
	}
	return
}
