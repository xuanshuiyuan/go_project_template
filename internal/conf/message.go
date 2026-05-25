// @Author xuanshuiyuan
package conf

var MessagePushType = map[string]string{
	"SmsLoginCode": "AlismsGgc-LoginCode",
}

var WxApp = map[string]map[string]interface{}{
	"Ggc": map[string]interface{}{
		"AppId":              "",
		"Secret":             "",
		"MiniprogramAppid":   "",
		"GetAccessTokenCall": "GgcWxAppAccessToken",
		"VehicleEvaluationOrderAuditResults": map[string]interface{}{
			"Code":                "",
			"MiniprogramPagepath": "/pages/content/index",
			"Message":             "车辆估价报告生成通知",
		},
	},
}

var Sms = map[string]map[string]interface{}{
	"AlismsGgc": map[string]interface{}{
		"accessKeyId":     "",
		"accessKeySecret": "",
		"SignName":        "test",
		"Channel":         "NewAlismsGgc",
		"LoginCode": map[string]interface{}{
			"Code":    "SMS_137820295",
			"Message": "验证码%s，您正在登录，若非本人操作，请勿泄露。",
			"Params":  "{\"code\":\"%s\"}",
		},
	},
}
